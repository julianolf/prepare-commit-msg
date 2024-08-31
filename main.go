package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/julianolf/prepare-commit-msg/ai/anthropic"
	"github.com/julianolf/prepare-commit-msg/ai/openai"
)

type AI interface {
	CommitMessage(string) (string, error)
}

type Args struct {
	Filename string
	Source   string
	SHA      string
}

type Config struct {
	AI     string
	APIKey string
	System string
}

const prefix = "~/"

var (
	ai  string
	sys string
	cfg string
)

func init() {
	flag.StringVar(&ai, "ai", "anthropic", "Specifies the AI model to use.")
	flag.StringVar(&sys, "sys", "", "Specifies the system prompt to provide instructions to the AI.")
	flag.StringVar(&cfg, "config", prefix+"prepare-commit-msg.json", "Path to the configuration file.")

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

func readConfig() (*Config, error) {
	conf := &Config{AI: ai, System: sys}

	if strings.HasPrefix(cfg, prefix) {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		cfg = filepath.Join(home, strings.TrimPrefix(cfg, prefix))
	}

	info, err := os.Stat(cfg)
	if err != nil {
		return conf, nil
	}
	if !info.Mode().IsRegular() {
		return nil, fmt.Errorf("%s is not a regular file\n", cfg)
	}

	data, err := os.ReadFile(cfg)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, conf)
	if err != nil {
		return nil, err
	}

	return conf, nil
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

	conf, err := readConfig()
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(2)
	}

	var cli AI
	switch conf.AI {
	case "openai":
		cli = openai.New(conf.APIKey, conf.System)
	default:
		cli = anthropic.New(conf.APIKey, conf.System)
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
