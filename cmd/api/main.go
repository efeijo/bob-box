package main

import (
	"flag"
	"log/slog"
	"os"

	"bobbox/internal/commands"
)

func main() {
	initCmd := flag.NewFlagSet("watch", flag.ExitOnError)

	if len(os.Args) < 2 {
		slog.Error("expected 'watch'")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "watch":
		rootPath := initCmd.String("path", "", "path to the folder we want to sync")
		configFilePath := initCmd.String("metadata_file", "", "path to store the metadata file")

		initCmd.Parse(os.Args[2:])

		commands.Init(rootPath, configFilePath)
	default:
		flag.PrintDefaults()
		slog.Error("expected 'watch'")
		os.Exit(1)
	}
}
