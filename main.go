package main

import (
	"fmt"
	"os"
)

func load(args []string) {
	if len(args) != 3 {
		fmt.Fprintf(os.Stderr, "usage: %s %s DIR\n", os.Args[0], os.Args[1])
		os.Exit(1)
	}
	conf := LoadConfig()
	db, err := ConnectMySQL(conf)
	if err != nil {
		panic(err)
	}
	db.PostFiles(args[2])
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: %s COMMAND\n", os.Args[0])
		os.Exit(1)
	}

	if os.Args[1] == "run" {
		Run()
	}

	if os.Args[1] == "load" {
		load(os.Args)
	}

	if os.Args[1] == "setup" {
		conf := LoadConfig()
		SetupMySQL(conf)
	}
}
