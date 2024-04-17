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
	"time"
)

// Structure to hold the llm request
type LlmRequest struct {
	Model   string `json:"model"`
	Prompt  string `json:"prompt"`
	Stream  bool   `json:"stream"`
	Context []int  `json:"context,omitempty"`
}

// Structure to hold the llm response
type LlmResponse struct {
	Response           string    `json:"response"`
	Context            []int     `json:"context,omitempty"`
	CreatedAt          time.Time `json:"created_at"`
	Done               bool      `json:"done"`
	TotalDuration      int       `json:"total_duration"`
	LoadDuration       int       `json:"load_duration"`
	PromptEvalDuration int       `json:"prompt_eval_duration"`
	EvalCount          int       `json:"eval_count"`
	EvalDuration       int       `json:"eval_duration"`
}

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
func Generate(llmRequest LlmRequest, modelName, llmURL string) (LlmResponse, error) {
	body, errMarshal := json.Marshal(llmRequest)
	if errMarshal != nil {
		return LlmResponse{}, errMarshal
	}
	resp, errPost := http.Post(llmURL+"/api/generate", "application/json;charset=utf-8", bytes.NewBuffer(body))
	if errPost != nil {
		return LlmResponse{}, errPost
	}
	defer resp.Body.Close()

	responseData, errRead := io.ReadAll(resp.Body)
	if errRead != nil {
		return LlmResponse{}, errRead
	}
	mapResponse := LlmResponse{}
	errUnmarshal := json.Unmarshal(responseData, &mapResponse)
	if errUnmarshal != nil {
		return LlmResponse{}, errUnmarshal
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

	firstRequest := LlmRequest{
		Model:  modelName,
		Prompt: "[Brief] Who is James T Kirk?",
		Stream: false,
	}

	body, err := Generate(firstRequest, modelName, ollamaURL)
	if err != nil {
		log.Fatal("😡", err)
	}
	fmt.Println("🤖", body.Response)
	//fmt.Println("🤖", body.Context)

	secondRequest := LlmRequest{
		Model:   modelName,
		Prompt:  "Who is his best friend?",
		Stream:  false,
		Context: body.Context,
	}

	body, err = Generate(secondRequest, modelName, ollamaURL)
	if err != nil {
		log.Fatal("😡", err)
	}
	fmt.Println("🤖", body.Response)

}
