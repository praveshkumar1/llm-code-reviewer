package reviewer

type ReviewResult struct {
	File   string  `json:"file"`
	Issues []Issue `json:"issues"`
}

type Issue struct {
	Type        string `json:"type"`
	Severity    string `json:"severity"`
	LineStart  int    `json:"line_start"`
	LineEnd    int    `json:"line_end"`
	Description string `json:"description"`
	Suggestion  string `json:"suggestion"`
}
