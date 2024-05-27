package commands

import (
	"bytes"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"sync"

	"github.com/fsnotify/fsnotify"
)

func Init(cmd *flag.FlagSet) {
	// Handle init
	rootPath := cmd.String("path", "", "path to the folder we want to sync") // is not working any ideas why ?
	cmd.Parse(os.Args[2:])

	dirFs := os.DirFS(*rootPath)

	// create hidden file with metadata for folder
	file, err := os.Create(*rootPath + ".bobbox")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	folders := []string{*rootPath}
	buff := bytes.NewBuffer([]byte{})
	err = fs.WalkDir(dirFs, ".", func(path string, d fs.DirEntry, err error) error {
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
			_, err = buff.WriteString(fmt.Sprintf("%s %d\n", path, fileInfo.Size()))
			if err != nil {
				fmt.Println(err)
			}
		} else {
			folders = append(folders, *rootPath+path)
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}

	file.Write(buff.Bytes())

	wg := &sync.WaitGroup{}
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

		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			// listen for events
			defer wg.Done()
			for {
				select {
				case event := <-watcher.Events:
					fmt.Println("event:", event)
				case err := <-watcher.Errors:
					fmt.Println("error:", err)
				}
			}
		}(wg)
	}

	wg.Wait()
}
