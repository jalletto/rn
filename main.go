package main

import (
	"log"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type app struct {
	*tview.Application
	*tview.Pages
	currentNode *tview.TreeNode //directory that is currently "open"
	rootDir     string          //root dir; where we are running rn from
}

func (a *app) getCurrentNode() *tview.TreeNode {

	if a.currentNode == nil {
		return nil
	}

	return a.currentNode
}

func (a *app) setCurrentNode(n *tview.TreeNode) {
	a.currentNode = n
}

func (a *app) getRoodDir() string {
	return a.rootDir
}

func getWD() string {

	currentDir, err := os.Getwd()

	if err != nil {
		log.Fatalf("Error getting current directory: %v\n", err)
	}
	return currentDir
}

func main() {
	workingDir := getWD()

	app := &app{
		Application: tview.NewApplication(),
		Pages:       tview.NewPages(),
		currentNode: nil,
		rootDir:     workingDir,
	}

	app.EnableMouse(true)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		if event.Rune() == 'q' {
			app.Stop()
		}

		return event
	})

	app.AddPage("Home", buildHomePage(app), true, true)

	// Start the application
	if err := app.SetRoot(app.Pages, true).Run(); err != nil {
		log.Fatalf("Error starting application: %v", err)
	}

}

// TODO
// I want to be able to press 'm' and get an option to move the file.
// - I can type in a path
// - I can choose from a list of common paths
// - I can create a common path to add to the list
// Add append and prepend all in Batch.
// I want to be able to enter RegEx to select
// I want to be able to Batch delete based on selections
// Add labels to the batch rename table columns
// It would be cool to be able to send the list of file names to chatgpt and have it suggest the changes.

// Found Bugs

// DONE
// I want to see the dirs in my current root.
// I want to be able to see the files in those dirs
// I want to be able to rename the dirs and files one at a time from the tree
// I want to be able to press 'd' to delete a file
// I want to be able to rename all the files in a directory at once.
// - If I click on a dir in the tree I am taken to new screen where I can see a list of all files
// - There is a form where I can set a find pattern and a string to replace it with.
// - I can click a button to batch rename all the files

// Bugs Done
// If you rename a directory and then try to rename a file in that directory you get an error.
// - Cause: we are relying on the node's reference for the file path. This get's set when we first generate the tree but is not updated after that.
// - Possible solution: If we change the reference to a struct that contains all the file data, we can check if we are renaming a dir or not. If we are, then we can get the nodes children and recursively update all child nodes. We could also just completely regenerate the tree. Not sure which is more of a pain in the ass.
