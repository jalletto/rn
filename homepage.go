package main

import (
	// "log"
	// "os"
	// "path/filepath"

	"log"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func getWD() string {

	currentDir, err := os.Getwd()

	if err != nil {
		log.Fatalf("Error getting current directory: %v\n", err)
	}
	return currentDir
}

func buildHomePage(app *app) *tview.Flex {

	// Set Up containers
	homePage := tview.NewFlex().SetDirection(tview.FlexRow)
	homePageTop := tview.NewFlex()
	menu := tview.NewTextView().
		SetTextColor(tcell.ColorGreen).
		SetText("(r) To Rename Current Selection\n(q) to quit")

	treeView := newTreeView(getWD())
	treeView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		switch event.Rune() {
		case 'r':
			treeView.renderRenameForm(homePageTop, app)
		case 'o':
			app.SwitchToPage("Batch Rename")
		}
		return event
	})

	// Assemble Layout
	app.SetFocus(treeView)
	homePageTop.AddItem(treeView, 0, 1, true)
	homePage.
		AddItem(homePageTop, 0, 4, true).
		AddItem(menu, 0, 1, false)

	return homePage
}
