package gui

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"text/template"
	"text/template/parse"

	"github.com/coyim/gotk3adapter/gtka"
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/coyim/gotk3extra"
)

const (
	presentationFormatTypeNickname    = "nickname"
	presentationFormatTypeAffiliation = "affiliation"
	presentationFormatTypeRole        = "role"

	presentationLeftDelim  = "{{"
	presentationRightDelim = "}}"
)

var presentationFormatTypes = []string{
	presentationFormatTypeNickname,
	presentationFormatTypeAffiliation,
	presentationFormatTypeRole,
}

var presentationFormatsHighlight = map[string]infoBarHighlightType{
	presentationFormatTypeNickname:    infoBarHighlightNickname,
	presentationFormatTypeAffiliation: infoBarHighlightAffiliation,
	presentationFormatTypeRole:        infoBarHighlightRole,
}

var presentationFormatFuncsMap = template.FuncMap{
	presentationFormatTypeNickname:    true,
	presentationFormatTypeAffiliation: true,
	presentationFormatTypeRole:        true,
}

type presentationTextFormat struct {
	typ    string
	value  string
	start  int
	length int
}

func (f *presentationTextFormat) startIndex() int {
	return f.start
}

func (f *presentationTextFormat) endIndex() int {
	return f.startIndex() + f.length
}

const presentationTextFormatterParseName = "parse"

type presentationTextFormatter struct {
	text string

	formats     []*presentationTextFormat
	formatsLock sync.Mutex

	leftDelim  string
	rightDelim string

	stringBuilder *strings.Builder
}

func newPresentationTextFormatter(text string) *presentationTextFormatter {
	formatter := &presentationTextFormatter{
		text:          text,
		stringBuilder: &strings.Builder{},
		leftDelim:     presentationLeftDelim,
		rightDelim:    presentationRightDelim,
	}

	if err := formatter.parse(); err != nil {
		formatter.error(err)
	}

	return formatter
}

func (f *presentationTextFormatter) parse() error {
	trees, err := parse.Parse(
		presentationTextFormatterParseName,
		f.text,
		f.leftDelim,
		f.rightDelim,
		presentationFormatFuncsMap,
	)

	if err != nil {
		return err
	}

	tree, ok := trees[presentationTextFormatterParseName]
	if !ok {
		return errors.New("invalid tree key")
	}

	return f.parseTree(tree.Copy())
}

func (f *presentationTextFormatter) parseTree(tree *parse.Tree) error {
	for _, node := range tree.Root.Nodes {
		if err := f.parseTreeNode(node); err != nil {
			return err
		}
	}
	return nil
}

func (f *presentationTextFormatter) parseTreeNode(node parse.Node) error {
	switch n := node.(type) {
	case *parse.TextNode:
		f.writeText(n.String())
	case *parse.ActionNode:
		for _, cmd := range n.Pipe.Cmds {
			firstWord := cmd.Args[0]
			switch n := firstWord.(type) {
			case *parse.IdentifierNode:
				f.evalIdentifierNode(n, cmd.Args)
			default:
				return errors.New("invalid node type")
			}
		}
	}

	return nil
}

const trimChars = "`\""

func (f *presentationTextFormatter) evalIdentifierNode(node parse.Node, args []parse.Node) {
	if args != nil {
		args = args[1:]
	}

	realVal := strings.Trim(args[0].String(), trimChars)

	switch typ := node.String(); typ {
	case presentationFormatTypeNickname,
		presentationFormatTypeAffiliation,
		presentationFormatTypeRole:
		if err := f.addFormBasedOnNode(typ, node, realVal); err != nil {
			f.error(err)
		}

	default:
		f.writeText(realVal)
	}
}

func (f *presentationTextFormatter) writeText(content string) (int, error) {
	return f.stringBuilder.WriteString(content)
}

func (f *presentationTextFormatter) error(err error) {
	// [ip] We need to log the given error
}

func (f *presentationTextFormatter) String() string {
	return f.stringBuilder.String()
}

func (f *presentationTextFormatter) addFormBasedOnNode(formatType string, node parse.Node, text string) error {
	f.formatsLock.Lock()
	defer f.formatsLock.Unlock()

	textLengthBeforeUpdate := len(f.String())

	length, err := f.writeText(text)
	if err != nil {
		return err
	}

	format := &presentationTextFormat{
		typ:    formatType,
		value:  text,
		start:  textLengthBeforeUpdate,
		length: length,
	}

	f.formats = append(f.formats, format)

	return nil
}

const textFormat = "%s"

// [ip] The following implementation is going to change to avoid using "gtka" and "gotk3extra" inside the "gui" package
func (f *presentationTextFormatter) formatLabel(label gtki.Label) {
	f.formatsLock.Lock()
	defer f.formatsLock.Unlock()

	label.SetText(fmt.Sprintf(textFormat, f))

	allAttributes := gotk3extra.PangoAttrListNew()

	for _, format := range f.formats {
		if highlightType, ok := presentationFormatsHighlight[format.typ]; ok {
			copy := newInfoBarHighlightAttributes(highlightType)
			copyAttributesTo(allAttributes, copy, format.startIndex(), format.endIndex())
		}
	}

	gotk3extra.LabelSetAttributes(gtka.UnwrapLabel(label), allAttributes)
}

func copyAttributesTo(toAttrList, fromAttrList *gotk3extra.PangoAttrList, startIndex, endIndex int) {
	for _, attr := range gotk3extra.PangoGetAttributesFromList(fromAttrList) {
		attr.StartIndex(uint(startIndex))
		attr.EndIndex(uint(endIndex))
		toAttrList.Insert(attr)
	}
}
