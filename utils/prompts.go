package utils

import (
	"errors"

	"github.com/sashabaranov/go-openai"
)

func ProcessCompletionPrompt(req *openai.CompletionRequest, systemPrompt string) error {
	// req.Prompt can be a string or an array of strings
	// if its a string, prepend the systemPrompt to it
	// if its an array, prepend the systemPrompt to each element
	// then reassign it to req.Prompt
	// if its neither, return an error

	switch prompt := req.Prompt.(type) {
	case string:
		req.Prompt = systemPrompt + "\n\n" + prompt
	case []string:
		for i, p := range prompt {
			prompt[i] = systemPrompt + "\n\n" + p
		}
		req.Prompt = prompt
	default:
		return errors.New("invalid prompt type")
	}

	return nil
}

func ProcessChatCompletionPrompt(req *openai.ChatCompletionRequest, systemPrompt string) {
	// loop through messages, check if a prompt with role "system" exists
	// if it does, replace with the systemPrompt
	// if it doesn't, prepend the systemPrompt to the messages

	for i, message := range req.Messages {
		if message.Role == "system" {
			req.Messages[i].Content = systemPrompt
			return
		}
	}

	// add it to the beginning of the messages
	req.Messages = append([]openai.ChatCompletionMessage{{Role: "system", Content: systemPrompt}}, req.Messages...)
}
