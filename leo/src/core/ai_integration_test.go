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
		state *TypeLeo.ValidateTopicResponse
	}{
		{"computers", &TypeLeo.ValidateTopicResponse{IsValid: false, Reason: "Reason 1"}},
		{"thermodynamics", &TypeLeo.ValidateTopicResponse{IsValid: true, Reason: "Reason 2"}},
		{"physics", &TypeLeo.ValidateTopicResponse{IsValid: false, Reason: "Reason 3"}},
		{"nonsense_1", &TypeLeo.ValidateTopicResponse{IsValid: false, Reason: "Reason 4"}},
		{"nonsense_2", &TypeLeo.ValidateTopicResponse{IsValid: false, Reason: "Reason 5"}},
	}

	for i, tc := range testCases {
		addTestDelay(t, i == 0) // Skip delay for first test

		t.Run("topic: "+tc.topic, func(t *testing.T) {
			t.Logf("Testing topic: %s", tc.topic)
			model, ctx, cancel := setupTestModel(t)
			defer cancel()
			resp, err := ValidateTopic(ctx, model, tc.topic)
			t.Logf("Topic: %s, Response: %v, Reason: %s", tc.topic, resp.IsValid, resp.Reason)
			if err != nil {
				t.Errorf("Failed to generate content schema: %v", err)
			}
			if resp == nil {
				t.Errorf("Expected non-nil response, got nil")
			}
			if resp.IsValid != tc.state.IsValid {
				t.Errorf("Expected response %v, got %v", tc.state.IsValid, resp.IsValid)
			}
			t.Logf("Completed testing topic: %s", tc.topic)
		})
	}
}

func TestRecommendAlternateTopics_Integration(t *testing.T) {
	testCases := []struct {
		name          string
		topic         string
		expectedCount int
		validateFn    func([]TypeLeo.AlternateTopicSuggestionResponse) error
	}{
		{
			name:          "broad topic - computers",
			topic:         "computers",
			expectedCount: 3,
			validateFn: func(suggestions []TypeLeo.AlternateTopicSuggestionResponse) error {
				// Validate that IDs are 1-3
				for i, suggestion := range suggestions {
					if suggestion.ID != i+1 {
						return fmt.Errorf("expected ID %d, got %d", i+1, suggestion.ID)
					}
					if suggestion.Subject == "" {
						return fmt.Errorf("empty subject string found at index %d", i)
					}
				}
				return nil
			},
		},
		{
			name:          "specific topic - quantum computing",
			topic:         "quantum computing",
			expectedCount: 3,
			validateFn: func(suggestions []TypeLeo.AlternateTopicSuggestionResponse) error {
				seen := make(map[string]bool)
				for _, suggestion := range suggestions {
					if suggestion.Subject == "" {
						return fmt.Errorf("empty topic string found")
					}
					// Check for duplicates
					if seen[suggestion.Subject] {
						return fmt.Errorf("duplicate topic found: %s", suggestion.Subject)
					}
					seen[suggestion.Subject] = true
				}
				return nil
			},
		},
		{
			name:          "nonsense topic",
			topic:         "xyzabc123",
			expectedCount: 3,
			validateFn: func(suggestions []TypeLeo.AlternateTopicSuggestionResponse) error {
				// Even for nonsense topics, we should get valid STEM suggestions
				for _, suggestion := range suggestions {
					if suggestion.Subject == "" {
						return fmt.Errorf("empty topic string found")
					}
					if suggestion.Subject == "xyzabc123" {
						return fmt.Errorf("model returned the nonsense topic")
					}
				}
				return nil
			},
		},
	}

	for i, tc := range testCases {
		addTestDelay(t, i == 0) // Skip delay for first test

		t.Run(tc.name, func(t *testing.T) {
			t.Logf("Testing topic: %s", tc.topic)
			model, ctx, cancel := setupTestModel(t)
			defer cancel()

			suggestions, err := RecommendAlternateTopics(ctx, model, tc.topic)

			// Log the response for debugging
			t.Logf("Topic: %s, Suggestions: %+v", tc.topic, suggestions)

			if err != nil {
				t.Fatalf("Failed to get alternate topics: %v", err)
			}

			// Validate response length
			if len(suggestions) != tc.expectedCount {
				t.Errorf("Expected %d suggestions, got %d", tc.expectedCount, len(suggestions))
			}

			// Run custom validation function
			if err := tc.validateFn(suggestions); err != nil {
				t.Errorf("Validation failed: %v", err)
			}

			t.Logf("Completed testing topic: %s", tc.topic)
		})
	}
}
