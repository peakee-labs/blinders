package suggest

import (
	"context"
	"errors"

	openai "github.com/sashabaranov/go-openai"
)

type GPTSuggesterOptions struct {
	prompter         Prompter
	chatModel        string
	textModel        string
	nChat            int
	nText            int
	modelTemperature float32
}

var DefaultSuggesterOptions = GPTSuggesterOptions{
	prompter:         NewMessageSuggestionPrompter(),
	chatModel:        openai.GPT3Dot5TurboInstruct,
	textModel:        openai.GPT3Dot5TurboInstruct,
	nChat:            2,
	nText:            1,
	modelTemperature: 0.4,
}

type GPTSuggester struct {
	client *openai.Client
	GPTSuggesterOptions
}

func (s *GPTSuggester) ChatCompletion(
	ctx context.Context,
	userData UserData,
	msgs []Message,
	prompter ...Prompter,
) ([]string, error) {
	var (
		suggestions = []string{}
		prompt      = ""
		err         error
	)

	switch len(prompter) {
	case 1:
		p := prompter[0]
		err = p.Update(userData, msgs)
		if err != nil {
			break
		}
		prompt, err = p.Build()
	default:
		err = s.prompter.Update(userData, msgs)
		if err != nil {
			break
		}
		prompt, err = s.prompter.Build()
	}

	if err != nil {
		return suggestions, err
	}

	req := openai.CompletionRequest{
		Model:       s.chatModel,
		Prompt:      prompt,
		N:           s.nChat,
		Temperature: s.modelTemperature,
	}
	rsp, err := s.client.CreateCompletion(ctx, req)
	if err != nil {
		return suggestions, err
	}

	if len(rsp.Choices) == 0 {
		return suggestions, errors.New("gptSuggester: got empty reply from server")
	}

	for _, choice := range rsp.Choices {
		suggestions = append(suggestions, choice.Text)
	}

	return suggestions, nil
}

func (s *GPTSuggester) TextCompletion(
	ctx context.Context,
	_ UserData,
	prompt string,
) ([]string, error) {
	req := openai.CompletionRequest{
		Model:       s.textModel,
		Prompt:      prompt,
		N:           s.nText,
		Temperature: s.modelTemperature,
	}
	rsp, err := s.client.CreateCompletion(ctx, req)
	if err != nil {
		return nil, err
	}

	suggestions := []string{}
	if len(rsp.Choices) == 0 {
		return suggestions, errors.New("gptSuggester: got empty reply from server")
	}

	for _, choice := range rsp.Choices {
		suggestions = append(suggestions, choice.Text)
	}
	return suggestions, nil
}

func NewGPTSuggester(client *openai.Client, opts ...Option) (*GPTSuggester, error) {
	gptSuggester := &GPTSuggester{
		client:              client,
		GPTSuggesterOptions: DefaultSuggesterOptions,
	}
	for _, opt := range opts {
		opt(gptSuggester)
	}

	return gptSuggester, nil
}

type Option func(s any)

func optionAdapter(closer func(s *GPTSuggester)) Option {
	return func(i any) {
		switch s := i.(type) {
		case *GPTSuggester:
			closer(s)
		}
	}
}

func WithTemperature(temperature float32) Option {
	return optionAdapter(func(s *GPTSuggester) {
		s.modelTemperature = temperature
	})
}

func WithNText(N int) Option {
	return optionAdapter(func(s *GPTSuggester) {
		s.nText = N
	})
}

func WithNChat(N int) Option {
	return optionAdapter(func(s *GPTSuggester) {
		s.nChat = N
	})
}

func WithTextModel(model string) Option {
	return optionAdapter(func(s *GPTSuggester) {
		s.textModel = model
	})
}

func WithChatModel(model string) Option {
	return optionAdapter(func(s *GPTSuggester) {
		s.chatModel = model
	})
}

func WithPrompter(prompter Prompter) Option {
	return optionAdapter(func(s *GPTSuggester) {
		s.prompter = prompter
	})
}
