# prepare-commit-msg

> An AI-powered Git commit message generator.

This tool automatically generates commit messages based on the staged changes in your Git repository.

It is designed to be used as a Git hook but can also function as a command-line program. When used from the command line, it can either output the generated message to standard output or save it to a file if an argument is provided.

## AI Support

Currently, this tool only supports the [Anthropic API](https://docs.anthropic.com), but more providers will be added in the future.

## Installation

**Requirements:** [Go](https://go.dev) 1.16+.

To install, run:

```sh
go install github.com/julianolf/prepare-commit-msg@latest
```

## Configuration

To access the Anthropic API, you need an API key, which you can generate in the [developer console](https://console.anthropic.com).

After obtaining your API key, export it as an environment variable:

```sh
export ANTHROPIC_API_KEY='xyzabc123'
```

You can add this to your `.rc` file so it is automatically exported when you start a new shell.

For details on how to configure it as a Git hook, refer to the documentation:
- Git hooks: [prepare-commit-msg](https://git-scm.com/docs/githooks#_prepare_commit_msg)
- Git config: [core.hooksPath](https://git-scm.com/docs/git-config#Documentation/git-config.txt-corehooksPath)
