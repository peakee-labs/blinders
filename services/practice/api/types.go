package practiceapi

type SimpleReview struct {
	Word        string   `json:"word"`
	Explain     string   `json:"explain"`
	ExpandWords []string `json:"expandWords"`
}
