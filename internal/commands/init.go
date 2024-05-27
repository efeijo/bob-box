package commands

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"sync"

	"github.com/fsnotify/fsnotify"
)

var hiddenFileName = ".bobbox"

// handle init command
func Init(rootPath *string) {
	dirFs := os.DirFS(*rootPath)

	// create hidden file with metadata for folder
	file, err := os.Create(*rootPath + hiddenFileName)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	metadataFile := NewMetadataFile(file)

	folders := []string{*rootPath}

	// walk through the root folder and get all folders and files
	err = fs.WalkDir(dirFs, ".", walkDirGettingFolders(*rootPath, &folders, metadataFile))
	if err != nil {
		fmt.Println(err)
	}

	metadataFile.PersistFile()

	watchFolders(metadataFile, folders)

}

func walkDirGettingFolders(rootPath string, folders *[]string, metadata *MetadataFile) func(dir string, d fs.DirEntry, err error) error {
	return func(dir string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}
		fileInfo, err := d.Info()
		if err != nil {
			fmt.Println(err)
			return err
		}

		if !fileInfo.IsDir() {
			metadata.AddFile(dir, fileInfo.Size())
			if err != nil {
				fmt.Println(err)
			}
		} else {
			*folders = append(*folders, rootPath+dir)
		}
		return nil
	}
}

func watchFolders(metadata *MetadataFile, folders []string) {

	eventsChan := make(chan fsnotify.Event)
	for _, folder := range folders {
		// creates a file watcher
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			fmt.Println(err)
		}

		err = watcher.Add(folder)
		if err != nil {
			fmt.Println(err)
		}

		go func(chanEvents chan fsnotify.Event) {
			// listen for events
			for {
				select {
				case event := <-watcher.Events:
					eventsChan <- event
				case err := <-watcher.Errors:
					fmt.Println("error:", err)
				}
			}
		}(eventsChan)
	}

	for event := range eventsChan {
		switch event.Op {
		case fsnotify.Create:
			handleCreate(event.Name)
			size, err := os.ReadFile(event.Name)
			if err != nil {
				fmt.Println(err)
			}
			metadata.AddFile(event.Name, int64(len(size)))
		case fsnotify.Rename, fsnotify.Remove:
			handleRenameDeleteMoved(event.Name)
			metadata.RemoveFile(event.Name)
		case fsnotify.Write, fsnotify.Chmod:
			handleUpdate(event.Name)
			size, err := os.ReadFile(event.Name)
			if err != nil {
				fmt.Println(err)
			}
			metadata.AddFile(event.Name, int64(len(size)))
		}
		fmt.Println("event", event.String())
		// FIXME: persist metadata causes a lot of writes to the file so updates are triggered
		metadata.PersistFile()
	}
}

func handleCreate(pathToFile string) {
	fmt.Println("create", pathToFile)
}

func handleRenameDeleteMoved(pathToFile string) {
	fmt.Println("renamed", pathToFile)
}

func handleUpdate(pathToFile string) {
	fmt.Println("update", pathToFile)
}

type MetadataFile struct {
	mu    sync.Mutex
	Files map[string]int64 `json:"files"`
	file  *os.File
}

func NewMetadataFile(file *os.File) *MetadataFile {
	return &MetadataFile{
		Files: make(map[string]int64),
		file:  file,
	}
}

func (m *MetadataFile) PersistFile() error {
	jsonBytes, err := json.Marshal(m.Files)
	if err != nil {
		return err
	}
	_, err = m.file.Write(jsonBytes)
	return err
}

func (m *MetadataFile) AddFile(path string, size int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Files[path] = size
}

func (m *MetadataFile) Close() error {
	return m.file.Close()
}

func (m *MetadataFile) RemoveFile(path string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.Files, path)
}
