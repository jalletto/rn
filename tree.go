package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type TreeView struct {
	*tview.TreeView
}

func buildTree(node *tview.TreeNode, dirPath string) error {
	files, err := os.ReadDir(dirPath)

	if err != nil {
		return err
	}

	for _, file := range files {
		nodeText := file.Name()
		// fullPath := filepath.Join(dirPath, nodeText)
		if file.IsDir() {
			childNode := tview.NewTreeNode(nodeText).
				SetSelectable(true).
				SetExpanded(false).
				SetColor(tcell.ColorBlue).
				SetReference(dirPath)
			node.AddChild(childNode)
			if err := buildTree(childNode, filepath.Join(dirPath, nodeText)); err != nil {
				return err
			}
		} else {
			node.AddChild(tview.NewTreeNode(nodeText).SetReference(dirPath))
		}
	}

	return nil
}

func reSetAllChildNodes(node *tview.TreeNode) error {

	dirPath := filepath.Join(node.GetReference().(string), node.GetText())

	for _, child := range node.GetChildren() {
		if len(child.GetChildren()) != 0 {
			child.SetReference(dirPath)
			if err := reSetAllChildNodes(child); err != nil {
				return err
			}
		} else {
			child.SetReference(dirPath)
		}
	}
	return nil
}

func newRootNode(currentDir string) *tview.TreeNode {

	rootNode := tview.NewTreeNode(currentDir).
		SetSelectable(true).
		SetExpanded(true).
		SetReference(currentDir)

	return rootNode

}

func newTreeView(currentDir string) TreeView {

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

	T := TreeView{
		treeView,
	}

	return T

}

func (tree *TreeView) regenerateTree(currentDir string) {

	rootNode := newRootNode(currentDir)

	if err := buildTree(rootNode, currentDir); err != nil {
		log.Fatalf("Error building tree: %v", err)
	}

	tree.SetRoot(rootNode)

}
