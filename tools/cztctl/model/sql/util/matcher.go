package util

import (
	"os"
	"path/filepath"
)

// MatchFiles returns the match values by globbing patterns
func MatchFiles(in string) ([]string, error) {
	dir, pattern := filepath.Split(in)
	abs, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}

	files, err := os.ReadDir(abs)
	if err != nil {
		return nil, err
	}
	var res []string
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		name := file.Name()
		match, err := filepath.Match(pattern, name)
		if err != nil {
			return nil, err
		}

		if !match {
			continue
		}

		res = append(res, filepath.Join(abs, name))
	}
	return res, nil
}
