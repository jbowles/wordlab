package wordlab

/*
Functions to traverse directories and get file names.
*/

import (
	"os"
	"path/filepath"
)

// FileHandler contains the directory path, list of file paths, and function to create full file paths.
type FileHandler struct {
	DirName       string
	DirPath       string
	FullFilePaths []string
	FileInfo      []os.FileInfo
	FullPathFn    func(string, string, string) string
}

var separator = string(filepath.Separator)

func NewFileHandler(dirPath string) *FileHandler {
	Log.Debug("Creating new FileHandler for Directory path at %s:", dirPath)
	handler := &FileHandler{
		DirName:    dirPath + separator,
		DirPath:    dirPath,
		FullPathFn: func(dirpath, sep, filename string) string { return dirpath + sep + filename },
	}
	handler.setFileNames()
	return handler
}

func (handle *FileHandler) setFileNames() {
	handle.getFileInfo()
	Log.Debug("number of files %d:", len(handle.FileInfo))
	for _, file := range handle.FileInfo {
		if file.Mode().IsRegular() {
			handle.FullFilePaths = append(
				handle.FullFilePaths,
				handle.FullPathFn(handle.DirPath, separator, file.Name()),
			)
		}
	}
}

func (handle *FileHandler) getFileInfo() {
	Log.Debug("GetFileInfo for new FileHandler %s:", handle.DirPath)
	d, err := os.Open(handle.DirPath)
	if err != nil {
		Log.Critical("%s", err)
		os.Exit(1)
	}
	defer d.Close()

	files, err := d.Readdir(-1)
	if err != nil {
		Log.Critical("%s", err)
		os.Exit(1)
	}
	handle.FileInfo = files
}

func (handle *FileHandler) FileByteSize() map[string]int64 {
	fbs := make(map[string]int64)
	for _, file := range handle.FileInfo {
		if file.Mode().IsRegular() {
			fbs[file.Name()] = file.Size()
		}
	}
	return fbs
}
