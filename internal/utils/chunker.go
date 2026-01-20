package utils

const DefaultChunkSize = 3000 // characters (safe for 14B models)

func ChunkDiff(diff string, maxSize int) []string {
	if maxSize <= 0 {
		maxSize = DefaultChunkSize
	}

	var chunks []string
	var current string

	parts := splitByFile(diff)

	for _, part := range parts {
		if len(current)+len(part) > maxSize {
			if current != "" {
				chunks = append(chunks, current)
			}
			current = part
		} else {
			current += part
		}
	}

	if current != "" {
		chunks = append(chunks, current)
	}

	return chunks
}

func splitByFile(diff string) []string {
	var files []string
	current := ""

	lines := splitLines(diff)

	for _, line := range lines {
		if startsWith(line, "diff --git") && current != "" {
			files = append(files, current)
			current = ""
		}
		current += line + "\n"
	}

	if current != "" {
		files = append(files, current)
	}

	return files
}

/* helpers */

func splitLines(s string) []string {
	var lines []string
	start := 0

	for i := range s {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}

	if start < len(s) {
		lines = append(lines, s[start:])
	}

	return lines
}

func startsWith(s, prefix string) bool {
	if len(s) < len(prefix) {
		return false
	}
	return s[:len(prefix)] == prefix
}
