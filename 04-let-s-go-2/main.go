package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/schema"
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

	llm, err := ollama.New(
		ollama.WithModel(modelName), 
		ollama.WithServerURL(ollamaURL),
		//ollama.WithVerbose(true),
	)
	if err != nil {
		log.Fatal(err)
	}

	systemInstructions := []llms.ContentPart{
		llms.TextContent{Text: "You are a Star Trek expert."},
	}

	// ✋ with gemma:2b the LLM cannot answer,
	// with gemma the LLM can answer
	firstQuestionParts := []llms.ContentPart{
		llms.TextContent{Text: "Who is James T Kirk?"},
		//llms.TextContent{Text: "Who is Jonathan Archer?"},

	}

	content := []llms.MessageContent{
		{Role: schema.ChatMessageTypeSystem, Parts: systemInstructions},
		{Role: schema.ChatMessageTypeHuman, Parts: firstQuestionParts},
	}

	fmt.Println("📝", content)
	fmt.Println("🙂:", content[len(content)-1].Parts[0])
	fmt.Println("🤖:")

	// Generate the response
	var txtResponse string

	_, _ = llm.GenerateContent(ctx, content,
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			fmt.Print(string(chunk))
			txtResponse += string(chunk)
			return nil
		}))

	fmt.Println("")
	fmt.Println("")

	responseParts := []llms.ContentPart{
		llms.TextContent{Text: txtResponse},
	}

	secondQuestionParts := []llms.ContentPart{
		llms.TextContent{Text: "Who is his best friend?"},
	}

	content = append(content,
		llms.MessageContent{ // former context
			Role: schema.ChatMessageTypeAI, Parts: responseParts,
		},
		llms.MessageContent{ // new question
			Role: schema.ChatMessageTypeHuman, Parts: secondQuestionParts,
		},
	)

	fmt.Println("📝 content:", )

	for idx, item := range content { // item is llms.MessageContent
		fmt.Println(" -",idx,  item)
	}

	fmt.Println("🙂:", content[len(content)-1].Parts[0])
	fmt.Println("🤖:")

	_, _ = llm.GenerateContent(ctx, content,
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			fmt.Print(string(chunk))
			return nil
		}))

	fmt.Println("")

	

}
