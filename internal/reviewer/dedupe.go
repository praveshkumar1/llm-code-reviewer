package reviewer

import "fmt"

func DeduplicateIssues(issues []Issue) []Issue {
	seen := make(map[string]bool)
	var unique []Issue

	for _, issue := range issues {
		key := fmt.Sprintf(
			"%s-%d-%d-%s",
			issue.Type,
			issue.LineStart,
			issue.LineEnd,
			issue.Description,
		)

		if !seen[key] {
			seen[key] = true
			unique = append(unique, issue)
		}
	}

	return unique
}
