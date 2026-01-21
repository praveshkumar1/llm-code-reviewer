package reviewer

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"code-reviewer/internal/git"
	"code-reviewer/internal/llm"
	"code-reviewer/internal/utils"
)

// Run starts the code review process for a given git repository
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

	fmt.Println(">>> Reviewer started")

	// 🔥 FILE-AWARE CHUNKING
	// Split diff by files, then further split each file if large
	fileChunks := utils.SplitDiffByFile(diff, 3000)
	fmt.Printf("Diff split into %d files\n\n", len(fileChunks))

	// Map[file] → list of all issues from all chunks of that file
	fileIssues := make(map[string][]Issue)

	for file, chunks := range fileChunks {
		fmt.Printf("📄 Reviewing file: %s (%d chunks)\n", file, len(chunks))

		// Combine all chunks for this file into a single string
		combinedDiff := strings.Join(chunks, "\n")

		prompt := fmt.Sprintf(`
You are a strict code reviewer.

Review the following git diff.
Return ONLY valid JSON in this format:

{
  "file": "%s",
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

Diff:
%s
`, file, combinedDiff)

		output, err := llm.SendToLLM("qwen2.5-coder:14b", prompt)
		if err != nil {
			fmt.Println("❌ LLM error:", err)
			continue
		}

		fmt.Println("LLM raw output:", output)

		var result ReviewResult
		if err := json.Unmarshal([]byte(output), &result); err != nil {
			fmt.Println("❌ JSON parse error:", err)
			fmt.Println("Raw:", output)
			continue
		}

		// 🔹 Aggregate all issues per file
		fileIssues[result.File] = append(fileIssues[result.File], result.Issues...)
	}

	// 🔥 FINAL CONSOLIDATION: merge issues for each file
	var finalResults []ReviewResult
	for file, issues := range fileIssues {
		// Optional: run high-level LLM refinement per file
		final, err := FinalLLMReview(file, issues)
		if err != nil {
			fmt.Println("❌ Final LLM review error for", file, ":", err)
			// fallback: include raw issues if refinement fails
			finalResults = append(finalResults, ReviewResult{
				File:   file,
				Issues: issues,
			})
			continue
		}
		finalResults = append(finalResults, final)
	}

	// 🔹 Print final consolidated JSON
	fmt.Println("\n✅ FINAL CONSOLIDATED REVIEW")
	json.NewEncoder(os.Stdout).Encode(finalResults)

	return nil
}
