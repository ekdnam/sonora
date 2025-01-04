package main

import (
	"context"
	"fmt"
	"leo/src/core"
	"leo/src/utils"
	"log"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func main() {

	ctx := context.Background()
	apiKey, err := utils.LoadConfig(".env", "GEMINI_API_KEY")
	if err != nil {
		log.Fatal(err)
	}
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	model, err := core.GetModel(client, core.GenerativeModelConfig{
		ModelName:       "gemini-2.0-flash-thinking-exp-1219",
		Temperature:     0.5,
		TopP:            0.95,
		TopK:            40,
		MaxOutputTokens: 8192,
	})
	if err != nil {
		log.Fatal(err)
	}
	topic := "cosmology"
	level := "advanced"
	resp, err := core.GeneratePlan(ctx, model, topic, level)
	if err != nil {
		log.Fatal(err)
	}
	printResponse(resp)

}

func printResponse(resp *genai.GenerateContentResponse) {
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				fmt.Println(part)
			}
		}
	}
	fmt.Println("---")
}
