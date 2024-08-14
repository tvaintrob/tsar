package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (t *TsarTUI) newReplaceInput() *tview.InputField {
	input := tview.NewInputField()
	input.
		SetFieldBackgroundColor(tcell.ColorDefault).
		SetTitle("Replace").
		SetTitleAlign(tview.AlignLeft)

	enableBorderColors(input.Box)
	input.SetChangedFunc(t.onReplaceChange)
	input.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			t.SetFocus(t.filesList)
			return nil
		case tcell.KeyBacktab:
			t.SetFocus(t.searchInput)
			return nil
		}
		return event
	})

	return input
}

func (t *TsarTUI) onReplaceChange(text string) {
	fileIndex := t.filesList.GetCurrentItem()
	item := t.matches[fileIndex]

	t.renderDiff(item.file, item.matches)
}
