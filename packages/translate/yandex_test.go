package translate

import (
	"blinders/packages/utils"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/joho/godotenv"
)

var (
	translator Translator
	skipEnv    = "CI"
)

func init() {
	err := godotenv.Load("../../.env")
	if err != nil {
		fmt.Println("Error loading .env file:", err)
	}

	translator = YandexTranslator{APIKey: os.Getenv("YANDEX_API_KEY")}
}

func TestTranslateWordEN_VI(t *testing.T) {
	utils.SkipTestOnEvironment(t, skipEnv)
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
	utils.SkipTestOnEvironment(t, skipEnv)
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
	utils.SkipTestOnEvironment(t, skipEnv)
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
	utils.SkipTestOnEvironment(t, skipEnv)
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
