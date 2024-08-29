package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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

const prefix = "~/"

var (
	ai  string
	cfg string
)

func init() {
	flag.StringVar(&ai, "ai", "anthropic", "AI to use")
	flag.StringVar(&cfg, "config", prefix+"prepare-commit-msg.json", "Configuration file")

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

func readConfig() ([]byte, error) {
	if strings.HasPrefix(cfg, prefix) {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		cfg = filepath.Join(home, strings.TrimPrefix(cfg, prefix))
	}

	info, err := os.Stat(cfg)
	if err != nil {
		return nil, nil
	}
	if !info.Mode().IsRegular() {
		return nil, nil
	}

	data, err := os.ReadFile(cfg)
	if err != nil {
		return nil, err
	}

	return data, nil
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

	_, err = readConfig()
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(2)
	}

	var cli AI
	switch ai {
	default:
		cli = anthropic.New()
	}

	msg, err := cli.CommitMessage(diff)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(3)
	}

	out := os.Stdout
	if args.Filename != "" {
		out, err = os.Create(args.Filename)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(4)
		}
		defer out.Close()
	}

	fmt.Fprintln(out, msg)
}
