package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/schema"
)

/*
LLM=tinydolphin go run main.go
LLM=tinyllama go run main.go
LLM=gemma go run main.go

LLM=deepseek-coder:instruct
*/
func main() {

	ctx := context.Background()

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
		ollama.WithFormat("json"),
	)
	if err != nil {
		log.Fatal("😡 when creating the LLM object:", err)
	}

	// Define the specification of the functions
	var functions = []llms.FunctionDefinition{
		{
			Name:        "hello",
			Description: "When you want to say hello to a given person, give the name of this person",
			Parameters: json.RawMessage(`{
				"type": "object", 
				"properties": {
					"name": {"type": "string", "description": "the name of the person"} 
				}, 
				"required": ["name"]
			}`),
		},
		{
			Name:        "addNumbers",
			Description: "Make an addition of the given numbers",
			Parameters: json.RawMessage(`{
				"type": "object", 
				"properties": {
					"numbers": {"type": "array", "description": "the list of the numbers to add"} 
				}, 
				"required": ["numbers"]
			}`),
			// ✋ ex: if you use array type, the final output is tools and not tool
			// try with another LLM to check the behaviour
			// it works better with gemma and it's ok too with gemma:2b
			// LLM=gemma:2b go run main.go
		},
	}
	/*
		"numbers": {"type": "array", "description": "the list of the numbers to add"}
	*/
	toolsList, err := json.Marshal(functions)
	if err != nil {
		log.Fatal("😡 tools to json string:", err)
	}
	toolsListStr := string(toolsList)

	systemMessageStr := fmt.Sprintf(`You have access to the following tools:

	%s
	
	To use a tool, respond with a JSON object with the following structure: 
	{
		"tool": <name of the called tool>,
		"tool_input": <parameters for the tool matching the above JSON schema>
	}
	`, toolsListStr)

	var msgs []llms.MessageContent

	// system message defines the available tools.
	msgs = append(msgs, llms.TextParts(
		schema.ChatMessageTypeSystem,
		systemMessageStr,
	))

	msgs = append(msgs, llms.TextParts(
		schema.ChatMessageTypeHuman,
		"Please, add 18 and 24",
		//"Hey, say hello to Bob Morane",
		//"Hey, say hello to Bob",
	))
	// it does not work with `Bob Morane`

	fmt.Println("📝:", "\n", msgs)

	resp, err := llm.GenerateContent(ctx, msgs)
	if err != nil {
		log.Fatal("😡 issue when generating content:", err)
	}

	choice1 := resp.Choices[0]

	/*
		fmt.Println("👋", resp.Choices)
		fmt.Println("")
		fmt.Println("🤖:", "\n", choice1)
	*/

	fmt.Println("📦:", "\n", choice1.Content) // sometime tools, sometime tool

	selectedTool := map[string]any{}
	errTool := json.Unmarshal([]byte(choice1.Content), &selectedTool)
	if errTool != nil {
		log.Fatal("😡 issue when unmarshalling:", errTool)
	}

	//fmt.Println("🛠️:", "\n", selectedTool["tool"])
	fmt.Println("🛠️:", "\n", selectedTool)

	switch toolName := selectedTool["tool"].(string); toolName {
	case "hello":
		name := selectedTool["tool_input"].(map[string]any)["name"]
		fmt.Println(hello(name.(string)))
	case "addNumbers":
		numbers := selectedTool["tool_input"]
		fmt.Println(numbers)
	default:
		fmt.Println("😢", "tool not found:", toolName)
	}

}

func hello(name string) string {
	return "👋 Hello " + name + " 🙂"
}
