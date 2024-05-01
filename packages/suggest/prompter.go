package suggest

import (
	"errors"
	"fmt"
	"strings"
)

type MessageSuggestionPrompter struct {
	embed    string // embed string that could be use to embed user's context and messages to make complete prompt
	userData UserData
	messages []Message
}

func (p MessageSuggestionPrompter) Build() (string, error) {
	msgs := []string{}
	for _, msg := range p.messages {
		msgs = append(msgs, fmt.Sprintf("%s: %s", msg.Sender, msg.Content))
	}
	return fmt.Sprintf(
		p.embed,
		p.userData.Learning.Lang,
		p.userData.Learning.Level,
		strings.Join(msgs, "\n"),
	), nil
}

func (p *MessageSuggestionPrompter) Update(objs ...any) error {
	for _, obj := range objs {
		switch doc := obj.(type) {
		case UserData:
			p.userData = doc
		case []Message:
			p.messages = doc
		default:
			return errors.New("messageSuggestionPrompter: expected(common.UserContext, []common.Message) got unknown")
		}
	}
	return nil
}

func NewMessageSuggestionPrompter() *MessageSuggestionPrompter {
	embed, _ := randomEmbed()

	return &MessageSuggestionPrompter{
		embed: embed,
	}
}
