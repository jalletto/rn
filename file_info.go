package main

import (
	"io/fs"
	"log"
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
