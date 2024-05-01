package suggest

import (
	"context"
)

type Prompter interface {
	Build() (string, error)
	Update(...interface{}) error
}

type Suggester interface {
	ChatCompletion(context.Context, UserData, []Message, ...Prompter) ([]string, error)
	TextCompletion(context.Context, UserData, string) ([]string, error)
}
