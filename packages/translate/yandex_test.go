package translate

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/joho/godotenv"
)

var translator Translator

func init() {
	err := godotenv.Load("../../.env")
	if err != nil {
		fmt.Println("Error loading .env file", err)
	}

	translator = YandexTranslator{APIKey: os.Getenv("YANDEX_API_KEY")}
}

func TestTranslateWordEN_VI(t *testing.T) {
	text := "absolutely"
	expectedResult := "tuyệt đối"
	fmt.Printf("translate \"%s\" to vietnamese, expect \"%s\"\n", text, expectedResult)

	result, err := translator.Translate(text, EnVi)
	if err != nil {
		t.Error(err)
	}

	if !strings.EqualFold(result, expectedResult) {
		t.Errorf("received = \"%s\", expect \"%s\"", result, expectedResult)
	}
}

func TestTranslateSentenceEN_VI(t *testing.T) {
	text := "hello, My name is Peakee"
	expectedResult := "xin chào, tên tôi là Peakee"
	fmt.Printf("translate \"%s\" to vietnamese, expect \"%s\"\n", text, expectedResult)

	result, err := translator.Translate(text, EnVi)
	if err != nil {
		t.Error(err)
	}

	if !strings.EqualFold(result, expectedResult) {
		t.Errorf("received = \"%s\", expect \"%s\"", result, expectedResult)
	}
}

func TestTranslateWordVI_EN(t *testing.T) {
	text := "tuyệt đối"
	expectedResult := "absolutely"
	fmt.Printf("translate \"%s\" to vietnamese, expect \"%s\"\n", text, expectedResult)

	result, err := translator.Translate(text, ViEn)
	if err != nil {
		t.Error(err)
	}

	if !strings.EqualFold(result, expectedResult) {
		t.Errorf("received = \"%s\", expect \"%s\"", result, expectedResult)
	}
}

func TestTranslateSentenceVI_EN(t *testing.T) {
	text := "xin chào, tên tôi là Peakee"
	expectedResult := "hello, My name is Peakee"
	fmt.Printf("translate \"%s\" to vietnamese, expect \"%s\"\n", text, expectedResult)

	result, err := translator.Translate(text, ViEn)
	if err != nil {
		t.Error(err)
	}

	if !strings.EqualFold(result, expectedResult) {
		t.Errorf("received = \"%s\", expect \"%s\"", result, expectedResult)
	}
}
