package main

import (
	"fmt"
	"os"

	"github.com/tvaintrob/tsar/internal/search"
	"github.com/tvaintrob/tsar/internal/tui"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:            "tsar",
		Usage:           "terminal search and replace",
		Suggest:         true,
		HideHelpCommand: true,
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "dir", Value: ".", Aliases: []string{"d"}, Usage: "directory root to operate on"},
			&cli.StringFlag{Name: "pattern", Aliases: []string{"p"}, Usage: "initial search pattern"},
			&cli.StringFlag{Name: "replace", Aliases: []string{"r"}, Usage: "initial replace content"},
		},
		Action: func(ctx *cli.Context) error {
			root := ctx.String("dir")
			files, err := search.ListFiles(root)
			if err != nil {
				return err
			}

			tuiApp := tui.NewApp(files)
			return tuiApp.Run()
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println()
		fmt.Println(err)
		os.Exit(1)
	}
}
