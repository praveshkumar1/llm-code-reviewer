package reviewer

func BuildFinalReview(results []ReviewResult) map[string][]Issue {
	byFile := AggregateResults(results)

	for file, issues := range byFile {
		byFile[file] = DeduplicateIssues(issues)
	}

	return byFile
}