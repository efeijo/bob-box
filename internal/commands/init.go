package commands

import (
	"fmt"
	"io/fs"
	"os"
	"os/signal"

	"github.com/fsnotify/fsnotify"

	"bobbox/internal/metadata"
)

var hiddenFileName = "/.bobbox.json"

// handle init command
func Init(rootPath, configFilePath *string) {
	dirFs := os.DirFS(*rootPath)

	// create metadata file FIXME:
	metaFile, err := metadata.NewMetadataFile(*configFilePath, hiddenFileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer metaFile.Close()

	folders := []string{*rootPath}

	// walk through the root folder and get all folders and files
	err = fs.WalkDir(dirFs, ".", walkDirGettingFolders(*rootPath, &folders, nil))
	if err != nil {
		fmt.Println(err)
	}

	// watch folders for changes
	watchFolders(nil, folders)

}

func walkDirGettingFolders(rootPath string, folders *[]string, metadata *metadata.MetadataFile) func(dir string, d fs.DirEntry, err error) error {
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
			return nil
		} else {
			*folders = append(*folders, rootPath+dir)
		}
		return nil
	}
}

func watchFolders(metadata *metadata.MetadataFile, folders []string) {
	eventsChan := make(chan fsnotify.Event)
	for _, folder := range folders {
		// creates a file watcher
		watchFolder(folder, eventsChan)
	}

	// wait for events
	waitsForEvents(eventsChan, metadata)

}

func watchFolder(folder string, eventsChan chan fsnotify.Event) {
	go func(chanEvents chan fsnotify.Event) {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			fmt.Println(err)
			return
		}

		err = watcher.Add(folder)
		if err != nil {
			fmt.Println(err)
			return
		}
		// listen for events on the folder
		for {
			select {
			case event := <-watcher.Events:
				eventsChan <- event
			case err := <-watcher.Errors:
				fmt.Println("error:", err)
				return
			}
		}
	}(eventsChan)
}

func waitsForEvents(eventsChan chan fsnotify.Event, metadata *metadata.MetadataFile) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	// wait for events
	for {
		select {
		// file system events from each watched folder
		case event := <-eventsChan:
			//fmt.Println("event:", event)
			handleEvent(event, metadata)
		// os interrupts
		case <-sigChan:
			//fmt.Println("SIGINT received. Exiting...")
			return
		}
	}

}

func handleEvent(event fsnotify.Event, metadata *metadata.MetadataFile) error {
	switch event.Op {
	case fsnotify.Create:
		handleCreate(event)
	case fsnotify.Rename, fsnotify.Remove:
		handleRenameDeleteMoved(event)
	case fsnotify.Write, fsnotify.Chmod:
		handleUpdate(event)
	}

	return nil
}

func handleCreate(event fsnotify.Event) {
	//fmt.Println("create", event.Name, event.Op)
}

func handleRenameDeleteMoved(event fsnotify.Event) {
	//fmt.Println("renamed", event.Name, event.Op)
}

func handleUpdate(event fsnotify.Event) {
	//fmt.Println("update", event.Name, event.Op)
}
