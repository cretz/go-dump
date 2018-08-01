package main

import (
	"flag"
	"log"
	"os"

	"github.com/golang/protobuf/proto"

	"github.com/cretz/go-dump/dump"
	"github.com/golang/protobuf/jsonpb"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	jsonOut := false
	flag.BoolVar(&jsonOut, "json", false, "Output JSON instead of binary")
	flag.Parse()
	pkgs, err := dump.FromArgs(flag.Args(), false)
	if err != nil {
		return err
	}
	if jsonOut {
		marshaler := &jsonpb.Marshaler{Indent: "  "}
		return marshaler.Marshal(os.Stdout, pkgs)
	}
	byts, err := proto.Marshal(pkgs)
	if err == nil {
		_, err = os.Stdout.Write(byts)
	}
	return err
}
