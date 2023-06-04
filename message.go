package main

import (
	"context"

	openai "github.com/sashabaranov/go-openai"
)

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
