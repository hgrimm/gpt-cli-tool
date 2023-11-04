# GPT-CLI-TOOL

Gpt-cli-tool is a command-line tool that translates pseudo-commands into executable commands across various operating system platforms and command shells. It utilizes the OpenAI API to interpret and transform given commands.

## Repository

The source code is hosted on GitHub: [https://github.com/hgrimm/gpt-cli-tool](https://github.com/hgrimm/gpt-cli-tool)

## Features

- Transforms pseudo-commands into executable shell commands.
- Supports multiple operating system platforms and command shells.
- Uses OpenAI's powerful language models for command interpretation.

## Installation

Before installing `gpt-cli-tool`, make sure you have Go installed on your system.

```bash
git clone https://github.com/hgrimm/gpt-cli-tool.git
cd gpt-cli-tool
go build
```

## Setting OpenAI API Key

Set your OpenAI API key as follows:

- Linux and MacOS: `export OPENAI_API_KEY=key`
- Windows PowerShell: `$Env:OPENAI_API_KEY = 'key'`

The API keys can be generated at [here](https://platform.openai.com/account/api-keys).

## Usage

Run the tool from the command line. The environment variable OPENAI_API_KEY must be set with the API key.

```sh
gpt-cli-tool [options] <pseudo command
```

## Options

- `-v` Verbose output
- `-m <model>` Specifies the OpenAI GPT model (default: gpt-3.5-turbo) [gpt-3.5-turbo, gpt-4, gpt-4-0613, ...]
- `-V` Display the version of the program

## Examples

```sh
gpt-cli-tool Find files larger than 42kB in the current directory
```

```sh
go run main.go -v -m gpt-4 "Find files larger than 42kB in the current directory"
```

## Acknowledgements

This tool was inspired by the [Unconventional Coding YouTube channel](https://www.youtube.com/@unconv). Special thanks for the idea to create this utility. The inspiration came from their video, which can be seen [here](https://www.youtube.com/watch?v=3LJ30aeT0uY).

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contact

Herwig Grimm - herwig.grimm@gmail.com

