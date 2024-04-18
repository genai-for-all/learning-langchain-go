// 🚧 this is a work in progress

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/prompts"
	"github.com/tmc/langchaingo/schema"
)

func removeIndent(s string) string {

	lines := strings.Split(s, "\n")
	for i, l := range lines {
		lines[i] = strings.TrimPrefix(l, "    ")
	}
	return strings.Join(lines, "\n")
}

/*
curl http://host.docker.internal:11434/api/pull -d '{"name": "deepseek-coder:1.3b"}'

LLM=deepseek-coder:1.3b go run main.go

LLM=tinyllama go run main.go
*/
func main() {

	//ctx := context.Background()

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
	)
	if err != nil {
		log.Fatal("😡 when creating the LLM object:", err)
	}

	content1 := `<doc>
	Michael Burnham is the main character on the Star Trek series, Discovery.  
	She's a human raised on the logical planet Vulcan by Spock's father.  
	Burnham is intelligent and struggles to balance her human emotions with Vulcan logic.  
	She's become a Starfleet captain known for her determination and problem-solving skills.
	Originally played by actress Sonequa Martin-Green
	</doc>`

	content2 := `<doc>
	James T. Kirk, also known as Captain Kirk, is a fictional character from the Star Trek franchise.  
	He's the iconic captain of the starship USS Enterprise, 
	boldly exploring the galaxy with his crew.  
	Originally played by actor William Shatner, 
	Kirk has appeared in TV series, movies, and other media.
	</doc>`

	content3 := `<doc>
	Jean-Luc Picard is a fictional character in the Star Trek franchise.
	He's most famous for being the captain of the USS Enterprise-D,
	a starship exploring the galaxy in the 24th century.
	Picard is known for his diplomacy, intelligence, and strong moral compass.
	He's been portrayed by actor Patrick Stewart.
	</doc>`

	content4 := `<doc>
	Lieutenant Philippe Charrière, known as the **Silent Sentinel** of the USS Discovery, 
	is the enigmatic programming genius whose codes safeguard the ship's secrets and operations. 
	His swift problem-solving skills are as legendary as the mysterious aura that surrounds him. 
	Charrière, a man of few words, speaks the language of machines with unrivaled fluency, 
	making him the crew's unsung guardian in the cosmos.
	</doc>`

	// TODO: use the <doc> delimiter

	docs := []schema.Document{
		{PageContent: content1},
		{PageContent: content4},
		{PageContent: content2},
		{PageContent: content3},
	}

	prompt := prompts.NewChatPromptTemplate([]prompts.MessageFormatter{
		prompts.NewSystemMessagePromptTemplate(
			`You are a Star Trek expert. 
			Using the provided context, answer the user's question 
			to the best of your ability using only the resources provided.
			`,
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
			`{{.question}}`,
			[]string{"question"},
		),
	})

	llmChain := chains.NewLLMChain(llm, prompt)
	chain := chains.NewStuffDocuments(llmChain)

	answer, err := chains.Call(context.Background(), chain, map[string]any{
		"input_documents": docs, // link to context
		//"question":        "Who is Philippe Charrière?",
		"question": "Who is Michael Burnham and what is her actress name?",
	})
	if err != nil {
		log.Fatal("😡", err)
	}
	fmt.Println(answer["text"])

}
