package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func buildHomePage(app *app) *tview.Flex {

	// Set Up containers
	homePageTop := tview.NewFlex()
	menu := tview.NewTextView().
		SetTextColor(tcell.ColorGreen).
		SetText("(r) Rename Current Selection\n(o) Open Folder\n(q) Quit")

	treeView := newTreeView(app.getRoodDir())
	treeView.SetChangedFunc(func(n *tview.TreeNode) {
		app.setCurrentNode(n)
	})

	treeView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		switch event.Rune() {
		case 'r':
			treeView.renderRenameForm(homePageTop, app)
		case 'o':

			if app.getCurrentNode().GetReference().(*fileInfo).isDir {
				app.AddAndSwitchToPage("Batch Rename", buildBatchRenamePage(app), true)

			}
		}
		return event
	})

	// Assemble Layout
	app.SetFocus(treeView)
	homePage := tview.NewFlex().SetDirection(tview.FlexRow)
	homePageTop.AddItem(treeView, 0, 1, true)
	homePage.
		AddItem(homePageTop, 0, 4, true).
		AddItem(menu, 0, 1, false)

	return homePage
}
