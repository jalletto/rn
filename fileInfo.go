package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/rivo/tview"
)

type fileInfo struct {
	path       string
	name       string
	isDir      bool
	modtime    time.Time
	size       int64
	parentNode *tview.TreeNode
}

func (f *fileInfo) fullPath() string {
	return filepath.Join(f.path, f.name)
}

func (f *fileInfo) getParentDir() string {
	return filepath.Base(f.path)
}

func (f *fileInfo) deleteReferenceFile() {

	if f.isDir {
		deleteDir(f.fullPath())
	} else {
		deleteFile(f.fullPath())
	}
}

func newFileInfo(file fs.DirEntry, dirPath string, parentNode *tview.TreeNode) *fileInfo {

	info, err := file.Info()

	if err != nil {
		log.Fatalf("Error getting file info: %v", err)
	}

	return &fileInfo{
		path:       dirPath,
		name:       file.Name(),
		isDir:      file.IsDir(),
		size:       info.Size(),
		modtime:    info.ModTime(),
		parentNode: parentNode,
	}
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

func deleteFile(path string) {

	err := os.Remove(path)
	if err != nil {
		fmt.Println("Error deleting file:", err)
		return
	}

}

func deleteDir(path string) {

	err := os.RemoveAll(path)
	if err != nil {
		fmt.Println("Error deleting directory:", err)
		return
	}

}
