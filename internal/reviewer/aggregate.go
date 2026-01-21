package reviewer

func AggregateResults(results []ReviewResult) map[string][]Issue {
	final := make(map[string][]Issue)

	for _, r := range results {
		final[r.File] = append(final[r.File], r.Issues...)
	}

	return final
}
