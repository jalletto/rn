package main

import (
	"fmt"
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

	node := app.getCurrentNode()
	fileInfo := node.GetReference().(*fileInfo)
	fullPath := fileInfo.fullPath()
	files, err := os.ReadDir(fullPath)
	batchRenameFileTable := tview.NewTable().SetBorders(true).SetFixed(5, 2)

	if err != nil {
		log.Fatalf("Error reading directory files %v", err)
	}

	i := 0
	for _, file := range files {

		if !file.IsDir() {
			reference := &renameValues{currentName: file.Name(), proposedName: file.Name()}
			batchRenameFileTable.SetCell(i, 0,
				tview.NewTableCell(file.Name()).
					SetReference(reference).
					SetTextColor(tcell.ColorWhite).
					SetExpansion(1))

			batchRenameFileTable.SetCell(i, 1,
				tview.NewTableCell(file.Name()).
					SetReference(reference).
					SetTextColor(tcell.ColorWhite).
					SetExpansion(2))

			i++
		}

	}

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
			app.SwitchToPage("Home")
		})
	batchRenamePage := tview.NewFlex().SetDirection(tview.FlexRow)
	batchRenamePage.
		AddItem(batchRenameForm, 0, 1, true).
		AddItem(batchRenameFileTable, 0, 2, false)

	return batchRenamePage

}

func renameFiles(table *tview.Table, path string) {
	for row := 0; row < table.GetRowCount(); row++ {
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

	for row := 0; row < table.GetRowCount(); row++ {
		cell := table.GetCell(row, col)
		reference := cell.GetReference().(*renameValues)
		currentName := reference.currentName

		if find == "" {
			cell.GetReference().(*renameValues).proposedName = currentName
			cell.SetText(currentName)
		} else if strings.Contains(currentName, find) {

			newFileName := strings.Replace(currentName, find, replace, -1)
			cell.GetReference().(*renameValues).proposedName = newFileName
			cell.SetText(newFileName)

		} else {
			if reference.proposedName != currentName {
				reference.proposedName = currentName
				cell.SetReference(reference)
				cell.SetText(currentName)
			}
		}

	}
}
