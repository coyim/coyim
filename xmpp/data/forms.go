package data

import "encoding/xml"

// FormCallback is the type of a function called to process a form. The
// argument is a list of pointers to FormField types. The function should type
// cast the elements, prompt the user and fill in the result field in each
// struct.
type FormCallback func(title, instructions string, fields []interface{}) error

// Form contains the definition for a data submission
type Form struct {
	XMLName      xml.Name     `xml:"jabber:x:data x"`
	Type         string       `xml:"type,attr"`
	Title        string       `xml:"title,omitempty"`
	Instructions string       `xml:"instructions,omitempty"`
	Fields       []FormFieldX `xml:"field"`
}

// FormFieldX contains form field information
type FormFieldX struct {
	XMLName  xml.Name            `xml:"field"`
	Desc     string              `xml:"desc,omitempty"`
	Var      string              `xml:"var,attr"`
	Type     string              `xml:"type,attr,omitempty"`
	Label    string              `xml:"label,attr,omitempty"`
	Required *FormFieldRequiredX `xml:"required"`
	Values   []string            `xml:"value"`
	Options  []FormFieldOptionX  `xml:"option"`
	Media    []FormFieldMediaX   `xml:"media"`
}

// FormFieldMediaX contains form field media information
type FormFieldMediaX struct {
	XMLName xml.Name    `xml:"urn:xmpp:media-element media"`
	URIs    []MediaURIX `xml:"uri"`
}

// MediaURIX contains information about a Media URI
type MediaURIX struct {
	XMLName  xml.Name `xml:"urn:xmpp:media-element uri"`
	MIMEType string   `xml:"type,attr,omitempty"`
	URI      string   `xml:",chardata"`
}

// FormFieldRequiredX contains information about whether a form field is required
type FormFieldRequiredX struct {
	XMLName xml.Name `xml:"required"`
}

// FormFieldOptionX contains a form field option
type FormFieldOptionX struct {
	Label string `xml:"var,attr,omitempty"`
	Value string `xml:"value"`
}

// FormField is the type of a generic form field. One should type cast to a
// specific type of field before processing.
type FormField struct {
	// Label is a human readable label for this field.
	Label string
	// Type is the XMPP-internal type of this field. One should type cast
	// rather than inspect this.
	Type string
	// Name gives the internal name of the field.
	Name string
	// Required specifies is the field is required.
	Required bool
	// Media contains one of more items of media associated with this
	// field and, for each item, one or more representations of it.
	Media [][]Media
}

// Media contains a specific media uri and data
type Media struct {
	MIMEType string
	// URI contains a URI to the data. It may be empty if Data is not.
	URI string
	// Data contains the raw data itself. It may be empty if URI is not.
	Data []byte
}

// FixedFormField is used to indicate a section heading. It's for the form to
// send data to the user rather than the other way around.
type FixedFormField struct {
	FormField

	Text string
}

// BooleanFormField is for a yes/no answer. The Result member should be set to
// the user's answer.
type BooleanFormField struct {
	FormField

	Result bool
}

// TextFormField is for the entry of a single textual item. The Result member
// should be set to the data entered.
type TextFormField struct {
	FormField

	Default string
	Result  string

	// Private is true if this is a password or other sensitive entry.
	Private bool
}

// MultiTextFormField is for the entry of a several textual items. The Results
// member should be set to the data entered.
type MultiTextFormField struct {
	FormField

	Defaults []string
	Results  []string
}

// SelectionFormField asks the user to pick a single element from a set of
// choices. The Result member should be set to an index of the Values array.
type SelectionFormField struct {
	FormField

	Values []string
	Ids    []string
	Result int
}

// MultiSelectionFormField asks the user to pick a subset of possible choices.
// The Result member should be set to a series of indexes of the Results array.
type MultiSelectionFormField struct {
	FormField

	Values  []string
	Ids     []string
	Results []int
}
