package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log-analyzer-bot/internal/model"
	"math"
	"net/http"
	"time"
)

const apiURL = "https://generativelanguage.googleapis.com/v1beta/models/gemini-flash-latest:generateContent"

func AskGemini(apiKey string, prompt string) (string, error) {
	maxRetries := 3

	for attempt := 1; attempt <= maxRetries; attempt++ {
		reqBody := model.GeminiRequest{
			Contents: []model.Content{
				{
					Parts: []model.Part{
						{Text: prompt},
					},
				},
			},
		}
		// setup request
		// serialize request body to JSON
		jsonData, err := json.Marshal(reqBody)
		if err != nil {
			return "", fmt.Errorf("error saat marshaling request body: %v", err)
		}

		req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))

		if err != nil {
			return "", fmt.Errorf("error saat creating request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("x-goog-api-key", apiKey)

		// send request
		resp, err := http.DefaultClient.Do(req)

		if err != nil {
			return "", fmt.Errorf("error saat send request: %v", err)
		}

		respBody, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return "", fmt.Errorf("error saat reading response body: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			if resp.StatusCode >= 500 {
				if attempt < maxRetries {
					delay := time.Duration(math.Pow(2, float64(attempt-1))) * time.Second
					time.Sleep(delay)
					continue
				}
				return "", fmt.Errorf("error saat memanggil Gemini: %s", respBody)
			}
			if resp.StatusCode >= 400 && resp.StatusCode < 500 {
				return "", fmt.Errorf("error saat memanggil Gemini: %s", respBody)
			}
		}

		var geminiResp model.Response
		// unmarshal response body into GeminiResponse struct
		if err := json.Unmarshal(respBody, &geminiResp); err != nil {
			return "", fmt.Errorf("error saat unmarshaling response: %v", err)
		}

		if len(geminiResp.Candidates) > 0 && len(geminiResp.Candidates[0].Content.Parts) > 0 {
			return geminiResp.Candidates[0].Content.Parts[0].Text, nil
		}
	}

	return "", nil
}
