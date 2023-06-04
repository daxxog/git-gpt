package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/alecthomas/kong"
	openai "github.com/sashabaranov/go-openai"
	"gopkg.in/yaml.v2"
)

// Config represents OpenAI API configuration
type Config struct {
	Token string `yaml:"token"`
}

var CLI struct {
	Commit struct {
		All          bool `short:"a" help:"Stage and commit all changes (including unstaged). Under the hood, this passes the -a flag into git commit."`
		SkipAmendMsg bool `short:"m" help:"Commit directly with the AI-generated message without amendment."`
	} `cmd:"" help:"Commits files using an AI-generated message based on diff. See git-gpt -h commit for more details."`
}

func loadConfig() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(fmt.Sprintf("%s/.config/git-gpt/openai.yaml", home))
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func generateCommitMessage(diff string, config *Config) (string, error) {
	client := openai.NewClient(config.Token)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: "Write a commit message for these changes:\n" + diff,
				},
			},
		},
	)

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

func main() {
	ctx := kong.Parse(&CLI)
	switch ctx.Command() {
	case "commit":

		config, err := loadConfig()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		var diffCmd *exec.Cmd
		// If -a is passed, get all changes (including unstaged)
		if CLI.Commit.All {
			diffCmd = exec.Command("git", "diff")
		} else {
			// Otherwise, only get staged changes
			diffCmd = exec.Command("git", "diff", "--cached")
		}
		diffOutput, err := diffCmd.Output()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Print the diff output to the user
		fmt.Println(string(diffOutput))

		msg, err := generateCommitMessage(string(diffOutput), config)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		var commitCmd *exec.Cmd
		// First commit with the generated message
		// If -a is passed, include it in the command
		if CLI.Commit.All {
			commitCmd = exec.Command("git", "commit", "-a", "-m", msg)
		} else {
			commitCmd = exec.Command("git", "commit", "-m", msg)
		}

		// Set the command output to our standard output
		commitCmd.Stdout = os.Stdout
		commitCmd.Stderr = os.Stderr

		err = commitCmd.Run()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// If -m flag is not set, open the editor to let user amend the commit message
		if !CLI.Commit.SkipAmendMsg {
			commitCmd = exec.Command("git", "commit", "--amend")

			// Set the command output to our standard output
			commitCmd.Stdout = os.Stdout
			commitCmd.Stderr = os.Stderr

			err = commitCmd.Run()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

	default:
		panic(ctx.Command())
	}
}
