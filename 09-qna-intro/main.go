// 🚧 this is a work in progress

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

/*
curl http://host.docker.internal:11434/api/pull -d '{"name": "deepseek-coder:1.3b"}'

LLM=deepseek-coder:1.3b go run main.go

curl http://host.docker.internal:11434/api/pull -d '{"name": "tinyllama"}'
curl http://host.docker.internal:11434/api/pull -d '{"name": "tinydolphin"}'

LLM=tinyllama go run main.go
LLM=tinydolphin go run main.go
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
	)
	if err != nil {
		log.Fatal("😡 when creating the LLM object:", err)
	}

	prompt := "Who is Philippe Charrière?"

	completion, err := llms.GenerateFromSinglePrompt(ctx, llm, prompt)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(completion)

}
