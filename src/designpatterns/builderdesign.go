package designpatterns

import (
	"encoding/json"
	"encoding/xml"
)

type Message struct {
	Body   []byte
	Format string
}

type MessageBuilder interface {
	SetRecipient(recipient string)
	SetText(text string)
	Message() (*Message, error)
}

type JSONMessageBuilder struct {
	messageRecipient string
	messageText      string
}

func (b *JSONMessageBuilder) SetRecipient(recipient string) {
	b.messageRecipient = recipient
}

func (b *JSONMessageBuilder) SetText(text string) {
	b.messageText = text
}

func (b *JSONMessageBuilder) Message() (*Message, error) {
	m := make(map[string]string)
	m["recipient"] = b.messageRecipient
	m["message"] = b.messageText

	data, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	return &Message{Body: data, Format: "JSON"}, nil
}

type XMLMessageBuilder struct {
	messageRecipient string
	messageText      string
}

func (b *XMLMessageBuilder) SetRecipient(recipient string) {
	b.messageRecipient = recipient
}

func (b *XMLMessageBuilder) SetText(text string) {
	b.messageText = text
}

func (b *XMLMessageBuilder) Message() (*Message, error) {
	type XMLMessage struct {
		Recipient string `xml:"recipient"`
		Text      string `xml:"body"`
	}

	m := XMLMessage{
		Recipient: b.messageRecipient,
		Text:      b.messageText,
	}

	data, err := xml.Marshal(m)
	if err != nil {
		return nil, err
	}

	return &Message{Body: data, Format: "XML"}, nil
}

type Sender struct{}

func (s *Sender) BuildMessage(builder MessageBuilder) (*Message, error) {
	builder.SetRecipient("Santa Claus")
	builder.SetText("I have tried to be good all year and hope that you and " +
		"your reindeers will be able to deliver me a nice present.")
	return builder.Message()
}
