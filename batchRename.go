package main

import (
	// "log"
	// "os"
	// "path/filepath"

	"github.com/rivo/tview"
)

func buildBatchRenamePage() *tview.Flex {
	batchRenamePage := tview.NewFlex().SetDirection(tview.FlexRow)

	batchRenameForm := tview.NewForm().
		AddInputField("Find:", " ", 70, nil, nil).
		AddInputField("Replace:", " ", 70, nil, nil)

	batchRenamePage.
		AddItem(batchRenameForm, 0, 1, true)

	return batchRenamePage

}
