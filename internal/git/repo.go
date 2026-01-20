package git

import (
	"os"
	"path/filepath"
)

func IsGitRepo(repoPath string) bool {
	gitDir := filepath.Join(repoPath, ".git")
	info, err := os.Stat(gitDir)
	return err == nil && info.IsDir()
}
