// 🚧 this is a work in progress

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/prompts"
)

/*
LLM=tinydolphin go run main.go
LLM=tinyllama go run main.go
LLM=gemma go run main.go

LLM=deepseek-coder:instruct
LLM=deepseek-coder go run main.go
*/
func main() {

	ctx := context.Background()

	modelName := os.Getenv("LLM")
	if modelName == "" {
		modelName = "gemma"
	}

	ollamaURL := os.Getenv("OLLAMA_BASE_URL")
	if ollamaURL == "" {
		ollamaURL = "http://host.docker.internal:11434"
	}

	llm, err := ollama.New(
		ollama.WithModel(modelName),
		ollama.WithServerURL(ollamaURL),
		ollama.WithPredictRepeatLastN(64), // 👋
	)
	if err != nil {
		log.Fatal("😡 when creating the LLM object:", err)
	}

	prompt := prompts.NewChatPromptTemplate([]prompts.MessageFormatter{
		prompts.NewSystemMessagePromptTemplate(
			"You are a programming expert.",
			nil,
		),
		prompts.NewHumanMessagePromptTemplate(
			`{{.question}}`,
			[]string{"question"},
		),
	})

	promptText, _ := prompt.Format(map[string]any{
		"question": `write a simple hello world program in golang`,
	})

	fmt.Println("🤖 prompt 1", promptText)
	fmt.Println("📝 answer:")

	_, _ = llms.GenerateFromSinglePrompt(ctx, llm, promptText,
		llms.WithTemperature(0.0),
		//llms.WithStopWords([]string{"cargo"}),
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			fmt.Print(string(chunk))

			return nil
		}))

	fmt.Println("")

}
