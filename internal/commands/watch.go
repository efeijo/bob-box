package commands

import (
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"os/signal"
	"strings"

	"github.com/fsnotify/fsnotify"

	"bobbox/internal/metadata"
)

var hiddenFileName = "/.bobbox.json"

// handle init command
func Init(rootPath, configFilePath *string) {
	dirFs := os.DirFS(*rootPath)

	metaFile, err := metadata.NewMetadataFile(*configFilePath, hiddenFileName)
	if err != nil {
		slog.Error("error creating metadata file", err)
		return
	}

	folders := []string{*rootPath}

	// walk through the root folder and get all folders and files
	slog.Info("walking through root folder")
	err = fs.WalkDir(dirFs, ".", walkDirGettingFolders(*rootPath, &folders, metaFile))
	if err != nil {
		slog.Error("error walking through root folder", err)
	}

	// watch folders for changes
	slog.Info("watching folders")
	err = watchFolders(metaFile, folders)
	if err != nil {
		slog.Error("error watching folders", err)
		return
	}

}

func walkDirGettingFolders(rootPath string, folders *[]string, metadata *metadata.MetadataFile) func(dir string, d fs.DirEntry, err error) error {
	return func(dir string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		fileInfo, err := d.Info()
		if err != nil {
			return err
		}

		if !fileInfo.IsDir() {
			metadata.AddFileOrFolder(dir, rootPath+dir, fileInfo.Size())
			return nil
		} else {
			if dir == "." {
				metadata.AddFileOrFolder(dir, rootPath, -1)
			} else {
				metadata.AddFileOrFolder(dir, rootPath+dir, -1)
			}
			*folders = append(*folders, rootPath+dir)
		}
		return nil
	}
}

func watchFolders(metadata *metadata.MetadataFile, folders []string) error {
	eventsChan := make(chan fsnotify.Event)
	errChan := make(chan error)
	for _, folder := range folders {
		// creates a file watcher
		watchFolder(folder, eventsChan, errChan)
	}

	// wait for events
	err := waitsForEvents(eventsChan, errChan, metadata)
	if err != nil {
		return err
	}

	return metadata.PersistFile()
}

func watchFolder(folder string, eventsChan chan fsnotify.Event, errChan chan error) {
	slog.Info("watching folder", "folder_name", folder)
	go func(chanEvents chan fsnotify.Event, errChan chan error) {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			errChan <- err
			return
		}

		err = watcher.Add(folder)
		if err != nil {
			errChan <- err
			return
		}
		// listen for events on the folder
		for {
			select {
			case event := <-watcher.Events:
				eventsChan <- event
			case err := <-watcher.Errors:
				errChan <- err
				return
			}
		}
	}(eventsChan, errChan)
}

func waitsForEvents(eventsChan chan fsnotify.Event, errChan chan error, metadata *metadata.MetadataFile) error {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	// wait for events

	for {
		select {
		// file system events from each watched folder
		case event := <-eventsChan:
			slog.Info("event received", "event_name", event.Name, "event_op", event.Op)
			handleEvent(event, metadata, eventsChan, errChan)
		// os interrupts
		case <-sigChan:
			return nil
		case err := <-errChan:
			return err
		}
	}

}

func handleEvent(event fsnotify.Event, metadata *metadata.MetadataFile, eventsChan chan fsnotify.Event, errChan chan error) {
	switch event.Op {
	case fsnotify.Create:
		handleCreate(event, metadata, eventsChan, errChan)
	case fsnotify.Rename, fsnotify.Remove:
		handleRenameDeleteMoved(event, metadata)
	case fsnotify.Write, fsnotify.Chmod:
		handleUpdate(event, metadata, eventsChan, errChan)
	}
}

func handleCreate(event fsnotify.Event, metadata *metadata.MetadataFile, eventsChan chan fsnotify.Event, errChan chan error) {
	eventInfo, err := os.Stat(event.Name)
	if err != nil {
		errChan <- fmt.Errorf("error handling creating: %w", err)
	}

	ss := strings.Split(event.Name, "/")

	if eventInfo.IsDir() {
		// watch folder for changes
		watchFolder(event.Name, eventsChan, errChan)
		// add file or folder to metadata
		metadata.AddFileOrFolder(ss[len(ss)-1], event.Name, eventInfo.Size())
	}

	// add file or folder to metadata
	metadata.AddFileOrFolder(ss[len(ss)-1], event.Name, -1)
}

func handleRenameDeleteMoved(event fsnotify.Event, metadata *metadata.MetadataFile) {
	metadata.RemoveFileOrFolder(event.Name)
}

func handleUpdate(event fsnotify.Event, metadata *metadata.MetadataFile, _ chan fsnotify.Event, errChan chan error) {
	eventInfo, err := os.Stat(event.Name)
	if err != nil {
		errChan <- fmt.Errorf("error handling update: %w", err)
	}

	ss := strings.Split(event.Name, "/")

	metadata.AddFileOrFolder(ss[len(ss)-1], event.Name, eventInfo.Size())
}
