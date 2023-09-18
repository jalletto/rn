package main

import (
	"io/fs"
	"log"
	"os"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type renameValues struct {
	currentName  string
	proposedName string
}

func buildBatchRenamePage(app *app) *tview.Flex {

	// Page Vars
	node := app.getCurrentNode()
	fileInfo := node.GetReference().(*fileInfo)
	fullPath := fileInfo.fullPath()
	files, err := os.ReadDir(fullPath)

	// Build Table
	batchRenameFileTable := tview.NewTable().SetBorders(true).SetFixed(5, 2)

	if err != nil {
		log.Fatalf("Error reading directory files %v", err)
	}

	populateTable(batchRenameFileTable, files)

	// Build Form
	replace := ""
	find := ""
	batchRenameForm := tview.NewForm()
	batchRenameForm.
		AddInputField("Find:", "", 70, nil, func(text string) {

			find = text

			findAndReplace(batchRenameFileTable, 1, find, replace)

		}).
		AddInputField("Replace:", "", 70, nil, func(text string) {

			replace = text

			findAndReplace(batchRenameFileTable, 1, find, replace)

		}).
		AddButton("Rename All", func() {

			renameFiles(batchRenameFileTable, fullPath)
		}).
		AddButton("Go Back", func() {

			app.AddAndSwitchToPage("Home", buildHomePage(app), true)

		})

	// Build Layout
	batchRenamePage := tview.NewFlex().SetDirection(tview.FlexRow)
	batchRenamePage.
		AddItem(batchRenameForm, 0, 1, true).
		AddItem(batchRenameFileTable, 0, 2, false)

	return batchRenamePage

}

func populateTable(table *tview.Table, files []fs.DirEntry) {

	table.SetCell(0, 0,
		tview.NewTableCell("Original Name").
			SetTextColor(tcell.ColorYellow).
			SetBackgroundColor(tcell.ColorGrey).
			SetExpansion(1))

	table.SetCell(0, 1,
		tview.NewTableCell("New Name").
			SetTextColor(tcell.ColorYellow).
			SetBackgroundColor(tcell.ColorGrey).
			SetExpansion(1))

	i := 1
	for _, file := range files {

		if !file.IsDir() {
			reference := &renameValues{currentName: file.Name(), proposedName: file.Name()}
			table.SetCell(i, 0,
				tview.NewTableCell(file.Name()).
					SetReference(reference).
					SetTextColor(tcell.ColorWhite).
					SetExpansion(1))

			table.SetCell(i, 1,
				tview.NewTableCell(file.Name()).
					SetReference(reference).
					SetTextColor(tcell.ColorWhite).
					SetExpansion(2))

			i++
		}

	}
}

func renameFiles(table *tview.Table, path string) {
	for row := 1; row < table.GetRowCount(); row++ {
		cell := table.GetCell(row, 1)
		reference := cell.GetReference().(*renameValues)
		currentName := reference.currentName
		proposedName := reference.proposedName

		if currentName != proposedName {
			renameFile(currentName, proposedName, path, path)
			reference.currentName = proposedName
			cell.SetReference(reference)
			table.GetCell(row, 0).SetText(proposedName)
		}

	}
}

func findAndReplace(table *tview.Table, col int, find string, replace string) {

	for row := 1; row < table.GetRowCount(); row++ {
		cell := table.GetCell(row, col)
		reference := cell.GetReference().(*renameValues)
		currentName := reference.currentName

		if find == "" {
			reference.proposedName = currentName
			cell.SetText(currentName)
		} else if strings.Contains(currentName, find) {

			newFileName := strings.Replace(currentName, find, replace, -1)
			reference.proposedName = newFileName
			cell.SetText(newFileName)

		} else if reference.proposedName != currentName {
			reference.proposedName = currentName
			cell.SetText(currentName)
		}

	}
}
