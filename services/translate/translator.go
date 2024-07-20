package translate

type Languages string

const (
	EnVi Languages = "en-vi"
	ViEn Languages = "vi-en"
)

type Translator interface {
	Translate(text string, langs Languages) (string, error)
}
