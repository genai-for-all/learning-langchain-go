// Package main : API to send prompt to Ollama
package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
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
GetBytesBody returns the body of an HTTP request as a []byte.
  - It takes a pointer to an http.Request as a parameter.
  - It returns a []byte.
*/
func GetBytesBody(request *http.Request) []byte {
	body := make([]byte, request.ContentLength)
	request.Body.Read(body)
	return body
}

func main() {
	ctx := context.Background()

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}
	modelName := os.Getenv("LLM")
	if modelName == "" {
		modelName = "gemma"
	}
	ollamaURL := os.Getenv("OLLAMA_BASE_URL")
	if ollamaURL == "" {
		ollamaURL = "http://host.docker.internal:11434"
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

	repeat := 64
	//temperature := 1.0

	llm, err := ollama.New(
		ollama.WithModel(modelName),
		ollama.WithServerURL(ollamaURL),
		ollama.WithPredictRepeatLastN(repeat),
	)
	if err != nil {
		log.Fatal("😡 when creating the LLM object:", err)
	}

	llmEmbeder, err := embeddings.NewEmbedder(llm)
	if err != nil {
		log.Fatal("😡 when creating the Embedder object:", err)
	}

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

	mux := http.NewServeMux()

	fileServerHtml := http.FileServer(http.Dir("public"))
	mux.Handle("/", fileServerHtml)

	shouldIStopTheCompletion := false

	// Cancel/Stop the generation of the completion
	mux.HandleFunc("DELETE /cancel-request", func(response http.ResponseWriter, request *http.Request) {
		shouldIStopTheCompletion = true
		response.Write([]byte("🚫 Cancelling request..."))
	})

	mux.HandleFunc("/prompt", func(response http.ResponseWriter, request *http.Request) {

		// add a flusher
		flusher, ok := response.(http.Flusher)
		if !ok {
			response.Write([]byte("😡 Error: expected http.ResponseWriter to be an http.Flusher"))
		}

		body := GetBytesBody(request)

		// unmarshal the json data
		var data map[string]string

		err = json.Unmarshal(body, &data)
		if err != nil {
			response.Write([]byte("😡 Error: " + err.Error()))
		}

		question := data["question"]
		systemMessage := data["system"]

		// Search of similarity(ies)
		similarities, err := store.SimilaritySearch(
			ctx,
			question,
			10,
			vectorstores.WithScoreThreshold(0.8),

		)

		log.Println("🔎", similarities)

		if err != nil {
			log.Println("😡 when searching similaritie")
			response.Write([]byte("😡 when searching similaritie"))
		}

		docs := []schema.Document{}

		for _, similarity := range similarities {
			doc := schema.Document{
				PageContent: "<doc>" + similarity.PageContent + "</doc>",
			}
			log.Println("📝", doc.PageContent)
			docs = append(docs, doc)
		}

		prompt := prompts.NewChatPromptTemplate([]prompts.MessageFormatter{
			prompts.NewSystemMessagePromptTemplate(
				systemMessage,
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
				`Now, answer this question using the above context and docs:
				{{.question}}`,
				[]string{"question"},
			),
		})

		llmChain := chains.NewLLMChain(llm, prompt)
		chain := chains.NewStuffDocuments(llmChain)

		answer := ""
		_, err = chains.Call(
			context.Background(),
			chain, map[string]any{
				"input_documents": docs, // link to context
				"question":        question,
			},
			chains.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
				response.Write(chunk)
				flusher.Flush()
				if !shouldIStopTheCompletion {
					answer += string(chunk)
					return nil
				} else {
					return errors.New("🚫 Cancelling request")
				}
			}))

		response.Header().Set("Content-Type", "application/octet-stream")
		response.Header().Set("Transfer-Encoding", "chunked")

		if err != nil {
			shouldIStopTheCompletion = false
			response.Write([]byte("bye"))
		}

	})

	var errListening error
	log.Println("🌍 http server is listening on: " + httpPort)
	errListening = http.ListenAndServe(":"+httpPort, mux)

	log.Fatal(errListening)
}
