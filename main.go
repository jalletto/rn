package main

import (
	"log"
	"os"
	"path/filepath"

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
// I want to be able to press 'd' to delete a file - A dialog box should appear to ask if I'm sure.
// I want to be able to rename all the files in a directory at once.
// If I click on a dir in the tree I am taken to new screen where I can see a list of all files
// Here I can change the name of each file.
// I can click a button to batch rename all the files
// I can enter a pattern to select in all files
// I can enter a string to be used as a replacement for the selected pattern

// Bugs

// DONE
// Create a file class to hold all file data. Then update the node reference to use the file class.
// I want to be able to rename the dirs and files from the tree
// I want to see the dirs in the my current root.
// I want to be able to see the files in those dirs

// Bugs Done
// If you rename a directory and then try to rename a file in that directory you get an error.
// Cause: we are relying on the node's reference for the file path. This get's set when we first generate the tree but is not updated after that.
// Possible solution: If we change the reference to a struct that contains all the file data, we can check if we are renaming a dir or not. If we are, then we can get the nodes children and recursively update all child nodes. We could als just completely regenerate the tree. Not sure which is more of a pain in the ass.
