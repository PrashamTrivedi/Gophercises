package main

import (
	"fmt"
	"os"
	"path/filepath"
	"taskcmd"
	"taskdb"

	"github.com/mitchellh/go-homedir"
)

func main() {
	home, _ := homedir.Dir()
	dbPath := filepath.Join(home, "tasks.db")
	must(taskdb.Init(dbPath))
	must(taskcmd.RootCmd.Execute())
}

func must(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
