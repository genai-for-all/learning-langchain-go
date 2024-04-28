// Package main : API to send prompt to Ollama
package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/memory"
	"github.com/tmc/langchaingo/prompts"
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

	var httpPort = os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	ctx := context.Background()

	history := memory.NewChatMessageHistory()
	history.AddAIMessage(ctx, "this is the beginning of the conversation")

	modelName := os.Getenv("LLM")
	ollamaURL := os.Getenv("OLLAMA_BASE_URL")

	repeat := 64
	temperature := 1.0

	llm, err := ollama.New(
		ollama.WithModel(modelName),
		ollama.WithServerURL(ollamaURL),
		ollama.WithPredictRepeatLastN(repeat),
	)

	if err != nil {
		log.Fatal(err)
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

	mux.HandleFunc("GET /message-history", func(response http.ResponseWriter, request *http.Request) {
		messages, _ := history.Messages(ctx)
		response.Header().Set("Content-Type", "application/json")
		
		jsonMessages, err := json.Marshal(messages)
		if err != nil {
			response.Write([]byte("😡 Error: " + err.Error()))
		}
		response.Write([]byte(jsonMessages))
	})

	// Clear the history
	mux.HandleFunc("DELETE /clear-history", func(response http.ResponseWriter, request *http.Request) {
		// TODO: test the error
		history.Clear(ctx)
		response.Write([]byte("🗑️ history cleared"))
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

		err := json.Unmarshal(body, &data)
		if err != nil {
			response.Write([]byte("😡 Error: " + err.Error()))
		}

		question := data["question"]
		systemMessage := data["system"]
		contextMessage := data["context"]

		prompt := prompts.NewChatPromptTemplate([]prompts.MessageFormatter{
			prompts.NewSystemMessagePromptTemplate(
				systemMessage,
				nil,
			),
			
			prompts.NewSystemMessagePromptTemplate(
				contextMessage, nil,
			),
			
			// Insert history
			prompts.NewGenericMessagePromptTemplate(
				"history",
				"{{range .historyMessages}}{{.GetContent}}\n{{end}}",
				[]string{"history"},
			),
			prompts.NewHumanMessagePromptTemplate(
				`{{.question}}`,
				[]string{"question"},
			),
		})


		historyMessages, _ := history.Messages(ctx)
		textPrompt, _ := prompt.Format(map[string]any{
			"historyMessages": historyMessages,
			"question":        question,
		})

		response.Header().Set("Content-Type", "application/octet-stream")
		response.Header().Set("Transfer-Encoding", "chunked")

		answer := ""
		_, err = llms.GenerateFromSinglePrompt(ctx, llm, textPrompt,
			llms.WithTemperature(temperature),
			llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
				log.Println("📝", string(chunk))
				response.Write(chunk)
				flusher.Flush()
				if !shouldIStopTheCompletion {
					answer += string(chunk)
					return nil
				} else {
					return errors.New("🚫 Cancelling request")
				}

			}))

		if err != nil {
			// I should test the kind of the error
			shouldIStopTheCompletion = false
			response.Write([]byte("bye"))
		}

		history.AddUserMessage(ctx, question)
		history.AddAIMessage(ctx, answer)

	})

	var errListening error
	log.Println("🌍 http server is listening on: " + httpPort)
	errListening = http.ListenAndServe(":"+httpPort, mux)

	log.Fatal(errListening)
}
