package main

import (
	"context"
	"fmt"
	"os"

	"github.com/sashabaranov/go-openai"
)

func summarizegpt(msgs []openai.ChatCompletionMessage) openai.ChatCompletionMessage {
	req := openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "Summarize the following messages by only mentioning their topics and try to use specific vocabulary",
			},
		},
		MaxTokens: 1024,
	}
	req.Messages = append(req.Messages, msgs...)

	client := openai.NewClient(os.Getenv("API-KEY"))

	resp, err := client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
	}
	fmt.Printf("///////// THIS IS THE DIALBACK SUMMARIZE GPT FUNC \n%s\n\n", resp.Choices[0].Message.Content)
	fmt.Printf("usage %+v\n\n", resp.Usage)

	return openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: resp.Choices[0].Message.Content,
	}
}
