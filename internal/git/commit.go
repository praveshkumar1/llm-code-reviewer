package git

import (
	"bytes"
	"os/exec"
	"strings"
)

func LastCommitMessage(repoPath string) (string, error) {
	cmd := exec.Command("git", "-C", repoPath, "log", "-1", "--pretty=%B")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}
