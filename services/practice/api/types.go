package practiceapi

type SimplePracticeUnit struct {
	Word        string   `json:"word"`
	Explain     string   `json:"explain"`
	ExpandWords []string `json:"expandWords"`
}
