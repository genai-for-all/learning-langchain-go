// 🚧 this is a work in progress

package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/tmc/langchaingo/llms/ollama"
)

func removeIndent(s string) string {

	lines := strings.Split(s, "\n")
	for i, l := range lines {
		lines[i] = strings.TrimPrefix(l, "    ")
	}
	return strings.Join(lines, "\n")
}

/*
LLM=tinydolphin go run main.go
LLM=tinyllama go run main.go
LLM=gemma go run main.go

LLM=deepseek-coder:instruct
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
		ollama.WithFormat("json"),
	)
	if err != nil {
		log.Fatal("😡 when creating the LLM object:", err)
	}

}
