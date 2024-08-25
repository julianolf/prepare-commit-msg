package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/julianolf/prepare-commit-msg/ai/anthropic"
)

type AI interface {
	CommitMessage(string) (string, error)
}

var ai string

func init() {
	flag.StringVar(&ai, "ai", "anthropic", "AI to use")

	flag.Usage = func() {
		fmt.Fprintf(os.Stdout, "Usage: %s [options...] [output file] \n", os.Args[0])
		flag.PrintDefaults()
	}
}

func gitDiff() (string, error) {
	var out strings.Builder

	cmd := exec.Command("git", "diff", "--staged")
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return out.String(), nil
}

func main() {
	flag.Parse()

	var cli AI
	switch ai {
	default:
		cli = anthropic.New()
	}

	diff, err := gitDiff()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if strings.TrimSpace(diff) == "" {
		fmt.Println("Nothing to commit")
		os.Exit(0)
	}

	msg, err := cli.CommitMessage(diff)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	var out io.Writer
	args := flag.Args()
	switch len(args) {
	case 1:
		filename := args[0]
		file, err := os.Create(filename)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(3)
		}
		defer file.Close()
		out = file
	default:
		out = os.Stdout
	}

	fmt.Fprintln(out, msg)
}
