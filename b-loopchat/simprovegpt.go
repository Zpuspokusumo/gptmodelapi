package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/sashabaranov/go-openai"
)

// test system fro qubio
// You are Qubio, and you will respond to questions for Qubio.
// Qubio is a young helpful assistant doing he best to help people achieve their goals
// Although Qubio is young he is wise, resourceful, and smart, all to help people
// Please try to use a space theme in your replies, and repond cheerfully and helpfully
// If possible, please also use a magical and whimsical theme in your replies

func main() {
	client := openai.NewClient(os.Getenv("API-KEY"))

	fmt.Print("Input username> ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	err := scanner.Err()
	if err != nil {
		log.Fatal(err)
	}
	userid := scanner.Text()
	fmt.Println("user is: ", userid)

	fmt.Print("Input system content> ")
	scannerc := bufio.NewScanner(os.Stdin)
	scannerc.Scan()
	err = scannerc.Err()
	if err != nil {
		log.Fatal(err)
	}
	contx := scannerc.Text()
	fmt.Println("context is: ", contx)

	req := openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: contx,
			},
		},
		MaxTokens: 1024,
		User:      userid,
	}
	fmt.Println("Conversation")
	fmt.Println("---------------------")
	fmt.Print("> ")
	s := bufio.NewScanner(os.Stdin)
	i := 0
	for s.Scan() {
		i++
		fmt.Printf("text no %d\n", i)
		nextmsg := openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: s.Text(),
		}
		req.Messages = append(req.Messages, nextmsg)
		resp, err := client.CreateChatCompletion(context.Background(), req)
		if err != nil {
			fmt.Printf("ChatCompletion error: %v\n", err)
			continue
		}
		fmt.Printf("%s\n\n", resp.Choices[0].Message.Content)
		fmt.Printf("usage %+v\n\n", resp.Usage)
		//resets to last 2 messages (user msg and gpt resp)
		req.Messages = append([]openai.ChatCompletionMessage{}, req.Messages[0], nextmsg, resp.Choices[0].Message)
		fmt.Print("> ")
	}
}
