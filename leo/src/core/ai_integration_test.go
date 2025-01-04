//go:build integration
// +build integration

package core

import (
	"context"
	"leo/src/utils"
	"testing"
	"time"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func TestGenerateContent_Integration(t *testing.T) {
	// Skip if API key not set
	apiKey, err := utils.LoadConfig("../../.env", "GEMINI_API_KEY")
	if err != nil {
		t.Errorf("GEMINI_API_KEY not set, skipping integration test. Error: %v", err)
	}
	if apiKey == "" {
		t.Errorf("GEMINI_API_KEY not set, skipping integration test")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-2.0-flash-exp")

	t.Run("real API call", func(t *testing.T) {
		// Add timeout to context
		ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		prompt := "What is 2+2?"
		resp, err := GenerateContent(ctx, model, prompt)
		t.Logf("Response: %v", resp)
		if err != nil {
			t.Errorf("Failed to generate content: %v", err)
		}

		if resp == nil {
			t.Fatal("Expected response, got nil")
		}

		if len(resp.Candidates) == 0 {
			t.Error("Expected at least one candidate in response")
		}

		// Print response for manual inspection
		if len(resp.Candidates) > 0 && resp.Candidates[0].Content != nil {
			t.Logf("Response: %v", resp.Candidates[0].Content.Parts[0])
		}
	})

	t.Run("with invalid prompt", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		// Empty prompt should cause an error
		resp, err := GenerateContent(ctx, model, "")
		t.Logf("Response: %v", resp)
		if err == nil {
			t.Error("Expected error for empty prompt, got nil")
		}

		if resp != nil {
			t.Errorf("Expected nil response for error case, got %v", resp)
		}
	})

	t.Run("stream response", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		prompt := "Count from 1 to 5 slowly"
		resp, err := GenerateContentStream(ctx, model, prompt)

		// Since GenerateContentStream currently returns nil, nil when successful
		if err != nil {
			t.Errorf("Failed to generate content stream: %v", err)
		}
		if resp != nil {
			t.Errorf("Expected nil response for streaming case, got %v", resp)
		}
	})
}
