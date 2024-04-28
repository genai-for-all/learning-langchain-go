// 🚧 this is a work in progress

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/prompts"
	"github.com/tmc/langchaingo/schema"
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
		chroma.WithNameSpace("tada-namespace"),
	)

	if err != nil {
		log.Fatal("😡 when creating the Chroma vector store:", err)
	}

	//question := "What are the actions of the game?"
	//question := "Tell me more about Goblins"
	question := "who are the monsters of the game?"

	// Search of similarity(ies)
	similarities, err := store.SimilaritySearch(
		ctx,
		question,
		3, // initially: 3 👋 change it
		vectorstores.WithScoreThreshold(0.7), // 👋 change it

	)
	if err != nil {
		log.Fatal("😡 when searching similarities:", err)
	}
	docs := []schema.Document{}

	for _, similarity := range similarities {
		doc := schema.Document{
			PageContent: "<doc>" + similarity.PageContent + "</doc>",
		}
		docs = append(docs, doc)

	}
	fmt.Println("👋", docs)

	prompt := prompts.NewChatPromptTemplate([]prompts.MessageFormatter{
		prompts.NewSystemMessagePromptTemplate(
			`You are the dungeon master, 
			expert at interpreting and answering questions based on provided sources.
			Using the provided context, answer the user's question 
			to the best of your ability using only the resources provided. 
			Be verbose!`,
			nil,
		),
		prompts.NewSystemMessagePromptTemplate(
			"<context> {{.context}} </context>",
			[]string{"context"}, // link because chains.NewStuffDocuments
		),
		/*
			StuffDocuments is a chain that combines documents with a separator and uses
			the stuffed documents in an LLMChain. The input values to the llm chain
			contains all input values given to this chain, and the stuffed document as
			a string in the key specified by the "DocumentVariableName" field that is
			by default set to "context".

			See https://github.com/tmc/langchaingo/blob/main/chains/stuff_documents.go#L18
		*/
		prompts.NewHumanMessagePromptTemplate(
			`Now, answer this question using the above context:
			{{.question}}`,
			[]string{"question"},
		),
	})

	llmChain := chains.NewLLMChain(llm, prompt)
	chain := chains.NewStuffDocuments(llmChain)

	answer, err := chains.Call(context.Background(), chain, map[string]any{
		"input_documents": docs, // link to context
		"question":        question,
	})
	if err != nil {
		log.Fatal("😡", err)
	}
	fmt.Println("")
	fmt.Println("question:", question)
	fmt.Println("answer:")
	fmt.Println(answer["text"])

}
