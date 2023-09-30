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
