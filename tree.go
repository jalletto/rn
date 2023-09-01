package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Node

func reSetAllChildNodes(node *tview.TreeNode) error {

	dirPath := filepath.Join(getNodeReference(node).path, node.GetText())

	for _, child := range node.GetChildren() {

		reference := getNodeReference(child)

		reference.path = dirPath

		if len(child.GetChildren()) != 0 {
			child.SetReference(reference)
			if err := reSetAllChildNodes(child); err != nil {
				return err
			}
		} else {
			child.SetReference(reference)
		}
	}
	return nil
}

func getNodeReference(node *tview.TreeNode) *fileInfo {
	return node.GetReference().(*fileInfo)
}

func getParentNode(node *tview.TreeNode) *tview.TreeNode {
	return getNodeReference(node).parentNode
}

// TREE

func buildTree(node *tview.TreeNode, dirPath string) error {
	files, err := os.ReadDir(dirPath)

	if err != nil {
		return err
	}

	for _, file := range files {
		fileName := file.Name()

		reference := newFileInfo(file, dirPath, node)

		if file.IsDir() {
			childNode := tview.NewTreeNode(fileName).
				SetSelectable(true).
				SetExpanded(false).
				SetColor(tcell.ColorBlue).
				SetReference(reference)

			node.AddChild(childNode)
			if err := buildTree(childNode, filepath.Join(dirPath, fileName)); err != nil {
				return err
			}
		} else {
			node.AddChild(tview.NewTreeNode(fileName).SetReference(reference))
		}
	}

	return nil
}

func newRootNode(currentDir string) *tview.TreeNode {

	rootNode := tview.NewTreeNode(currentDir).
		SetSelectable(true).
		SetExpanded(true)

	return rootNode

}

func renameNodeAndFile(node *tview.TreeNode, newName string) {

	oldFileName := node.GetText()
	reference := getNodeReference(node)
	path := reference.path

	renameFile(oldFileName, newName, path, path)

	node.SetText(newName)
	reference.name = newName

	if len(node.GetChildren()) != 0 {
		reSetAllChildNodes(node)
	}

}

func newTreeView(currentDir string) *tview.TreeView {

	rootNode := newRootNode(currentDir)

	if err := buildTree(rootNode, currentDir); err != nil {
		log.Fatalf("Error building tree: %v", err)
	}

	treeView := tview.NewTreeView().
		SetRoot(rootNode).
		SetCurrentNode(rootNode)

	treeView.SetSelectedFunc(func(node *tview.TreeNode) {
		node.SetExpanded(!node.IsExpanded())
	})

	return treeView
}

func renderRenameForm(app *app, onSubmitHandler func()) *tview.Form {
	renameForm := tview.NewForm()
	node := app.getCurrentNode()

	newFileName := node.GetText()

	renameForm.AddInputField("Name:", node.GetText(), 50, nil, func(newName string) {
		newFileName = newName
	})

	renameForm.AddButton("Rename", func() {

		renameNodeAndFile(node, newFileName)

		renameForm.Clear(true)

		if onSubmitHandler != nil {

			onSubmitHandler()
		}

	})
	return renameForm

}
