package metadata

import (
	"encoding/json"
	"log/slog"
	"os"
	"sync"
)

type MetadataFile struct {
	mu           sync.Mutex
	file         string
	currentState map[string]FileInfo
}

type FileType string

const (
	File FileType = "file"
	Dir  FileType = "dir"
)

type FileInfo struct {
	Name     string   `json:"name"`
	FullPath string   `json:"full_path"`
	Size     int64    `json:"size"`
	FileType FileType `json:"file_type"`
}

func NewMetadataFile(configFilePath, hiddenFileName string) (*MetadataFile, error) {
	return &MetadataFile{
		currentState: make(map[string]FileInfo),
		file:         configFilePath + hiddenFileName,
	}, nil
}

func (m *MetadataFile) PersistFile() error {
	file, err := os.Create(m.file)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			slog.Error("error closing file", cerr)
		}
	}()

	jsonBytes, err := json.Marshal(m.currentState)
	if err != nil {
		return err
	}
	_, err = file.Write(jsonBytes)
	return err
}

func (m *MetadataFile) AddFileOrFolder(path string, fullPath string, size int64) {
	fileInfo := FileInfo{
		Name:     path,
		FullPath: fullPath,
		Size:     size,
		FileType: File,
	}
	if size == -1 {
		fileInfo.FileType = Dir
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	m.currentState[fullPath] = fileInfo
}

func (m *MetadataFile) RemoveFileOrFolder(path string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.currentState, path)
}
