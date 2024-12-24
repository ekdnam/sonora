package main

import (
	"context"
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

}
