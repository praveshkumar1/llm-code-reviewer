package git

import (
	"bytes"
	"os/exec"
)

func GetDiff(repoPath string) (string, error) {
	cmd := exec.Command("git", "diff")
	cmd.Dir = repoPath

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return out.String(), nil
}
