package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

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

	llm, err := ollama.New(ollama.WithModel(modelName), ollama.WithServerURL(ollamaURL))
	if err != nil {
		log.Fatal(err)
	}

	prompt := "Who is Michael Burnham?"

	//prompt := "Who are the main caracters of StarTrek Discovery?"


	_, err = llms.GenerateFromSinglePrompt(ctx, llm, prompt,
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			fmt.Print(string(chunk))

			return nil
		}))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("")

}
