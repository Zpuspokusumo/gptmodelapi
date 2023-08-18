package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
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
		Stream:    true,
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
		stream, err := client.CreateChatCompletionStream(context.Background(), req)
		if err != nil {
			fmt.Printf("ChatCompletion error: %v\n", err)
			continue
		}
		fmt.Printf("Stream Resp: ")
		var chatresponse string
		chatresponse = ""
		for {
			response, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				fmt.Println("\nStream finished")
				break
			}

			if err != nil {
				fmt.Printf("\nStream error: %v\n", err)
				break
			}

			fmt.Printf(response.Choices[0].Delta.Content)
			chatresponse += response.Choices[0].Delta.Content
		}

		//fmt.Printf("%s\n\n", resp.Choices[0].Message.Content)
		//fmt.Printf("usage %+v\n\n", resp.Usage)
		//resets to last 2 messages (user msg and gpt resp)
		req.Messages = append([]openai.ChatCompletionMessage{}, req.Messages[0], nextmsg, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: chatresponse,
		})
		fmt.Print("> ")
	}
}
