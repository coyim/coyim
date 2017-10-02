package plist

import (
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"math"
	"runtime"
	"strings"
	"time"
)

const xmlDOCTYPE = `<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
`

type xmlPlistGenerator struct {
	writer     io.Writer
	xmlEncoder *xml.Encoder
}

func (p *xmlPlistGenerator) generateDocument(root cfValue) {
	io.WriteString(p.writer, xml.Header)
	io.WriteString(p.writer, xmlDOCTYPE)

	plistStartElement := xml.StartElement{
		Name: xml.Name{
			Space: "",
			Local: "plist",
		},
		Attr: []xml.Attr{{
			Name: xml.Name{
				Space: "",
				Local: "version"},
			Value: "1.0"},
		},
	}

	p.xmlEncoder.EncodeToken(plistStartElement)

	p.writePlistValue(root)

	p.xmlEncoder.EncodeToken(plistStartElement.End())
	p.xmlEncoder.Flush()
}

func (p *xmlPlistGenerator) writeDictionary(dict *cfDictionary) {
	dict.sort()
	startElement := xml.StartElement{Name: xml.Name{Local: "dict"}}
	p.xmlEncoder.EncodeToken(startElement)
	for i, k := range dict.keys {
		p.xmlEncoder.EncodeElement(k, xml.StartElement{Name: xml.Name{Local: "key"}})
		p.writePlistValue(dict.values[i])
	}
	p.xmlEncoder.EncodeToken(startElement.End())
}

func (p *xmlPlistGenerator) writeArray(a *cfArray) {
	startElement := xml.StartElement{Name: xml.Name{Local: "array"}}
	p.xmlEncoder.EncodeToken(startElement)
	for _, v := range a.values {
		p.writePlistValue(v)
	}
	p.xmlEncoder.EncodeToken(startElement.End())
}

func (p *xmlPlistGenerator) writePlistValue(pval cfValue) {
	if pval == nil {
		return
	}

	defer p.xmlEncoder.Flush()

	if dict, ok := pval.(*cfDictionary); ok {
		p.writeDictionary(dict)
		return
	} else if a, ok := pval.(*cfArray); ok {
		p.writeArray(a)
		return
	} else if uid, ok := pval.(cfUID); ok {
		p.writeDictionary(&cfDictionary{
			keys: []string{"CF$UID"},
			values: []cfValue{
				&cfNumber{
					signed: false,
					value:  uint64(uid),
				},
			},
		})
		return
	}

	// Everything here and beyond is encoded the same way: <key>value</key>
	key := ""
	var encodedValue interface{} = pval

	switch pval := pval.(type) {
	case cfString:
		key = "string"
	case *cfNumber:
		key = "integer"
		if pval.signed {
			encodedValue = int64(pval.value)
		} else {
			encodedValue = pval.value
		}
	case *cfReal:
		key = "real"
		encodedValue = pval.value
		switch {
		case math.IsInf(pval.value, 1):
			encodedValue = "inf"
		case math.IsInf(pval.value, -1):
			encodedValue = "-inf"
		case math.IsNaN(pval.value):
			encodedValue = "nan"
		}
	case cfBoolean:
		key = "false"
		b := bool(pval)
		if b {
			key = "true"
		}
		encodedValue = ""
	case cfData:
		key = "data"
		encodedValue = xml.CharData(base64.StdEncoding.EncodeToString([]byte(pval)))
	case cfDate:
		key = "date"
		encodedValue = time.Time(pval).In(time.UTC).Format(time.RFC3339)
	}

	if key != "" {
		err := p.xmlEncoder.EncodeElement(encodedValue, xml.StartElement{Name: xml.Name{Local: key}})
		if err != nil {
			panic(err)
		}
	}
}

func (p *xmlPlistGenerator) Indent(i string) {
	p.xmlEncoder.Indent("", i)
}

func newXMLPlistGenerator(w io.Writer) *xmlPlistGenerator {
	mw := mustWriter{w}
	return &xmlPlistGenerator{mw, xml.NewEncoder(mw)}
}

type xmlPlistParser struct {
	reader             io.Reader
	xmlDecoder         *xml.Decoder
	whitespaceReplacer *strings.Replacer
	ntags              int
}

func (p *xmlPlistParser) parseDocument() (pval cfValue, parseError error) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			if _, ok := r.(invalidPlistError); ok {
				parseError = r.(error)
			} else {
				// Wrap all non-invalid-plist errors.
				parseError = plistParseError{"XML", r.(error)}
			}
		}
	}()
	for {
		if token, err := p.xmlDecoder.Token(); err == nil {
			if element, ok := token.(xml.StartElement); ok {
				pval = p.parseXMLElement(element)
				if p.ntags == 0 {
					panic(invalidPlistError{"XML", errors.New("no elements encountered")})
				}
				return
			}
		} else {
			// The first XML parse turned out to be invalid:
			// we do not have an XML property list.
			panic(invalidPlistError{"XML", err})
		}
	}
}

func (p *xmlPlistParser) parseXMLElement(element xml.StartElement) cfValue {
	var charData xml.CharData
	switch element.Name.Local {
	case "plist":
		p.ntags++
		for {
			token, err := p.xmlDecoder.Token()
			if err != nil {
				panic(err)
			}

			if el, ok := token.(xml.EndElement); ok && el.Name.Local == "plist" {
				break
			}

			if el, ok := token.(xml.StartElement); ok {
				return p.parseXMLElement(el)
			}
		}
		return nil
	case "string":
		p.ntags++
		err := p.xmlDecoder.DecodeElement(&charData, &element)
		if err != nil {
			panic(err)
		}

		return cfString(charData)
	case "integer":
		p.ntags++
		err := p.xmlDecoder.DecodeElement(&charData, &element)
		if err != nil {
			panic(err)
		}

		s := string(charData)
		if len(s) == 0 {
			panic(errors.New("invalid empty <integer/>"))
		}

		if s[0] == '-' {
			s, base := unsignedGetBase(s[1:])
			n := mustParseInt("-"+s, base, 64)
			return &cfNumber{signed: true, value: uint64(n)}
		} else {
			s, base := unsignedGetBase(s)
			n := mustParseUint(s, base, 64)
			return &cfNumber{signed: false, value: n}
		}
	case "real":
		p.ntags++
		err := p.xmlDecoder.DecodeElement(&charData, &element)
		if err != nil {
			panic(err)
		}

		n := mustParseFloat(string(charData), 64)
		return &cfReal{wide: true, value: n}
	case "true", "false":
		p.ntags++
		p.xmlDecoder.Skip()

		b := element.Name.Local == "true"
		return cfBoolean(b)
	case "date":
		p.ntags++
		err := p.xmlDecoder.DecodeElement(&charData, &element)
		if err != nil {
			panic(err)
		}

		t, err := time.ParseInLocation(time.RFC3339, string(charData), time.UTC)
		if err != nil {
			panic(err)
		}

		return cfDate(t)
	case "data":
		p.ntags++
		err := p.xmlDecoder.DecodeElement(&charData, &element)
		if err != nil {
			panic(err)
		}

		str := p.whitespaceReplacer.Replace(string(charData))

		l := base64.StdEncoding.DecodedLen(len(str))
		bytes := make([]uint8, l)
		l, err = base64.StdEncoding.Decode(bytes, []byte(str))
		if err != nil {
			panic(err)
		}

		return cfData(bytes[:l])
	case "dict":
		p.ntags++
		var key *string
		keys := make([]string, 0, 32)
		values := make([]cfValue, 0, 32)
		for {
			token, err := p.xmlDecoder.Token()
			if err != nil {
				panic(err)
			}

			if el, ok := token.(xml.EndElement); ok && el.Name.Local == "dict" {
				if key != nil {
					panic(errors.New("missing value in dictionary"))
				}
				break
			}

			if el, ok := token.(xml.StartElement); ok {
				if el.Name.Local == "key" {
					var k string
					p.xmlDecoder.DecodeElement(&k, &el)
					key = &k
				} else {
					if key == nil {
						panic(errors.New("missing key in dictionary"))
					}
					keys = append(keys, *key)
					values = append(values, p.parseXMLElement(el))
					key = nil
				}
			}
		}

		if len(keys) == 1 && keys[0] == "CF$UID" && len(values) == 1 {
			if integer, ok := values[0].(*cfNumber); ok {
				return cfUID(integer.value)
			}
		}

		return &cfDictionary{keys: keys, values: values}
	case "array":
		p.ntags++
		values := make([]cfValue, 0, 10)
		for {
			token, err := p.xmlDecoder.Token()
			if err != nil {
				panic(err)
			}

			if el, ok := token.(xml.EndElement); ok && el.Name.Local == "array" {
				break
			}

			if el, ok := token.(xml.StartElement); ok {
				values = append(values, p.parseXMLElement(el))
			}
		}
		return &cfArray{values}
	}
	err := fmt.Errorf("encountered unknown element %s", element.Name.Local)
	if p.ntags == 0 {
		// If out first XML tag is invalid, it might be an openstep data element, ala <abab> or <0101>
		panic(invalidPlistError{"XML", err})
	}
	panic(err)
}

func newXMLPlistParser(r io.Reader) *xmlPlistParser {
	return &xmlPlistParser{r, xml.NewDecoder(r), strings.NewReplacer("\t", "", "\n", "", " ", "", "\r", ""), 0}
}
