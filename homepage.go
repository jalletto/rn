package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func buildHomePage(app *app) *tview.Flex {

	// Set Up containers
	homePageTop := tview.NewFlex()
	renameFormContainer := tview.NewFlex()
	menu := tview.NewTextView().
		SetTextColor(tcell.ColorGreen).
		SetText("(r) Rename Current Selection\n(d) Delete Selection\n(o) Open Folder\n(q) Quit")

	treeView := newTreeView(app.getRoodDir())

	treeView.SetChangedFunc(func(n *tview.TreeNode) {
		app.setCurrentNode(n)
	})

	treeView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		switch event.Rune() {
		case 'r': // rename selected

			renameForm := renderRenameForm(app, func() {
				renameFormContainer.Clear()
				app.SetFocus(treeView)
			})

			renameFormContainer.AddItem(renameForm, 0, 1, true)
			app.SetFocus(renameForm)

		case 'o': // open dir for batch rename

			if app.getCurrentNode().GetReference().(*fileInfo).isDir {
				app.AddAndSwitchToPage("Batch Rename", buildBatchRenamePage(app), true)
			}
		case 'd': // delete selected
			node := app.getCurrentNode()
			fileInfo := getNodeReference(node)
			path := fileInfo.fullPath()
			parent := getParentNode(node)

			if fileInfo.isDir {
				deleteDir(path)
			} else {
				deleteFile(path)
			}
			treeView.Move(1)
			parent.RemoveChild(node)

		}
		return event
	})

	// Assemble Layout
	app.SetFocus(treeView)
	homePage := tview.NewFlex().SetDirection(tview.FlexRow)
	homePageTop.
		AddItem(treeView, 0, 1, true).
		AddItem(renameFormContainer, 0, 1, true)
	homePage.
		AddItem(homePageTop, 0, 4, true).
		AddItem(menu, 0, 1, false)

	return homePage
}
