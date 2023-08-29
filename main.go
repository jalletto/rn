package main

import (
	"log"
	"os"
	"path/filepath"

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

func renameFile(oldFileName string, newFileName string, oldPath string, newPath string) {

	if oldFileName != newFileName {

		oldFilePath := filepath.Join(oldPath, oldFileName)
		newFilePath := filepath.Join(newPath, newFileName)

		err := os.Rename(oldFilePath, newFilePath)

		if err != nil {
			log.Fatalf("Error renaming file: %v", err)
		}

	}

}

func renameNodeAndFile(node *tview.TreeNode, newName string) {

	oldFileName := node.GetText()
	path := node.GetReference().(string)

	renameFile(oldFileName, newName, path, path)

	node.SetText(newName)

	if len(node.GetChildren()) != 0 {
		reSetAllChildNodes(node)
	}

}

func main() {

	// Create a new tview application
	app := tview.NewApplication().
		EnableMouse(true)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		if event.Rune() == 'q' {
			app.Stop()
		}

		return event
	})

	// Set up Primitives
	currentDir := getWD()

	pages := tview.NewPages()
	flexRow := tview.NewFlex().SetDirection(tview.FlexRow)
	flexCol := tview.NewFlex()
	renameForm := tview.NewForm()
	treeView := newTreeView(currentDir)
	menu := tview.NewTextView().
		SetTextColor(tcell.ColorGreen).
		SetText("(r) To Rename Current Selection\n(q) to quit")

	// debugOutputList := tview.NewList().ShowSecondaryText(false) // For Debugging
	// debugOutputList.AddItem("Debug", " ", 43, nil)

	app.SetFocus(treeView)

	treeView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		if event.Rune() == 'r' { //rename a file or dir
			node := treeView.GetCurrentNode()

			newFileName := node.GetText()

			renameForm.AddInputField("Path:", node.GetReference().(string), 50, nil, nil)

			renameForm.AddInputField("Name:", node.GetText(), 50, nil, func(newName string) {
				newFileName = newName
			})

			renameForm.AddButton("Rename", func() {

				renameNodeAndFile(node, newFileName)

				renameForm.Clear(true)
				flexCol.RemoveItem(renameForm)

				app.SetFocus(treeView)

			})

			flexCol.AddItem(renameForm, 0, 1, true)
			app.SetFocus(renameForm)

		}

		return event
	})

	// Layout
	flexCol.AddItem(treeView, 0, 1, true)
	flexRow.
		AddItem(flexCol, 0, 4, true).
		AddItem(menu, 0, 1, false)
	pages.AddPage("Tree View", flexRow, true, true)

	// Start the application
	if err := app.SetRoot(pages, true).Run(); err != nil {
		log.Fatalf("Error starting application: %v", err)
	}

}

// TODO
// Create a file class to hold all file data. Then update the node reference to use the file class.
// I want to be able to rename all the files in a directory at once.
// If I click on a dir in the tree I am taken to new screen where I can see a list of all files
// Here I can change the name of each file.
// I can click a button to batch rename all the files
// I can enter a pattern to select in all files
// I can enter a string to be used as a replacement for the selected pattern

// Bugs

// DONE
// I want to be able to rename the dirs and files from the tree
// I want to see the dirs in the my current root.
// I want to be able to see the files in those dirs

// Bugs Done
// If you rename a directory and then try to rename a file in that directory you get an error.
// Cause: we are relying on the node's reference for the file path. This get's set when we first generate the tree but is not updated after that.
// Possible solution: If we change the reference to a struct that contains all the file data, we can check if we are renaming a dir or not. If we are, then we can get the nodes children and recursively update all child nodes. We could als just completely regenerate the tree. Not sure which is more of a pain in the ass.
