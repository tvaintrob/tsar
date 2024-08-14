package tui

import (
	"errors"
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/tvaintrob/tsar/internal/search"
	"github.com/tvaintrob/tsar/internal/utils"
)

func (t *TsarTUI) newSearchInput() *tview.InputField {
	input := tview.NewInputField()
	input.
		SetChangedFunc(t.onSearchChange).
		SetFieldBackgroundColor(tcell.ColorDefault).
		SetTitle("Search").
		SetTitleAlign(tview.AlignLeft)

	enableBorderColors(input.Box)
	input.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			t.SetFocus(t.replaceInput)
			return nil
		case tcell.KeyBacktab:
			t.SetFocus(t.filesList)
			return nil
		}
		return event
	})

	return input
}

func (t *TsarTUI) onSearchChange(text string) {
	// we should avoid matching an empty pattern,
	// this will just match anything
	if len(text) == 0 {
		t.filesList.Clear()
		return
	}

	var err error
	matches, err := search.FindMatches(text, t.projectFiles)

	// if the input is not a valid regexp we ignore the error to allow the user to finish writing it
	if err != nil && !errors.Is(err, search.ErrInvalidRegexp) {
		// TODO: better error handling in the components
		// TODO: is probably needed, for now just panic.
		panic(err)
	}

	t.filesList.Clear()

	// TODO: add sorting based on filename and number of matches
	groupedMatches := utils.GroupBy(matches, func(m search.Match) string {
		return m.Filename
	})

	t.matches = make([]fileMatch, 0, len(groupedMatches))
	for file, matches := range groupedMatches {
		itemText := fmt.Sprintf("%s (%d)", file, len(matches))
		t.matches = append(t.matches, fileMatch{file: file, matches: matches})
		t.filesList.AddItem(itemText, "", 0, nil)
	}
}
