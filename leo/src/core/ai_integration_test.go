//go:build integration
// +build integration

package core

import (
	"context"
	TypeLeo "leo/src/typeLeo"
	"leo/src/utils"
	"testing"
	"time"

	"fmt"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

var DELAY_TIME = 120

// setupTestModel creates and returns a configured test model
func setupTestModel(t *testing.T) (*genai.GenerativeModel, context.Context, context.CancelFunc) {
	apiKey, err := utils.LoadConfig("../../.env", "GEMINI_API_KEY")
	if err != nil {
		t.Fatalf("GEMINI_API_KEY not set, skipping integration test. Error: %v", err)
	}
	if apiKey == "" {
		t.Fatal("GEMINI_API_KEY not set, skipping integration test")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	model, err := GetModel(client, TypeLeo.GenerativeModelConfig{
		ModelName:       "gemini-1.5-pro",
		Temperature:     0.7,
		TopP:            0.95,
		TopK:            40,
		MaxOutputTokens: 1024,
	})
	if err != nil {
		t.Fatalf("Failed to create model: %v", err)
	}
	ctx, cancel := context.WithTimeout(ctx, time.Duration(DELAY_TIME)*time.Second)
	return model, ctx, cancel
}

func addTestDelay(t *testing.T, skipDelay bool) {
	if skipDelay {
		return
	}

	fmt.Printf("\rWaiting for next test: ")
	for remaining := DELAY_TIME; remaining > 0; remaining-- {
		fmt.Printf("\rWaiting for next test: %d seconds...", remaining)
		time.Sleep(1 * time.Second)
	}
	fmt.Printf("\r") // Clear the line
}

func TestGenerateContent_Integration(t *testing.T) {
	model, ctx, cancel := setupTestModel(t)
	defer cancel()

	t.Run("successful generation", func(t *testing.T) {
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
		if len(resp.Candidates) > 0 && resp.Candidates[0].Content != nil {
			t.Logf("Response: %v", resp.Candidates[0].Content.Parts[0])
		}
	})

	addTestDelay(t, false)

	t.Run("with invalid prompt", func(t *testing.T) {
		resp, err := GenerateContent(ctx, model, "")
		if err == nil {
			t.Error("Expected error for empty prompt, got nil")
		}
		if resp != nil {
			t.Errorf("Expected nil response for error case, got %v", resp)
		}
	})
}

func TestGenerateContentStream_Integration(t *testing.T) {
	model, ctx, cancel := setupTestModel(t)
	defer cancel()

	t.Run("stream response", func(t *testing.T) {
		prompt := "Count from 1 to 5 slowly"
		resp, err := GenerateContentStream(ctx, model, prompt)

		if err != nil {
			t.Errorf("Failed to generate content stream: %v", err)
		}
		if resp != nil {
			t.Errorf("Expected nil response for streaming case, got %v", resp)
		}
	})
}

func TestValidateTopic_Integration(t *testing.T) {
	// model, ctx, cancel := setupTestModel(t)
	// defer cancel()

	testCases := []struct {
		topic string
		state bool
	}{
		{"computers", false},
		{"thermodynamics", true},
		{"physics", false},
		{"nonsense_1", false},
		{"nonsense_2", false},
	}

	for i, tc := range testCases {
		addTestDelay(t, i == 0) // Skip delay for first test

		t.Run("topic: "+tc.topic, func(t *testing.T) {
			t.Logf("Testing topic: %s", tc.topic)
			model, ctx, cancel := setupTestModel(t)
			defer cancel()
			resp, err := ValidateTopic(ctx, model, tc.topic)

			if err != nil {
				t.Errorf("Failed to generate content schema: %v", err)
			}
			if resp == nil {
				t.Errorf("Expected non-nil response, got nil")
			}
			stringResponse := utils.ConvertFromResponseToString(resp)
			if len(stringResponse) == 0 {
				t.Errorf("Expected non-empty response, got empty")
			}
			t.Logf("Response: %v", stringResponse)
			boolResponse, err := utils.ConvertFromStringToType(stringResponse[0], "bool")
			if err != nil {
				t.Errorf("Error while converting from string to bool for %s", stringResponse)
			}
			if boolResponse != tc.state {
				t.Errorf("Topic %s : Expected - %v, got %v", tc.topic, tc.state, boolResponse)
			}
			t.Logf("Completed testing topic: %s", tc.topic)
		})
	}
}
