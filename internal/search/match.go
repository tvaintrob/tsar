package search

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/fatih/color"
)

type Match struct {
	Filename    string
	Pattern     *regexp.Regexp
	MatchedText string
	Context     string
}

func (m Match) String() string {
	highlightedMatch := color.New(color.BgRed).Sprintf("%s", m.MatchedText)
	context := strings.ReplaceAll(m.Context, m.MatchedText, highlightedMatch)

	return fmt.Sprintf("File: %s\n\n%s\n", m.Filename, context)
}
