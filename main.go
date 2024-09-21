package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/julianolf/prepare-commit-msg/ai"
)

const version = "1.0.1"

const (
	ok = iota
	cfgErr
	msgErr
	outErr
)

type Args struct {
	Filename string
	Source   string
	SHA      string
}

var (
	aiModel      string
	systemPrompt string
	configFile   string
	versionFlag  bool
)

func init() {
	flag.StringVar(&aiModel, "ai", "", "Specifies the AI model to use. (default \"anthropic\")")
	flag.StringVar(&systemPrompt, "sys", "", "Specifies the system prompt to provide instructions to the AI.")
	flag.StringVar(&configFile, "config", ai.DefaultConfigFile, "Path to the configuration file.")
	flag.BoolVar(&versionFlag, "version", false, "Show version number and quit.")

	flag.Usage = func() {
		fmt.Fprintf(os.Stdout, "Usage: %s [options...] [output file] [commit source] [commit hash] \n", os.Args[0])
		flag.PrintDefaults()
	}
}

func parseArgs() *Args {
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

func loadConfig() (*ai.Config, error) {
	cfgFlags := &ai.Config{AI: aiModel, System: systemPrompt}
	cfgFile := &ai.Config{}
	cfgEnv := ai.ConfigFromEnv()

	_, err := os.Stat(configFile)
	if err == nil {
		cfgFile, err = ai.ConfigFromFile(configFile)
		if err != nil {
			return nil, err
		}
	}

	cfg := &ai.Config{}
	cfg.Update(cfgEnv)
	cfg.Update(cfgFile)
	cfg.Update(cfgFlags)

	return cfg, nil
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

func readCommitMsg(filename string) (string, error) {
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

func improveCommitMsg(cli ai.AI, args *Args) (string, error) {
	msg, err := readCommitMsg(args.Filename)
	if err != nil {
		return "", err
	}
	if msg == "" {
		return "", nil
	}
	return cli.RefineText(msg)
}

func generateCommitMsg(cli ai.AI) (string, error) {
	diff, err := gitDiff()
	if err != nil {
		return "", err
	}
	if diff == "" {
		return "", nil
	}
	return cli.CommitMessage(diff)
}

func prepareCommitMsg(cli ai.AI, args *Args) (string, error) {
	switch args.Source {
	case "merge", "squash", "commit":
		return "", nil
	case "message":
		return improveCommitMsg(cli, args)
	default:
		return generateCommitMsg(cli)
	}
}

func outputCommitMsg(msg string, args *Args) error {
	out := os.Stdout
	if args.Filename != "" {
		var err error
		out, err = os.Create(args.Filename)
		if err != nil {
			return err
		}
		defer out.Close()
	}
	fmt.Fprintln(out, msg)
	return nil
}

func showVersion() {
	fmt.Printf("%s %s\n", os.Args[0], version)
}

func main() {
	flag.Parse()
	if versionFlag {
		showVersion()
		os.Exit(ok)
	}

	args := parseArgs()
	conf, err := loadConfig()
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(cfgErr)
	}

	cli := ai.New(conf)
	msg, err := prepareCommitMsg(cli, args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(msgErr)
	}
	if msg == "" {
		os.Exit(ok)
	}

	err = outputCommitMsg(msg, args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(outErr)
	}
}
