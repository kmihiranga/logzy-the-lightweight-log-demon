package entity

// TextObject represents markdown text
type TextObject struct {
	Type string  `json:"type"`
	Text *string `json:"text"`
}

type HeaderText struct {
	Type  string `json:"type"`
	Text  string `json:"text"`
	Emoji bool   `json:"emoji"`
}

type MarkdownText struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type Element struct {
	Type     string       `json:"type"`
	Border   int          `json:"border"`
	Elements []TextObject `json:"elements"`
}

type Block struct {
	Type    string         `json:"type"`
	Text    Text           `json:"text,omitempty"`
	Fields  []MarkdownText `json:"fields,omitempty"`
	Element []Element      `json:"elements,omitempty"`
}

type SlackMessage struct {
	Text   string  `json:"text"`
	Blocks []Block `json:"blocks"`
}

func (ht *HeaderText) TextType() string {
	return ht.Type
}

func (to *TextObject) TextType() string {
	return to.Type
}

// Text interface to be implemented by different text types
type Text interface {
	TextType() string
}
