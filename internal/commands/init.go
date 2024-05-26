package commands

import (
	"bytes"
	"flag"
	"fmt"
	"io/fs"
	"os"
)

func Init(cmd *flag.FlagSet) {
	// Handle init
	path := cmd.String("path", "", "path to the folder we want to sync") // is not working any ideas why ?
	cmd.Parse(os.Args[2:])

	dirFs := os.DirFS(*path)

	// create hidden file with metadata for folder
	file, err := os.Create(*path + ".bobbox")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

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
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}

	file.Write(buff.Bytes())
}
