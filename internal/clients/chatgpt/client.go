package chatgpt

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"anki-api/internal/entity"
)

const (
	promptID   = "pmpt_6860d2e3a56c8194bd492e375ae2560607cfbd9015b3aff7" // saved prompt in OpenAI
	version    = "5"
	ChatGPTURL = "https://api.openai.com/v1/responses"
)

//curl https://api.openai.com/v1/responses \
//-H "Content-Type: application/json" \
//-H "Authorization: Bearer $OPENAI_API_KEY" \
//-d '{
//"prompt": {
//"id": "pmpt_6860d2e3a56c8194bd492e375ae2560607cfbd9015b3aff7",
//"version": "1",
//"variables": {
//"learning_language": "example learning_language",
//"native_language": "example native_language"
//}
//}
//}'

type Client struct {
	URL        string
	APIKey     string
	HTTPClient *http.Client
}

func New(apiKey string, timeout time.Duration) *Client {
	return &Client{
		URL:    ChatGPTURL,
		APIKey: apiKey,
		HTTPClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *Client) Response(ctx context.Context, message string, learningLang, nativeLang entity.Language) (entity.Flashcard, error) {
	reqBuild := map[string]interface{}{
		"prompt": map[string]interface{}{
			"id":      promptID,
			"version": version,
			"variables": map[string]interface{}{
				"learning_language": learningLang,
				"native_language":   nativeLang,
			},
		},
		"input": []interface{}{map[string]interface{}{
			"role": "user",
			"content": []interface{}{
				map[string]string{
					"type": "input_text",
					"text": message,
				},
			},
		},
		}}

	reqBytes, err := json.Marshal(reqBuild)
	if err != nil {
		return entity.Flashcard{}, fmt.Errorf("json.Marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.URL, bytes.NewBuffer(reqBytes))
	if err != nil {
		return entity.Flashcard{}, fmt.Errorf("http.NewRequestWithContext: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return entity.Flashcard{}, fmt.Errorf("client.Do: %w", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return entity.Flashcard{}, fmt.Errorf("io.ReadAll: %w", err)
	}

	fmt.Println(string(body)) // For debugging

	// Step 1: Parse outer response
	var result struct {
		Output []struct {
			Content []struct {
				Type string `json:"type"`
				Text string `json:"text"`
			} `json:"content"`
		} `json:"output"`
	}

	if err = json.Unmarshal(body, &result); err != nil {
		return entity.Flashcard{}, fmt.Errorf("json.Unmarshal outer: %w", err)
	}

	if len(result.Output) == 0 || len(result.Output[0].Content) == 0 {
		return entity.Flashcard{}, fmt.Errorf("no content in response")
	}

	rawText := result.Output[0].Content[0].Text
	if rawText == "" {
		return entity.Flashcard{}, fmt.Errorf("empty text in response")
	}

	// Step 2: Parse inner JSON string into Flashcard
	var flashcard entity.Flashcard
	if err = json.Unmarshal([]byte(rawText), &flashcard); err != nil {
		return entity.Flashcard{}, fmt.Errorf("json.Unmarshal flashcard: %w", err)
	}

	return flashcard, nil
}
