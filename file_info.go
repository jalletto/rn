package main

import (
	"io/fs"
	"log"
	"time"
)

type fileInfo struct {
	path    string
	name    string
	isDir   bool
	modtime time.Time
	size    int64
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
