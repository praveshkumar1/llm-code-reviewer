package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const OllamaURL = "http://localhost:11434/api/chat"

// ReviewRequest defines what we send to Ollama
type ReviewRequest struct {
	Model  string `json:"model"`  // "qwen2.5-coder:14b"
	Prompt string `json:"prompt"` // chunk content + instructions
}

// ReviewResponse defines what we expect from Ollama
type ReviewResponse struct {
	Output string `json:"output"` // raw LLM output
}

// SendToLLM sends a single chunk to Ollama and returns raw response
func SendToLLM(model string, chunk string) (string, error) {
	reqBody := ReviewRequest{
		Model:  model,
		Prompt: chunk,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(fmt.Sprintf("%s/run", OllamaURL), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var llmResp ReviewResponse
	if err := json.Unmarshal(body, &llmResp); err != nil {
		return "", fmt.Errorf("failed to parse LLM response: %w, raw: %s", err, string(body))
	}

	return llmResp.Output, nil
}
