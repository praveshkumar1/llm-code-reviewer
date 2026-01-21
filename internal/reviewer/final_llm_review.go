package reviewer

import (
	"encoding/json"
	"fmt"

	"code-reviewer/internal/llm"
)

func toPrettyJSON(v any) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "[]"
	}
	return string(b)
}

func FinalLLMReview(file string, issues []Issue) (ReviewResult, error) {
	// If no issues detected, add a default positive comment for this file
	if len(issues) == 0 {
		issues = append(issues, Issue{
			Type:        "style",
			Severity:    "low",
			LineStart:   0,
			LineEnd:     0,
			Description: "No issues detected",
			Suggestion:  "Looks good to me",
		})
	}

	// Optional: call LLM to refine duplicates, merge similar issues, etc.
	prompt := fmt.Sprintf(`
You are a senior software engineer.

Below are detected review issues for file: %s

Tasks:
- Merge duplicates
- Remove weak or repetitive points
- Improve clarity
- Keep only actionable issues
- For any diff block not covered by issues, add a comment: "Looks good to me"

Return ONLY valid JSON in this format:
{
  "file": "%s",
  "issues": [...]
}

Issues:
%s
`, file, file, toPrettyJSON(issues))

	output, err := llm.SendToLLM("qwen2.5-coder:14b", prompt)
	if err != nil {
		return ReviewResult{}, err
	}

	var result ReviewResult
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		return ReviewResult{}, err
	}

	// If LLM returned empty issues (rare), fallback to default positive
	if len(result.Issues) == 0 {
		result.Issues = []Issue{{
			Type:        "style",
			Severity:    "low",
			LineStart:   0,
			LineEnd:     0,
			Description: "No issues detected",
			Suggestion:  "Looks good to me",
		}}
	}

	return result, nil
}
