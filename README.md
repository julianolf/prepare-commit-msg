# prepare-commit-msg

> An AI-powered Git commit message generator.

This tool automatically generates commit messages based on the staged changes in your Git repository.

If a message is passed to the git-commit command with the `-m` or `-F` flags, the tool will perform a spelling and grammar check, fixing any issues found and making minor improvements for clarity.

It is designed to be used as a Git hook but can also function as a command-line program. When used from the command line, it can either output the generated message to standard output or save it to a file if an argument is provided.

## AI Support

Currently, this tool supports the [Anthropic](https://docs.anthropic.com) (Claude 3.5) and [OpenAI](https://platform.openai.com) (GPT-4) models.

## Installation

**Requirements:** [Go](https://go.dev) 1.16+.

To install, run:

```sh
go install github.com/julianolf/prepare-commit-msg@latest
```

## Configuration

To access the AI models, you need an API key, which can be generated on the providers' pages:
- Anthropic: [console](https://console.anthropic.com/settings/keys)
- OpenAI: [dashboard](https://platform.openai.com/api-keys)

After obtaining your API key, export it as an environment variable:

```sh
export ANTHROPIC_API_KEY='xyzabc123'
```
or
```sh
export OPENAI_API_KEY='xyzabc123'
```

You can add this to your `.rc` file so it is automatically exported when you start a new shell.

### Configuration File

Configurations can be defined using a JSON file. The following configurations are supported:
- `AI`: Specifies the AI model to use (anthropic or openai).
- `APIKey`: Specifies the API key to use when accessing the AI model.
- `System`: Specifies the system prompts to provide instructions to the AI. There are two different prompts that can be defined:
  - `GenMsg`: Used when generating commit messages from staged changes.
  - `FixMsg`: Used when improving a message provided as input.

```json
{
  "AI": "openai",
  "APIKey": "xyzabc123",
  "System": {
    "GenMsg": "Detailed instructions on how to generate the message",
    "FixMsg": "Detailed instructions on how to improve the message"
  }
}
```

**Notes:**
1. The keys are case-sensitive.
2. All configurations are optional; if a configuration is missing, default values are assumed.

The default location for the configuration file is different depending on the operating system.

On Unix systems, the default location is `$XDG_CONFIG_HOME/prepare-commit-msg/config.json` if the environment variable _XDG_CONFIG_HOME_ is defined, or `$HOME/.config/prepare-commit-msg/config.json` otherwise.

On macOS, it's located at `$HOME/Library/Application Support/prepare-commit-msg/config.json`.

On Windows, it's `%AppData%\prepare-commit-msg\config.json`.

You can explicitly pass a configuration file when running the program by using an argument flag:
```sh
prepare-commit-msg -config=my-config.json
```

### Git Hook

For details on how to configure it as a Git hook, refer to the documentation:
- Git hooks: [prepare-commit-msg](https://git-scm.com/docs/githooks#_prepare_commit_msg)
- Git config: [core.hooksPath](https://git-scm.com/docs/git-config#Documentation/git-config.txt-corehooksPath)

## Usage

For usage, run:

```sh
prepare-commit-msg -h
```
