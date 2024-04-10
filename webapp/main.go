// Package main : API to send prompt to Ollama
package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
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
	if httpPort == "" { httpPort = "8080" }
	
	ctx := context.Background()

	modelName := os.Getenv("LLM")
	ollamaURL := os.Getenv("OLLAMA_BASE_URL")

	llm, err := ollama.New(ollama.WithModel(modelName), ollama.WithServerURL(ollamaURL))
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	fileServerHtml := http.FileServer(http.Dir("public"))
	mux.Handle("/", fileServerHtml)

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

		response.Header().Set("Content-Type", "application/octet-stream")
		response.Header().Set("Transfer-Encoding", "chunked")
		//response.Header().Set("Connection", "Keep-Alive")
		//response.Header().Set("X-Content-Type-Options", "nosniff")
		
		_, err = llms.GenerateFromSinglePrompt(ctx, llm, question,
			llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
				log.Println("📝", string(chunk))
				response.Write(chunk)
				flusher.Flush()
				return nil
			}))
		if err != nil {
			log.Fatal(err)
		}

	})

	var errListening error
	log.Println("🌍 http server is listening on: " + httpPort)
	errListening = http.ListenAndServe(":"+httpPort, mux)

	log.Fatal(errListening)
}
