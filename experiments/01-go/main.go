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
)

func main() {

	modelName := os.Getenv("LLM")
	if modelName == "" {
		modelName = "gemma"
	}

	ollamaURL := os.Getenv("OLLAMA_BASE_URL")
	if ollamaURL == "" {
		ollamaURL = "http://host.docker.internal:11434"
	}

	mapRequest := map[string]interface{}{
		"model":  modelName,
		"prompt": "Who is James T Kirk?",
		"stream": false,
	}
	body, errMarshal := json.Marshal(mapRequest)
	if errMarshal != nil {
		log.Println("😡", errMarshal)
		os.Exit(1)
	}

	resp, errPost := http.Post(ollamaURL+"/api/generate", "application/json;charset=utf-8", bytes.NewBuffer(body))

	if errPost != nil {
		log.Println(errPost.Error())
		os.Exit(1)
	}

	defer resp.Body.Close()

	responseData, errRead := io.ReadAll(resp.Body)
	if errRead != nil {
		log.Println(errRead)
		os.Exit(1)
	}
	//fmt.Println(string(responseData))

	mapResponse := map[string]any{}
	errUnmarshal := json.Unmarshal(responseData, &mapResponse)
	if errUnmarshal != nil {
		log.Println(errUnmarshal)
		os.Exit(1)
	}

	fmt.Println(mapResponse["response"])
	fmt.Println(mapResponse["context"]) // []interface {}

	
	/*
	fmt.Println(mapResponse["created_at"])
	fmt.Println(mapResponse["done"])
	fmt.Println(mapResponse["total_duration"])
	fmt.Println(mapResponse["load_duration"])
	fmt.Println(mapResponse["prompt_eval_duration"])
	fmt.Println(mapResponse["eval_count"])
	fmt.Println(mapResponse["eval_duration"])
	*/

}
