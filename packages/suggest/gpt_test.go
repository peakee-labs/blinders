package suggest

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/assert"
)

var authToken = os.Getenv("OPENAI_API_KEY")

var suggester, _ = NewGPTSuggester(openai.NewClient(authToken))

func TestTextCompletion(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	prompt := "Just reply 'hello, world!'"
	suggestions, err := suggester.TextCompletion(ctx, UserData{}, prompt)
	assert.Nil(t, err)
	assert.Equal(t, suggester.nText, len(suggestions))

	fmt.Println(suggestions)
}

func TestSuggest(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	sender := "User1"
	receiver := "User2"
	userContext := newUserContext(
		sender,
		Language{
			Lang:  LangVi,
			Level: Advanced,
		},
		Language{
			Lang:  LangEn,
			Level: Beginner,
		},
	)
	msgs := []Message{
		*NewMessage(sender, receiver, "Hello, how are you?"),
		*NewMessage(receiver, sender, "Fine, how about you?"),
		*NewMessage(sender, receiver, "Too. Did you come to the class yesterday?"),
		*NewMessage(receiver, sender, "Yes, yesterday the teacher gave the students some homework."),
	}

	suggestions, err := suggester.ChatCompletion(ctx, userContext, msgs)
	assert.Nil(t, err)
	assert.NotNil(t, suggestions)
	assert.Equal(t, suggester.nChat, len(suggestions))

	// TODO: would be better to check the response format
	for _, suggestion := range suggestions {
		fmt.Printf("suggestion: %v\n", suggestion)
	}
}

func newUserContext(ID string, native Language, language Language) UserData {
	return UserData{
		UserID:   ID,
		Native:   native,
		Learning: language,
	}
}
