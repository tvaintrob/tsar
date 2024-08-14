package tui

import (
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
	tview.Borders.HorizontalFocus = tview.BoxDrawingsHeavyHorizontal
	tview.Borders.VerticalFocus = tview.BoxDrawingsHeavyVertical
	tview.Borders.TopLeftFocus = tview.BoxDrawingsHeavyDownAndRight
	tview.Borders.TopRightFocus = tview.BoxDrawingsHeavyDownAndLeft
	tview.Borders.BottomLeftFocus = tview.BoxDrawingsHeavyUpAndRight
	tview.Borders.BottomRightFocus = tview.BoxDrawingsHeavyUpAndLeft
}

func NewApp(files []string) *TsarTUI {
	app := TsarTUI{Application: tview.NewApplication(), projectFiles: files}

	app.output = tview.NewTextView()
	app.output.
		SetDynamicColors(true).
		SetTitle("Output").
		SetBorder(true)

	app.searchInput = app.newSearchInput()
	app.replaceInput = app.newReplaceInput()
	app.filesList = app.newFilesList()

	sidebar := tview.NewFlex().
		SetDirection(tview.FlexColumnCSS).
		AddItem(app.searchInput, 3, 0, true).
		AddItem(app.replaceInput, 3, 0, false).
		AddItem(app.filesList, 0, 1, false)

	layout := tview.NewFlex().
		SetDirection(tview.FlexRowCSS).
		AddItem(sidebar, 65, 0, true).
		AddItem(app.output, 0, 1, false)

	app.SetRoot(layout, true)
	return &app
}
