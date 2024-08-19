package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/tvaintrob/tsar/internal/search"
)

type fileMatch struct {
	file    string
	matches []search.Match
}

type TsarTUI struct {
	*tview.Application
	output       *tview.TextView
	searchInput  *tview.InputField
	replaceInput *tview.InputField
	filesList    *tview.List

	projectFiles []string
	matches      []fileMatch
}

func init() {
	tview.Borders.VerticalFocus = tview.BoxDrawingsHeavyVertical
	tview.Borders.HorizontalFocus = tview.BoxDrawingsHeavyHorizontal
	tview.Borders.TopLeftFocus = tview.BoxDrawingsHeavyDownAndRight
	tview.Borders.TopRightFocus = tview.BoxDrawingsHeavyDownAndLeft
	tview.Borders.BottomLeftFocus = tview.BoxDrawingsHeavyUpAndRight
	tview.Borders.BottomRightFocus = tview.BoxDrawingsHeavyUpAndLeft
}

func enableBorderColors(b *tview.Box) {
	b.SetBorder(true).
		SetBorderColor(tcell.ColorWhite).
		SetFocusFunc(func() { b.SetBorderColor(tcell.ColorYellow) }).
		SetBlurFunc(func() { b.SetBorderColor(tcell.ColorWhite) })
}

func NewApp(files []string, pattern, replace string) *TsarTUI {
	app := TsarTUI{Application: tview.NewApplication(), projectFiles: files}

	app.output = tview.NewTextView()
	app.output.
		SetDynamicColors(true).
		SetTitle("Preview").
		SetBorder(true)

	app.searchInput = app.newSearchInput(pattern)
	app.replaceInput = app.newReplaceInput(replace)
	app.filesList = app.newFilesList()

	helpText := tview.
		NewTextView().
		SetDynamicColors(true).
		SetText(`[#626262]esc: [#4A4A4A]Quit[-]    [#626262]tab: [#4A4A4A]Next Panel[-]    [#626262]shift+tab: [#4A4A4A]Previous Panel[-]    [#626262]space / enter: [#4A4A4A]Confirm Replacement[-]`)

	sidebar := tview.NewFlex().
		SetDirection(tview.FlexColumnCSS).
		AddItem(app.searchInput, 3, 0, true).
		AddItem(app.replaceInput, 3, 0, false).
		AddItem(app.filesList, 0, 1, false)

	mainLayout := tview.NewFlex().
		SetDirection(tview.FlexRowCSS).
		AddItem(sidebar, 88, 0, true).
		AddItem(app.output, 0, 1, false)

	layout := tview.NewFlex().
		SetDirection(tview.FlexColumnCSS).
		AddItem(mainLayout, 0, 1, true).
		AddItem(helpText, 1, 0, false)

	app.SetRoot(layout, true)

	if len(pattern) > 0 {
		app.onSearchChange(pattern)
	}

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			app.Stop()
			return nil
		}
		return event
	})

	return &app
}
