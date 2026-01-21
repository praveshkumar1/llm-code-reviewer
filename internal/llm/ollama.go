package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// OllamaGenerateURL is the local Ollama endpoint for model inference
const OllamaGenerateURL = "http://localhost:11434/api/generate"

// Request sent to Ollama
type GenerateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

// Top-level response from Ollama
type GenerateResponse struct {
	Model    string `json:"model"`
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

// SendToLLM sends a prompt to Ollama and returns the structured JSON response
func SendToLLM(model, prompt string) (string, error) {
	payload := GenerateRequest{
		Model:  model,
		Prompt: prompt,
		Stream: false,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	client := &http.Client{
		Timeout: 5 * time.Minute,
	}

	req, err := http.NewRequest("POST", OllamaGenerateURL, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var raw map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return "", err
	}

	// ✅ Extract Ollama response text HERE
	response, ok := raw["response"].(string)
	if !ok {
		return "", fmt.Errorf("ollama response missing 'response' field")
	}

	// Clean markdown
	response = strings.TrimSpace(response)
	response = strings.TrimPrefix(response, "```json")
	response = strings.TrimPrefix(response, "```")
	response = strings.TrimSuffix(response, "```")
	response = strings.TrimSpace(response)

	return response, nil
}
