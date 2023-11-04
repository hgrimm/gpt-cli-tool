/*
Gpt-cli-tool: Command line utility for translating pseudo-commands into executable shell commands.

Licensed under the MIT License. See LICENSE file in the repository root for full license text.
https://github.com/hgrimm/gpt-cli-tool

Author: Herwig Grimm <herwig.grimm@gmail.com>
Version: 0.1
Date: November 3, 2023

Usage: gpt-cli-tool [options] <pseudo command>

	This tool takes a 'pseudo command' as input and attempts to translate it
	into an executable command for various operating system platforms and command shells,
	leveraging the OpenAI API for command transformation and execution.

Options:

	-v              Enable verbose output.
	-m <model>      Specify the OpenAI model to use for command translation (default: gpt-3.5-turbo).
	-V              Display the version of gpt-cli-tool.

Example:

	gpt-cli-tool -v -m gpt-4 "show me all files with extension .txt"

The above example will output the translated command in verbose mode using the gpt-4 model.
*/
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/fatih/color"
	ps "github.com/mitchellh/go-ps"
)

var (
	verbose        bool
	displayVersion bool
	model          string
)

const (
	version    = "0.1"
	apiKeyInfo = "Goto https://platform.openai.com/account/api-keys to get your API key. Set the API key on CLI by 'export OPENAI_API_KEY=key' on Linux and MacOS or $Env:OPENAI_API_KEY = 'key' on Windows PowerShell\n\n"
)

func init() {
	flag.BoolVar(&verbose, "v", false, "verbose")
	flag.StringVar(&model, "m", "gpt-3.5-turbo", "OpenAI model (gpt-3.5-turbo, gpt-4, ...)\nFor further information, refer to https://platform.openai.com/docs/models/overview")
	flag.BoolVar(&displayVersion, "V", false, "display version")
}

func debugPrintf(format string, args ...interface{}) {
	if verbose {
		log.Printf(format, args...)
	}
}

type ChatCompletionRequest struct {
	Model            string                  `json:"model"`
	Messages         []ChatCompletionMessage `json:"messages"`
	MaxTokens        int                     `json:"max_tokens,omitempty"`
	Temperature      float32                 `json:"temperature,omitempty"`
	TopP             float32                 `json:"top_p,omitempty"`
	N                int                     `json:"n,omitempty"`
	Stream           bool                    `json:"stream,omitempty"`
	Stop             []string                `json:"stop,omitempty"`
	PresencePenalty  float32                 `json:"presence_penalty,omitempty"`
	FrequencyPenalty float32                 `json:"frequency_penalty,omitempty"`
	// LogitBias is must be a token id string (specified by their token ID in the tokenizer), not a word string.
	// incorrect: `"logit_bias":{"You": 6}`, correct: `"logit_bias":{"1639": 6}`
	// refs: https://platform.openai.com/docs/api-reference/chat/create#chat/create-logit_bias
	LogitBias    map[string]int       `json:"logit_bias,omitempty"`
	User         string               `json:"user,omitempty"`
	Functions    []FunctionDefinition `json:"functions,omitempty"`
	FunctionCall any                  `json:"function_call,omitempty"`
}

type ChatCompletionMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`

	// This property isn't in the official documentation, but it's in
	// the documentation for the official library for python:
	// - https://github.com/openai/openai-python/blob/main/chatml.md
	// - https://github.com/openai/openai-cookbook/blob/main/examples/How_to_count_tokens_with_tiktoken.ipynb
	Name string `json:"name,omitempty"`

	FunctionCall *FunctionCall `json:"function_call,omitempty"`
}

type FunctionCall struct {
	Name string `json:"name,omitempty"`
	// call function with arguments in JSON format
	Arguments string `json:"arguments,omitempty"`
}

type FunctionDefinition struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	// Parameters is an object describing the function.
	// You can pass json.RawMessage to describe the schema,
	// or you can pass in a struct which serializes to the proper JSON schema.
	// The jsonschema package is provided for convenience, but you should
	// consider another specialized library if you require more complex schemas.
	Parameters any `json:"parameters"`
}

func makeCommand(pseudoCommand, commandShell string) (string, map[string]interface{}) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	// check if OPENAI_API_KEY is set
	if apiKey == "" {
		fmt.Printf("Error: OPENAI_API_KEY is not set. %s", apiKeyInfo)
		os.Exit(1)
	}
	debugPrintf("OPENAI_API_KEY: %s", apiKey)

	url := "https://api.openai.com/v1/chat/completions"
	debugPrintf("API URL: %s", url)

	plattform := runtime.GOOS
	debugPrintf("plattform: %s", plattform)

	payload := ChatCompletionRequest{
		Model: model,
		Messages: []ChatCompletionMessage{
			{
				Role: "user",
				Content: "Convert this pseudo command into a real command that can be run on " +
					plattform + " and " + commandShell + " command shell. Note that the command might include misspelled, invalid or " +
					"imagined arguments or even imagined program names. Try your best to convert it " +
					"into an actual command that would do what the command seems to be intended to do.\n\n" +
					pseudoCommand + "\n\nRespond only with the command.",
			},
		},
		MaxTokens: 1000,
	}

	data, err := json.Marshal(payload)
	debugPrintf("\n%s\n", data)
	if err != nil {
		fmt.Println("Error marshaling payload:", err)
		os.Exit(1)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		fmt.Println("Error creating request:", err)
		os.Exit(1)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("Error reading response:", err)
		os.Exit(1)
	}

	var result map[string]interface{}
	json.Unmarshal(body, &result)
	// check if error is returned in result map
	if _, ok := result["error"]; ok {
		errorDetails := result["error"].(map[string]interface{})
		fmt.Printf("Error code %s: %s", errorDetails["code"], errorDetails["message"])
		os.Exit(1)
	}
	debugPrintf("result:\n%#v\n", result)
	choices := result["choices"].([]interface{})
	firstChoice := choices[0].(map[string]interface{})
	message := firstChoice["message"].(map[string]interface{})
	content := message["content"].(string)
	usage := result["usage"].(map[string]interface{})
	debugPrintf("usage:\n%#v\n", usage)

	return content, usage
}

func main() {

	flag.Usage = func() {
		thisCommand := filepath.Base(os.Args[0])

		parentProcessId := os.Getppid()
		parentProcess, _ := ps.FindProcess(parentProcessId)
		parentProcessName := parentProcess.Executable()

		fmt.Fprintf(os.Stderr, "Convert a pseudo command into a real command that can be run on "+runtime.GOOS+" and "+parentProcessName+" command shell.\n\n")
		fmt.Fprintf(os.Stderr, "%s-%s by Herwig Grimm <herwig.grimm@gmail.com>\n\n", thisCommand, version)
		fmt.Fprintf(os.Stderr, "Usage: "+thisCommand+" [-v] [-m <model>] <pseudo command>\n\n")
		fmt.Fprintf(os.Stderr, "Command requires API key from OpenAI. "+apiKeyInfo)
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if displayVersion {
		fmt.Println("Version: " + version)
		os.Exit(0)
	}

	parentProcessId := os.Getppid()
	parentProcess, _ := ps.FindProcess(parentProcessId)
	parentProcessName := parentProcess.Executable()
	debugPrintf("parent process name %s\n", parentProcessName)

	args := os.Args[1:]
	pseudo := strings.Join(args, " ")

	if verbose {
		// remove -v from pseudo command
		pseudo = strings.Replace(pseudo, "-v", "", 1)
	}

	// remove flag -m and model
	pseudo = strings.Replace(pseudo, "-m "+model, "", 1)

	// fmt.Printf("Pseudo command: %s\n", pseudo)

	if pseudo == "" {
		thisCommand := filepath.Base(os.Args[0])
		flag.PrintDefaults()
		fmt.Println("Example: " + thisCommand + " Find files larger than 42kB in the current directory")
		os.Exit(1)
	}

	command, usage := makeCommand(pseudo, parentProcessName)
	fmt.Printf("Looks insanely complicated? Don't panic. The answer is ...\n")
	fmt.Printf("number of tokens used (total_tokens): %.1f\n", usage["total_tokens"])
	color.Cyan(command)

	var confirmation string
	fmt.Print("Run? (y/n) ")
	fmt.Scanln(&confirmation)

	var cmd *exec.Cmd

	if confirmation == "y" {
		switch runtime.GOOS {
		case "windows":
			switch parentProcessName {
			case "powershell.exe":
				cmd = exec.Command("powershell.exe", "-c", command)
			case "cmd.exe":
				cmd = exec.Command("cmd.exe", "/C", command)
			default:
				fmt.Printf("Error: unsupported shell %s on %s", parentProcessName, runtime.GOOS)
				os.Exit(1)
			}
		default: //Mac & Linux
			switch parentProcessName {
			case "zsh":
				cmd = exec.Command("zsh", "-c", command)
			case "bash":
				cmd = exec.Command("bash", "-c", command)
			default:
				fmt.Printf("Error: unsupported shell %s on %s", parentProcessName, runtime.GOOS)
				os.Exit(1)
			}
		}

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			fmt.Println("Error: ", err)
		}
	}

}
