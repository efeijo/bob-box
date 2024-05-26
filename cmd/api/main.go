package main

import (
	"bobbox/internal/commands"
	"flag"
	"fmt"
	"os"
)

var (
	UserID string
)

func main() {
	loginCmd := flag.NewFlagSet("login", flag.ExitOnError)
	initCmd := flag.NewFlagSet("init", flag.ExitOnError)
	syncCmd := flag.NewFlagSet("sync", flag.ExitOnError)
	watchCmd := flag.NewFlagSet("watch", flag.ExitOnError)

	if len(os.Args) < 2 {
		fmt.Println("expected 'login', 'init', 'sync', or 'watch' subcommands")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "login":
		commands.Login(loginCmd)
	case "init":
		commands.Init(initCmd)
	case "sync":
		commands.Sync(syncCmd)
	case "watch":
		commands.Watch(watchCmd)
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}
}

type Storage interface {
	GetFiles()
	PutFiles()
}
