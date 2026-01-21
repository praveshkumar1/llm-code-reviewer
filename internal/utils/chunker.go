package utils

import "strings"

func ChunkDiff(diff string, maxLines int) []string {
	lines := strings.Split(diff, "\n")
	var chunks []string

	current := ""
	count := 0

	for _, line := range lines {
		if strings.HasPrefix(line, "diff --git") && current != "" {
			chunks = append(chunks, current)
			current = ""
			count = 0
		}

		current += line + "\n"
		count++

		if count >= maxLines {
			chunks = append(chunks, current)
			current = ""
			count = 0
		}
	}

	if current != "" {
		chunks = append(chunks, current)
	}

	return chunks
}

func SplitDiffByFile(diff string, maxLines int) map[string][]string {
	result := make(map[string][]string)

	lines := strings.Split(diff, "\n")
	var currentFile string
	var buffer []string

	flush := func() {
		if currentFile != "" && len(buffer) > 0 {
			result[currentFile] = append(
				result[currentFile],
				strings.Join(buffer, "\n"),
			)
		}
		buffer = nil
	}

	for _, line := range lines {
		if strings.HasPrefix(line, "diff --git") {
			flush()
			parts := strings.Split(line, " ")
			if len(parts) >= 3 {
				currentFile = strings.TrimPrefix(parts[2], "a/")
			}
		}
		buffer = append(buffer, line)
	}

	flush()
	return result
}

