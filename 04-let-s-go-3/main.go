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

	prompt1 := prompts.NewChatPromptTemplate([]prompts.MessageFormatter{
		prompts.NewSystemMessagePromptTemplate(
			"You are a Star Trek expert.",
			nil,
		),
		prompts.NewHumanMessagePromptTemplate(
			`[Brief] {{.question}}`,
			[]string{"question"},
		),
	})

	promptText1, _ := prompt1.Format(map[string]any{
		"question": "Who is James T Kirk?",
	})

	fmt.Println("🤖 prompt 1", promptText1)
	fmt.Println("📝 answer:")

	var answer string
	_, _ = llms.GenerateFromSinglePrompt(ctx, llm, promptText1,

		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			fmt.Print(string(chunk))
			answer += string(chunk)
			return nil
		}))

	prompt2 := prompts.NewChatPromptTemplate([]prompts.MessageFormatter{
		prompts.NewSystemMessagePromptTemplate(
			"You are a Star Trek expert.",
			nil,
		),
		prompts.NewHumanMessagePromptTemplate(`[Brief] {{.initialQuestion}}`, []string{"initialQuestion"}),
		prompts.NewAIMessagePromptTemplate(`Previous conversation history: {{.history}}`, []string{"history"}),
		prompts.NewHumanMessagePromptTemplate(
			`[Brief] {{.question}}`,
			[]string{"question"},
		),
	})

	promptText2, _ := prompt2.Format(map[string]any{
		"question": "Who is his best friend?",
		"history":  answer,
		"initialQuestion": "Who is James T Kirk?",
	})

	fmt.Println("")
	fmt.Println("")

	fmt.Println("✋Prompt 2:")

	for idx, item := range prompt2.Messages {
		fmt.Println(" -",idx,  item)
	}

	fmt.Println("🤖 prompt 2", promptText2)
	fmt.Println("📝 answer:")
	
	_, _ = llms.GenerateFromSinglePrompt(ctx, llm, promptText2,

		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			fmt.Print(string(chunk))

			return nil
		}))

}
