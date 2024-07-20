package translate

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const apiTemplate = "https://translate.yandex.net/api/v1.5/tr.json/translate?lang=%s&format=%s&text=%s&key=%s"

type YandexTranslator struct {
	APIKey string
}

func (t YandexTranslator) Translate(text string, langs Languages) (string, error) {
	url := fmt.Sprintf(apiTemplate, string(langs), "plain", url.QueryEscape(text), t.APIKey)

	res, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("error calling api: %v", err)
	}

	translateRes := YandexTranslateRes{}
	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	err = json.Unmarshal(resBytes, &translateRes)
	if err != nil || len(translateRes.Text) == 0 {
		return "", fmt.Errorf("error unmarshalling response body: %v", err)
	}

	return translateRes.Text[0], nil
}

type YandexTranslateRes struct {
	Code int      `json:"code"`
	Lang string   `json:"lang"`
	Text []string `json:"text"`
}
