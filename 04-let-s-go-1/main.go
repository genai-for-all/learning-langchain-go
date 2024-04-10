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

	prompt := prompts.NewChatPromptTemplate([]prompts.MessageFormatter{
		prompts.NewSystemMessagePromptTemplate(
			"You are a Star Trek expert.",
			nil,
		),
		prompts.NewHumanMessagePromptTemplate(
			`Please answer my question precisely. This is my question: {{.question}}`,
			[]string{"question"},
		),
	})

	promptText1, _ := prompt.Format(map[string]any{ //! usually you shoul handle the error
		"question": "Who is Jonathan Archer?",
	})

	fmt.Println("🤖 prompt 1", promptText1)
	fmt.Println("📝 answer:")

	_, _ = llms.GenerateFromSinglePrompt(ctx, llm, promptText1, //! usually you shoul handle the error

		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			fmt.Print(string(chunk))

			return nil
		}))

	fmt.Println("")
	fmt.Println("")

	promptText2, _ := prompt.Format(map[string]any{
		"question": "Who are his crewmates?",
	})

	fmt.Println("🤖 prompt 2", promptText2)
	fmt.Println("📝 answer:")

	_, _ = llms.GenerateFromSinglePrompt(ctx, llm, promptText2,

		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			fmt.Print(string(chunk))

			return nil
		}))

	fmt.Println("")

}
