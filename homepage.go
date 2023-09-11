package main

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func buildHomePage(app *app) *tview.Flex {

	askToDelete := true

	// Set Up containers
	homePageTop := tview.NewFlex()
	renameFormContainer := tview.NewFlex()
	menu := tview.NewTextView().
		SetTextColor(tcell.ColorGreen).
		SetText("(r) Rename Current Selection\n(d) Delete Selection\n(o) Open Folder (Batch Rename)\n(q) Quit")

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
			modal := newAskToDeleteModal(fileInfo.name)
			parent := getParentNode(node)

			modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {

				if strings.Contains(buttonLabel, "Delete") {

					fileInfo.deleteReferenceFile()
					treeView.Move(1)
					parent.RemoveChild(node)
				}

				if buttonLabel == "Delete and Don't Show Again" {
					askToDelete = false
				}

				homePageTop.RemoveItem(modal)
				app.SetFocus(treeView)

			})

			if askToDelete {

				homePageTop.AddItem(modal, 0, 1, true)
				app.SetFocus(modal)
			} else {

				fileInfo.deleteReferenceFile()
				treeView.Move(1)
				parent.RemoveChild(node)

			}
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

func newAskToDeleteModal(thingToDelete string) *tview.Modal {

	labels := []string{
		"Delete",
		"Cancel",
		"Delete and Don't Show Again",
	}

	m := tview.NewModal().
		SetText(fmt.Sprintf("Delete %s ?", thingToDelete)).
		AddButtons(labels)

	return m
}
