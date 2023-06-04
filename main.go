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
		AutoCommit bool `short:"a" help:"Auto commit flag"`
		SkipMsg    bool `short:"m" help:"Skip message flag"`
	} `cmd:"" help:"Commit files."`
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

		diffCmd := exec.Command("git", "diff")
		diffOutput, err := diffCmd.Output()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		msg, err := generateCommitMessage(string(diffOutput), config)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if CLI.Commit.AutoCommit {
			var commitCmd *exec.Cmd
			// First commit with the generated message
			commitCmd = exec.Command("git", "commit", "-a", "-m", msg)

			// Set the command output to our standard output
			commitCmd.Stdout = os.Stdout
			commitCmd.Stderr = os.Stderr

			err = commitCmd.Run()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			// If -m flag is not set, open the editor to let user amend the commit message
			if !CLI.Commit.SkipMsg {
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
		}

	default:
		panic(ctx.Command())
	}
}
