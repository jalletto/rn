package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type tree struct {
	*tview.TreeView
}

func buildTree(node *tview.TreeNode, dirPath string) error {
	files, err := os.ReadDir(dirPath)

	if err != nil {
		return err
	}

	for _, file := range files {
		fileName := file.Name()

		reference := newFileInfo(file, dirPath)

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

func newRootNode(currentDir string) *tview.TreeNode {

	rootNode := tview.NewTreeNode(currentDir).
		SetSelectable(true).
		SetExpanded(true)

	return rootNode

}

func renameNodeAndFile(node *tview.TreeNode, newName string) {

	oldFileName := node.GetText()
	path := getNodeReference(node).path

	renameFile(oldFileName, newName, path, path)

	node.SetText(newName)

	if len(node.GetChildren()) != 0 {
		reSetAllChildNodes(node)
	}

}

func newTreeView(currentDir string) tree {

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

	t := tree{
		treeView,
	}

	return t

}

func (tree *tree) regenerateTree(currentDir string) {

	rootNode := newRootNode(currentDir)

	if err := buildTree(rootNode, currentDir); err != nil {
		log.Fatalf("Error building tree: %v", err)
	}

	tree.SetRoot(rootNode)
}

func (tree *tree) renderRenameForm(container *tview.Flex, app *app) {
	renameForm := tview.NewForm()
	node := tree.GetCurrentNode()

	newFileName := node.GetText()

	renameForm.AddInputField("Path:", getNodeReference(node).path, 50, nil, nil)

	renameForm.AddInputField("Name:", node.GetText(), 50, nil, func(newName string) {
		newFileName = newName
	})

	renameForm.AddButton("Rename", func() {

		renameNodeAndFile(node, newFileName)

		renameForm.Clear(true)
		container.RemoveItem(renameForm)

		app.SetFocus(tree)

	})

	container.AddItem(renameForm, 0, 1, true)
	app.SetFocus(renameForm)
}
