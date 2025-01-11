package cmd

import (
	"os"
	"testing"
)

func cleanupTestFile(t *testing.T) {
	files := []string{
		"test_tasks.json",
		"corrupted_tasks.json",
		"nonexistent.json",
	}

	t.Cleanup(func() {
		for _, file := range files {
			os.Remove(file)
		}
	})
}
