// 🚧 this is a work in progress

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/vectorstores"
	"github.com/tmc/langchaingo/vectorstores/chroma"
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

	//fmt.Println("📚 docs", docs)

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

	llmEmbeder, err := embeddings.NewEmbedder(llm)
	if err != nil {
		log.Fatal("😡 when creating the Embedder object:", err)
	}

	chromaURL := os.Getenv("CHROMA_URL")
	if chromaURL == "" {
		chromaURL = "http://localhost:8000"
		//chromaURL = "http://host.docker.internal:8000"
	}

	// Create a new Chroma vector store.
	store, err := chroma.New(
		chroma.WithChromaURL(chromaURL),
		chroma.WithEmbedder(llmEmbeder),
		chroma.WithNameSpace("tada-namespace"),
	)

	if err != nil {
		log.Fatal("😡 when creating the Chroma vector store:", err)
	}


	// Search of similarity(ies)
	similarities, err := store.SimilaritySearch(
		ctx,
		//"what is the title of this game?",
		"who are the monsters of the game?",
		3,
		vectorstores.WithScoreThreshold(0.8), // 🤔

	)
	if err != nil {
		log.Fatal("😡 when searching similarities:", err)
	}

	fmt.Println("👋", similarities)

	for _, doc := range similarities {
		fmt.Println(doc.PageContent)
	}

}
