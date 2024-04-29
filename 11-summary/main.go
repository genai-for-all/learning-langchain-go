// 🚧 this is a work in progress

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/textsplitter"
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
		modelName = "tinyllama" // better result with gemma (several choices), but tinyllama is faster, phi3 is not good for this
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

	/*
		Load and split a markdown document
	*/
	file, err := os.Open("./rules.md")
	if err != nil {
		log.Fatal("😡 when reading the file:", err)
	}
	textLoader := documentloaders.NewText(file)

	//docs, err := textLoader.Load(ctx)
	//textsplitter.NewRecursiveCharacter()
	// the chunk size and overlap has a real impact on the vector similarity

	splitter := textsplitter.NewMarkdownTextSplitter(
		textsplitter.WithChunkSize(1536),
		textsplitter.WithChunkOverlap(128),
	)

	docs, err := textLoader.LoadAndSplit(ctx, splitter)
	if err != nil {
		log.Fatal("😡 when loading the docs:", err)
	}

	llmSummarizationChain := chains.LoadRefineSummarization(llm)

	fmt.Println("--------------------------------------------------")

	outputValues, err := chains.Call(
		ctx,
		llmSummarizationChain,
		map[string]any{"input_documents": docs},
		chains.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {

			fmt.Print(string(chunk))
			if len(chunk) == 0 {
				fmt.Println("")
				fmt.Println("--------------------------------------------------")
			}

			return nil
		}),
	)
	if err != nil {
		log.Fatal("😡 when summarizing the docs:", err)
	}
	out := outputValues["text"].(string)
	fmt.Println("📝 Last summary:")
	fmt.Println(out)

}
