package main

import (
	"flag"
	"fmt"
	"os"
)

var ai string

func init() {
	flag.StringVar(&ai, "ai", "anthropic", "AI to use")

	flag.Usage = func() {
		fmt.Fprintf(os.Stdout, "Usage: %s [options...] [output file] \n", os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()
	fmt.Println("Hello, world!")
}
