package search

import (
	"errors"
	"io"
	"os"
	"regexp"
	"strings"
)

var ErrInvalidRegexp = errors.New("failed to compile regexp")

// FindMatches searches the given files and returns all [Matches] for the given pattern
func FindMatches(pattern string, files []string) ([]Match, error) {
	var matches []Match

	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, ErrInvalidRegexp
	}

	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		content, err := io.ReadAll(f)
		if err != nil {
			return nil, err
		}

		lines := strings.Split(string(content), "\n")
		for i, line := range lines {
			matchesInLine := re.FindAllStringIndex(line, -1)
			for _, loc := range matchesInLine {
				context := ""
				if i-1 > 0 {
					context += lines[i-1] + "\n"
				}

				context += line
				if i+1 < len(lines) {
					context += "\n" + lines[i+1]
				}

				match := Match{
					Filename:    file,
					Pattern:     re,
					Context:     context,
					MatchedText: line[loc[0]:loc[1]],
				}
				matches = append(matches, match)
			}
		}
	}

	return matches, nil
}
