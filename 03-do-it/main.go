package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/prompts"
)

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

	llm, err := ollama.New(ollama.WithModel(modelName), ollama.WithServerURL(ollamaURL))
	if err != nil {
		log.Fatal(err)
	}

	/* From LangChainGo doc
	https://tmc.github.io/langchaingo/docs/modules/model_io/prompts/prompt_templates/#creating-prompt-templates-for-chat-messages
	*/

	//promptText := "Who is Kathryn Janeway?"

	
	prompt := prompts.NewChatPromptTemplate([]prompts.MessageFormatter{
		prompts.NewSystemMessagePromptTemplate(
			"You are a Star Trek expert.",
			nil,
		),
		prompts.NewHumanMessagePromptTemplate(
			`Who is {{.name}}?`,
			[]string{"name"},
		),
	})

	promptText, err := prompt.Format(map[string]any{
        "name": "Kathryn Janeway",
    })
    if err != nil {
        log.Fatal(err)
    }
	

    fmt.Println("🤖 prompt", promptText)

	_, err = llms.GenerateFromSinglePrompt(ctx, llm, promptText,
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			fmt.Print(string(chunk))

			return nil
		}))

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("")

}
