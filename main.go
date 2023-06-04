package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	openai "github.com/sashabaranov/go-openai"
	"gopkg.in/yaml.v2"
)

// Config represents OpenAI API configuration
type Config struct {
	Token string `yaml:"token"`
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
					Content: "Create a commit message for:\n" + diff,
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
	autoCommit := flag.Bool("a", false, "Auto commit flag")
	message := flag.Bool("m", false, "Skip message flag")
	flag.Parse()

	config, err := loadConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	diffCmd := exec.Command("git", "diff")
	diffOutput, err := diffCmd.Output()
	fmt.Println(string(diffOutput))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	msg, err := generateCommitMessage(string(diffOutput), config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if *autoCommit {
		var commitCmd *exec.Cmd
		if *message {
			commitCmd = exec.Command("git", "commit", "-a", "-m", msg)
		} else {
			commitCmd = exec.Command("git", "commit", "-a")
		}

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
