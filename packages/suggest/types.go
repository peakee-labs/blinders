package suggest

import (
	"fmt"
	"strings"
)

type Suggestion struct {
	Suggestions    []string
	RequestTokens  int
	ResponseTokens int
	Timestamp      int64 // Unix
}

type UserData struct {
	UserID   string   `json:"userID"`
	Native   Language `json:"nativeLanguage"`
	Learning Language `json:"learningLanguage"`
}

func (d UserData) String() string {
	str := strings.Builder{}
	str.WriteString(fmt.Sprintf("Native language: %s\n", d.Native))
	str.WriteString(fmt.Sprintf("Learning language: %s\n", d.Native))
	return str.String()
}

func GetUserData(userID string) (UserData, error) {
	return UserData{
		UserID: userID,
		Native: Language{
			Lang:  LangVi,
			Level: Intermediate,
		},
		Learning: Language{
			Lang:  LangEn,
			Level: Beginner,
		},
	}, nil
}

const (
	Beginner     Level = "Beginner"
	Intermediate Level = "Intermediate"
	Advanced     Level = "Advanced"
)

var (
	LangVi = Lang{Code: "vi", Name: "Vietnamese"}
	LangEn = Lang{Code: "en", Name: "English"}
)

type (
	Level string
	Lang  struct {
		Code string `json:"languageCode"` // ISO-[639-1] Code of language based
		Name string `json:"languageName"` // English name of language
	}
	Language struct {
		Lang  Lang  `json:"language"`
		Level Level `json:"languageLevel"`
	}
)

func (l Level) String() string {
	return string(l)
}

func (l Lang) String() string {
	return l.Name
}

func (c Language) String() string {
	return fmt.Sprintf("[language: %s, level: %s]", c.Lang, c.Level)
}
