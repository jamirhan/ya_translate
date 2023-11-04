package ya_translate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type TranslationResponse struct {
	Text                 string `json:"text"`
	DetectedLanguageCode string `json:"detectedLanguageCode"`
}

type TranslationRequest struct {
	FolderID           string   `json:"folderId"`
	Texts              []string `json:"texts"`
	TargetLanguageCode string   `json:"targetLanguageCode"`
}

type Client interface {
	Translate(to string, texts []string) ([]TranslationResponse, error)
}

var _ Client = (*ClientImpl)(nil)

type ClientImpl struct {
	Token    string
	Endpoint string
	FolderID string
}

var DefaultEndpoint = "https://translate.api.cloud.yandex.net"

func (c *ClientImpl) Translate(to string, texts []string) ([]TranslationResponse, error) {
	// I was too lazy to create a normal client so yeah...
	client := http.Client{}
	fullURL, err := url.JoinPath(c.Endpoint, "/translate/v2/translate")
	if err != nil {
		return nil, err
	}

	translationReq := TranslationRequest{
		FolderID:           c.FolderID,
		Texts:              texts,
		TargetLanguageCode: to,
	}
	reqJSON, err := json.Marshal(translationReq)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, fullURL, bytes.NewReader(reqJSON))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to translate (%d): %s", resp.StatusCode, string(body))
	}

	var translated struct {
		Translations []TranslationResponse `json:"translations"`
	}
	if err := json.Unmarshal(body, &translated); err != nil {
		return nil, err
	}

	return translated.Translations, nil
}
