package reviewer

import (
	"encoding/json"
	"fmt"

	"code-reviewer/internal/git"
	"code-reviewer/internal/llm"
	"code-reviewer/internal/utils"
)

func Run(repoPath string) error {
	if !git.IsGitRepo(repoPath) {
		return fmt.Errorf("not a git repository: %s", repoPath)
	}

	diff, err := git.GetDiff(repoPath)
	if err != nil {
		return err
	}

	if diff == "" {
		fmt.Println("No changes detected")
		return nil
	}

	// chunk diff
	chunks := utils.ChunkDiff(diff, 3000)
	fmt.Printf("Diff split into %d chunks\n\n", len(chunks))

	var allResults []ReviewResult

	for i, chunk := range chunks {
		fmt.Printf("🧩 Processing chunk %d\n", i+1)

		// build prompt for LLM
		prompt := fmt.Sprintf(`You are a code reviewer. Review the following diff and return JSON in this format:
{
  "file": "<filename>",
  "issues": [
    {
      "type": "bug|security|performance|style|best_practice",
      "severity": "low|medium|high|critical",
      "line_start": 12,
      "line_end": 18,
      "description": "what is wrong",
      "suggestion": "how to fix it"
    }
  ]
}

Here is the diff:
%s`, chunk)

		output, err := llm.SendToLLM("qwen2.5-coder:14b", prompt)
		if err != nil {
			fmt.Println("❌ LLM error:", err)
			continue
		}

		// parse LLM JSON output
		var result ReviewResult
		if err := json.Unmarshal([]byte(output), &result); err != nil {
			fmt.Println("❌ Failed to parse JSON from LLM:", err)
			fmt.Println("Raw output:", output)
			continue
		}

		allResults = append(allResults, result)
	}

	// Final JSON summary
	finalJSON, _ := json.MarshalIndent(allResults, "", "  ")
	fmt.Println("✅ Final Review JSON:")
	fmt.Println(string(finalJSON))

	return nil
}

