package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"time"
)

type fileInfo struct {
	path    string
	name    string
	isDir   bool
	modtime time.Time
	size    int64
}

func (f *fileInfo) fullPath() string {
	return filepath.Join(f.path, f.name)
}

func newFileInfo(file fs.DirEntry, dirPath string) *fileInfo {

	info, err := file.Info()

	if err != nil {
		log.Fatalf("Error getting file info: %v", err)
	}

	return &fileInfo{
		path:    dirPath,
		name:    file.Name(),
		isDir:   file.IsDir(),
		size:    info.Size(),
		modtime: info.ModTime(),
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
