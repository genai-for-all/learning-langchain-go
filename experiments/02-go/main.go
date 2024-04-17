package main

/*
go run main.go
go run main.go | jq -c '{ response, context }'
*/

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cast"
)

/*
Generate calls the generate API and returns `map[string]any, error`:

	{
		"response",
		"context",
		"created_at",
		"done",
		"total_duration",
		"load_duration",
		"prompt_eval_duration",
		"eval_count",
		"eval_duration",
	}
*/
func Generate(llmRequest map[string]interface{}, modelName, llmURL string) (map[string]any, error) {
	body, errMarshal := json.Marshal(llmRequest)
	if errMarshal != nil {
		return nil, errMarshal
	}
	resp, errPost := http.Post(llmURL+"/api/generate", "application/json;charset=utf-8", bytes.NewBuffer(body))
	if errPost != nil {
		return nil, errPost
	}
	defer resp.Body.Close()

	responseData, errRead := io.ReadAll(resp.Body)
	if errRead != nil {
		return nil, errRead
	}
	mapResponse := map[string]any{}
	errUnmarshal := json.Unmarshal(responseData, &mapResponse)
	if errUnmarshal != nil {
		return nil, errUnmarshal
	}
	return mapResponse, nil
}

func main() {

	modelName := os.Getenv("LLM")
	if modelName == "" {
		modelName = "gemma"
	}

	ollamaURL := os.Getenv("OLLAMA_BASE_URL")
	if ollamaURL == "" {
		ollamaURL = "http://host.docker.internal:11434"
	}

	firstRequest := map[string]interface{}{
		"model":  modelName,
		"prompt": "[Brief] Who is James T Kirk?",
		"stream": false,
	}

	response, err := Generate(firstRequest, modelName, ollamaURL)
	if err != nil {
		log.Fatal("😡", err)
	}
	fmt.Println("🤖", response["response"])

	secondRequest := map[string]interface{}{
		"model":   modelName,
		"prompt":  "Who is his best friend?",
		"stream":  false,
		"context": cast.ToIntSlice(response["context"]),
	}

	response, err = Generate(secondRequest, modelName, ollamaURL)
	if err != nil {
		log.Fatal("😡", err)
	}
	fmt.Println("🤖", response["response"])

}
