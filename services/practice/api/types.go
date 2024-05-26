package practiceapi

import (
	"fmt"

	"blinders/packages/db/practicedb"
)

type SimplePracticeUnit struct {
	Word        string   `json:"word"`
	Explain     string   `json:"explain"`
	ExpandWords []string `json:"expandWords"`
}

type RequestFlashCardBody struct {
	FrontText    string `json:"frontText"`
	FrontImgURL  string `json:"frontImgURL"`
	BackText     string `json:"backText"`
	BackImgURL   string `json:"backImgURL"`
	CollectionID string `json:"collectionId"`
}

func (b RequestFlashCardBody) Validate() error {
	if b.FrontText == "" {
		return fmt.Errorf("front text is empty")
	}
	if b.BackText == "" {
		return fmt.Errorf("back text is empty")
	}
	return nil
}

type RequestFlashCardCollection struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Cards       []RequestFlashCardBody `json:"cards"`
}

func (b RequestFlashCardCollection) Validate() error {
	if b.Name == "" {
		return fmt.Errorf("name is empty")
	}

	if b.Cards == nil || len(b.Cards) == 0 {
		return fmt.Errorf("cards are empty")
	}

	return nil
}

type ResponseFlashCardCollection struct {
	Metadata   practicedb.CardCollectionMetadata `json:"metadata"`
	FlashCards []*practicedb.FlashCard           `json:"flashcards"`
}
