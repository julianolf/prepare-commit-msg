package main

import (
	"bufio"
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
	RefineText(string) (string, error)
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

const (
	prefix  = "~/"
	version = "1.0.1"
)

var (
	ai  string
	sys string
	cfg string
	ver bool
)

func init() {
	flag.StringVar(&ai, "ai", "anthropic", "Specifies the AI model to use.")
	flag.StringVar(&sys, "sys", "", "Specifies the system prompt to provide instructions to the AI.")
	flag.StringVar(&cfg, "config", prefix+"prepare-commit-msg.json", "Path to the configuration file.")
	flag.BoolVar(&ver, "version", false, "Show version number and quit.")

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

func readFile(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var content strings.Builder
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}
		content.WriteString(line + "\n")
	}

	err = scanner.Err()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(content.String()), nil
}

func gitDiff() (string, error) {
	var out strings.Builder

	cmd := exec.Command("git", "diff", "--staged")
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(out.String()), nil
}

func main() {
	args := parseArgs()

	if ver {
		fmt.Printf("%s %s\n", os.Args[0], version)
		os.Exit(0)
	}

	var txt string
	var err error

	switch args.Source {
	case "merge", "squash", "commit":
		os.Exit(0)
	case "message":
		txt, err = readFile(args.Filename)
	default:
		txt, err = gitDiff()
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if txt == "" {
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

	var msg string

	switch args.Source {
	case "message":
		msg, err = cli.RefineText(txt)
	default:
		msg, err = cli.CommitMessage(txt)
	}

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
