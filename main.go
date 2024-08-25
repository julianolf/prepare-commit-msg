package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/julianolf/prepare-commit-msg/ai/anthropic"
)

type AI interface {
	CommitMessage(string) (string, error)
}

type Args struct {
	Filename string
	Source   string
	SHA      string
}

var ai string

func init() {
	flag.StringVar(&ai, "ai", "anthropic", "AI to use")

	flag.Usage = func() {
		fmt.Fprintf(os.Stdout, "Usage: %s [options...] [output file] [commit source] [commit hash] \n", os.Args[0])
		flag.PrintDefaults()
	}
}

func parseArgs() *Args {
	flag.Parse()
	args := flag.Args()

	switch len(args) {
	case 1:
		return &Args{Filename: args[0]}
	case 2:
		return &Args{Filename: args[0], Source: args[1]}
	case 3:
		return &Args{Filename: args[0], Source: args[1], SHA: args[2]}
	default:
		return &Args{}
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
	args := parseArgs()

	switch args.Source {
	case "message", "merge", "squash", "commit":
		os.Exit(0)
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

	var cli AI
	switch ai {
	default:
		cli = anthropic.New()
	}

	msg, err := cli.CommitMessage(diff)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	out := os.Stdout
	if args.Filename != "" {
		out, err = os.Create(args.Filename)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(3)
		}
		defer out.Close()
	}

	fmt.Fprintln(out, msg)
}
