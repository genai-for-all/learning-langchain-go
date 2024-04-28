// 🚧 this is a work in progress

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/textsplitter"
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
	// search how it works
	splitter := textsplitter.NewMarkdownTextSplitter(
		textsplitter.WithChunkSize(1536), 
		textsplitter.WithChunkOverlap(128),
	)

	docs, err := textLoader.LoadAndSplit(ctx, splitter)
	if err != nil {
		log.Fatal("😡 when loading the docs:", err)
	}

	//fmt.Println("📚 docs", docs)

	modelName := os.Getenv("LLM")
	if modelName == "" {
		modelName = "llama3"
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
		chroma.WithDistanceFunction("cosine"),
		chroma.WithNameSpace("tada-namespace"),
	)
	if err != nil {
		log.Fatal("😡 when creating the Chroma vector store:", err)
	}
	// Add documents to the vector store
	documentIds, err := store.AddDocuments(ctx, docs)
	if err != nil {
		log.Fatal("😡 when adding the document to the store:", err)
	}
	fmt.Println("📝", documentIds)
	fmt.Println("📚", docs)

}
