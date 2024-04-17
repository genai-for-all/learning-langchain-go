package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/memory"
	"github.com/tmc/langchaingo/prompts"
)

/*
LLM=tinydolphin go run main.go
LLM=tinyllama go run main.go
LLM=gemma go run main.go
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
		log.Fatal(err)
	}

	history := memory.NewChatMessageHistory()

	systemMessage := "You are a Star Trek expert."

	prompt1 := prompts.NewChatPromptTemplate([]prompts.MessageFormatter{
		prompts.NewSystemMessagePromptTemplate(
			systemMessage,
			nil,
		),
		prompts.NewHumanMessagePromptTemplate(
			`[Brief] {{.question}}`,
			[]string{"question"},
		),
	})

	//myFirstQuestion := "Who is Jean Luc Picard?"
	myFirstQuestion := "Who is James T Kirk?"

	promptText1, _ := prompt1.Format(map[string]any{
		"question": myFirstQuestion,
	})

	fmt.Println("🤖 prompt 1", promptText1)
	fmt.Println("📝 answer:")

	var answer string
	_, err = llms.GenerateFromSinglePrompt(ctx, llm, promptText1,

		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			fmt.Print(string(chunk))
			answer += string(chunk)
			return nil
		}))
	if err != nil {
		log.Fatal("😡", err)
	}

	// create an history
	history.AddUserMessage(ctx, myFirstQuestion)
	history.AddAIMessage(ctx, answer)

	prompt2 := prompts.NewChatPromptTemplate([]prompts.MessageFormatter{
		prompts.NewSystemMessagePromptTemplate(
			systemMessage,
			nil,
		),
		// Insert history
		prompts.NewGenericMessagePromptTemplate("history", "{{range .historyMessages}}{{.GetContent}}\n{{end}}", []string{"history"}),
		prompts.NewHumanMessagePromptTemplate(
			`[Brief] {{.question}}`,
			[]string{"question"},
		),
	})

	historyMessages, _ := history.Messages(ctx)
	mySecondQuestion := "Who is his best friend?"

	promptText2, _ := prompt2.Format(map[string]any{
		"historyMessages": historyMessages,
		"question":        mySecondQuestion,
	})

	fmt.Println("")
	fmt.Println("")

	fmt.Println("🤖 prompt 2", promptText2)
	fmt.Println("📝 answer:")

	_, err = llms.GenerateFromSinglePrompt(ctx, llm, promptText2,

		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			fmt.Print(string(chunk))
			return nil
		}))
	if err != nil {
		log.Fatal("😡", err)
	}
}
