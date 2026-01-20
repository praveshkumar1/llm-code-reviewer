package main

import (
	"fmt"
	"os"

	"code-reviewer/internal/reviewer"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: reviewer <repo-path>")
		os.Exit(1)
	}

	if err := reviewer.Run(os.Args[1]); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
