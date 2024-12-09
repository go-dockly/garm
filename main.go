package main

import (
	"flag"

	"github.com/algoboyz/garm/pkg/compile"
	"github.com/algoboyz/garm/pkg/dbg"
)

var (
	debug  bool
	target string
)

func init() {
	flag.BoolVar(&debug, "v", false, "debug mode")
	flag.StringVar(&target, "in", "test/add_simple.go", "src file for compilation")
}

func main() {
	flag.Parse()
	compiler := compile.New(dbg.NewDebugger(debug))

	_, err := compiler.Parse(target, debug)
	if err != nil {
		panic(err)
	}
	_, err = compiler.Generate()
	if err != nil {
		panic(err)
	}
}
