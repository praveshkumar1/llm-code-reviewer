package llm

import (
	"fmt"
	"strings"
)

func ExtractResponseText(output map[string]interface{}) (string, error) {
	raw, ok := output["response"]
	if !ok {
		return "", fmt.Errorf("missing 'response' field from LLM output")
	}

	text, ok := raw.(string)
	if !ok {
		return "", fmt.Errorf("'response' is not a string")
	}

	// Remove ```json ``` or ``` wrappers
	text = strings.TrimSpace(text)
	text = strings.TrimPrefix(text, "```json")
	text = strings.TrimPrefix(text, "```")
	text = strings.TrimSuffix(text, "```")
	text = strings.TrimSpace(text)

	return text, nil
}
