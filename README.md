# GPT-CLI-TOOL

Gpt-cli-tool is a command-line tool that translates pseudo-commands into executable commands across various operating system platforms and command shells. It utilizes the OpenAI API to interpret and transform given commands. 


## Repository

The source code is hosted on GitHub: [https://github.com/hgrimm/gpt-cli-tool](https://github.com/hgrimm/gpt-cli-tool)


## Features

- Transforms pseudo-commands into executable shell commands.
- Supports multiple operating system platforms and command shells.
- Uses OpenAI's powerful language models for command interpretation.
- Easy to deploy due to very low dependencies as the command is a single static binary Golang file.


## Ideas for extensions

- Add a sanitised command history for more context in the GPT prompt
- Better API error handling


## Installation

Before installing `gpt-cli-tool`, make sure you have Go installed on your system.

```bash
git clone https://github.com/hgrimm/gpt-cli-tool.git
cd gpt-cli-tool
go mod init grimm.world/gpt-cli-tool
go mod tidy
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
gpt-cli-tool [options] <pseudo command>
```

## Options

- `-v` Verbose output
- `-m <model>` Specifies the OpenAI GPT model (default: gpt-3.5-turbo) [gpt-3.5-turbo, gpt-4, gpt-4-0613, ...]
- `-V` Display the version of the program

## Examples

```sh
gpt-cli-tool Find files larger than 42kB in the current directory
```

```
Looks insanely complicated? Don't panic. The answer is ...
number of tokens used (total_tokens): 99.0
find . -type f -size +42k
Run? (y/n) 
```

```sh
gpt-cli-tool -v -m gpt-4 "Find files larger than 42kB in the current directory"
```

```
2023/11/04 09:31:30 parent process name zsh
2023/11/04 09:31:30 OPENAI_API_KEY: <key_removed>
2023/11/04 09:31:30 API URL: https://api.openai.com/v1/chat/completions
2023/11/04 09:31:30 plattform: darwin
2023/11/04 09:31:30 
{"model":"gpt-4","messages":[{"role":"user","content":"Convert this pseudo command into a real command that can be run on darwin and zsh command shell. Note that the command might include misspelled, invalid or imagined arguments or even imagined program names. Try your best to convert it into an actual command that would do what the command seems to be intended to do.\n\n  Find files larger than 42kB in the current directory\n\nRespond only with the command."}],"max_tokens":1000}
2023/11/04 09:31:31 result:
map[string]interface {}{"choices":[]interface {}{map[string]interface {}{"finish_reason":"stop", "index":0, "message":map[string]interface {}{"content":"find . -type f -size +42k", "role":"assistant"}}}, "created":1.69908669e+09, "id":"chatcmpl-8H6Ly33o9R0WevdQzXhNOCQd17mTQ", "model":"gpt-4-0613", "object":"chat.completion", "usage":map[string]interface {}{"completion_tokens":10, "prompt_tokens":90, "total_tokens":100}}
Looks insanely complicated? Don't panic. The answer is ...
number of tokens used (total_tokens): 100.0
find . -type f -size +42k
Run? (y/n) 
```

## Acknowledgements

This tool was inspired by the [Unconventional Coding YouTube channel](https://www.youtube.com/@unconv). Special thanks for the idea to create this utility. The inspiration came from their video, which can be seen [here](https://www.youtube.com/watch?v=3LJ30aeT0uY).

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contact

Herwig Grimm - herwig.grimm@gmail.com

