package tui

import (
	"bytes"
	"os"
	"os/exec"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/tvaintrob/tsar/internal/search"
)

func (t *TsarTUI) newFilesList() *tview.List {
	list := tview.NewList()
	list.
		ShowSecondaryText(false).
		SetTitle("Files").
		SetTitleAlign(tview.AlignLeft)

	enableBorderColors(list.Box)
	list.SetChangedFunc(t.onFileSelected)
	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			t.SetFocus(t.searchInput)
			return nil
		case tcell.KeyBacktab:
			t.SetFocus(t.replaceInput)
			return nil
		}
		return event
	})

	return list
}

func (t *TsarTUI) onFileSelected(index int, main string, secondary string, shortcut rune) {
	item := t.matches[index]
	t.renderDiff(item.file, item.matches)
}

func (t *TsarTUI) renderDiff(file string, matches []search.Match) {
	_, _, width, _ := t.output.GetInnerRect()
	t.output.Clear()
	content, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}

	content = matches[0].Pattern.ReplaceAll(content, []byte(t.replaceInput.GetText()))

	cmd := exec.Command("delta", "--width", strconv.Itoa(width), file, "-")
	cmd.Stdin = bytes.NewReader(content)
	cmd.Stdout = tview.ANSIWriter(t.output)
	cmd.Run()

}
