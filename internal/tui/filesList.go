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

var hasDelta bool

func init() {
	_, err := exec.LookPath("delta")
	hasDelta = err == nil
}

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

	list.SetSelectedFunc(func(index int, _, _ string, _ rune) {
		fileMatch := t.matches[index]
		content, err := os.ReadFile(fileMatch.file)
		if err != nil {
			panic(err)
		}

		match := fileMatch.matches[0]
		content = match.Pattern.ReplaceAll(content, []byte(t.replaceInput.GetText()))

		finfo, err := os.Stat(fileMatch.file)
		if err != nil {
			panic(err)
		}

		if err := os.WriteFile(fileMatch.file, content, finfo.Mode()); err != nil {
			panic(err)
		}

		// refresh the list
		t.output.Clear()
		t.onSearchChange(t.searchInput.GetText())
	})

	return list
}

func (t *TsarTUI) onFileSelected(index int, main string, secondary string, shortcut rune) {
	item := t.matches[index]
	t.renderDiff(item.file, item.matches)
}

func (t *TsarTUI) renderDiff(file string, matches []search.Match) {
	t.output.Clear()
	content, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}

	content = matches[0].Pattern.ReplaceAll(content, []byte(t.replaceInput.GetText()))
	cmd := exec.Command("diff", "--color=always", "-u", file, "-")
	if hasDelta {
		_, _, width, _ := t.output.GetInnerRect()
		cmd = exec.Command("delta", "--width", strconv.Itoa(width), file, "-")
	}

	cmd.Stdin = bytes.NewReader(content)
	cmd.Stdout = tview.ANSIWriter(t.output)
	_ = cmd.Run()
}
