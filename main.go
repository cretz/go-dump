package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/cretz/go-dump/dump"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	flag.Parse()
	if flag.NArg() != 1 {
		return fmt.Errorf("Expected dir got no arguments")
	}
	pkgs, err := dump.LoadDir(flag.Arg(0))
	if err != nil {
		return err
	}
	jsonable, err := dump.PackagesToJSONMap(pkgs)
	if err != nil {
		return err
	}
	byts, err := json.MarshalIndent(jsonable, "", "  ")
	if err != nil {
		return err
	}
	println(string(byts))
	return nil
}
