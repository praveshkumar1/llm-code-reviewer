package prompt

func ReviewPrompt(diff string) string {
	return `
You are a senior software engineer.

Review the following git diff.
Focus on:
1. Bugs
2. Security issues
3. Performance
4. Code quality
5. Best practices

Return each issue in this format:
- Severity:
- File:
- Line:
- Issue:
- Suggested Fix:

Git Diff:
` + diff
}
