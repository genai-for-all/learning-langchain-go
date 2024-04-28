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


func main() {

	ctx := context.Background()

	/*
		Load and split a markdown document
	*/
	filePath := os.Getenv("MD_FILE_PATH")
	if filePath == "" {
		filePath = "./rules.md"
	}
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal("😡 when reading the file:", err)
	}
	textLoader := documentloaders.NewText(file)

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
		chromaURL = "http://chroma-db:8000"
		//chromaURL = "http://host.docker.internal:8000"
	}
	chromaNameSpace := os.Getenv("CHROMA_NAMESPACE")
	if chromaNameSpace == "" {
		chromaNameSpace = "aethelgard-namespace"
	}

	log.Println("📦 creating Chroma vector store:", chromaNameSpace)

	// Create a new Chroma vector store.
	store, err := chroma.New(
		chroma.WithChromaURL(chromaURL),
		chroma.WithEmbedder(llmEmbeder),
		chroma.WithDistanceFunction("cosine"),
		chroma.WithNameSpace(chromaNameSpace),
	)
	if err != nil {
		log.Fatal("😡 when creating the Chroma vector store:", err)
	}
	// Add documents to the vector store
	documentIds, err := store.AddDocuments(ctx, docs)
	if err != nil {
		log.Fatal("😡 when adding the document to the store:", err)
	}
	fmt.Println("📚", docs)

	fmt.Println("📝", documentIds)
	

}
