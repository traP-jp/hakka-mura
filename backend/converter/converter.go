package converter

import (
	"context"
	_ "embed"
	"errors"
	"log"
	"os"

	"github.com/sashabaranov/go-openai"
)

type OpenAIConverter struct {
	cl *openai.Client
}

//go:embed prompt.txt
var systemPrompt string

func NewOpenAI() (*OpenAIConverter, error) {
	log.Println("SYSTEM PROMPT: ", systemPrompt)

	token, ok := os.LookupEnv("OPENAI_API_KEY")
	if !ok {
		return nil, errors.New("OPENAI_API_KEY is not set")
	}

	client := openai.NewClient(token)
	return &OpenAIConverter{cl: client}, nil
}

func (c *OpenAIConverter) ConvertMessage(ctx context.Context, originalMessage string) (string, error) {
	req := openai.ChatCompletionRequest{
		Model: openai.GPT6Luna,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: systemPrompt,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: originalMessage,
			},
		},
	}
	res, err := c.cl.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", err
	}

	return res.Choices[0].Message.Content, nil
}
