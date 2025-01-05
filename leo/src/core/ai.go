package core

import (
	"context"
	"fmt"
	"leo/src/prompts"
	TypeLeo "leo/src/typeLeo"
	"leo/src/utils"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/iterator"
)

// GetModel initializes and configures a generative AI model with the specified parameters.
//
// This function creates a new generative model instance with the provided configuration,
// performing validation and applying the specified generation parameters. It's a critical
// component for ensuring consistent and valid model initialization across the application.
//
// Args:
//   - client: A initialized Gemini API client instance
//   - config: A GenerativeModelConfig struct containing model parameters and settings
//
// Returns:
//   - *genai.GenerativeModel: A configured generative model instance ready for use
//   - error: An error if model validation fails or configuration is invalid
//
// Example usage:
//
//	model, err := GetModel(client, GenerativeModelConfig{
//	    ModelName: "gemini-1.5-pro",
//	    Temperature: 0.7,
//	    TopP: 0.95,
//	    TopK: 40,
//	    MaxOutputTokens: 1024,
//	})
//
// The function performs the following key operations:
// 1. Validates that the requested model name is in the allowed models list
// 2. Creates a new model instance with the specified name
// 3. Configures generation parameters (temperature, topP, topK, max tokens)
//
// Important: This function should be used as the standard way to obtain model
// instances to ensure consistent configuration and validation across the codebase.
// Direct model creation should be avoided to maintain reliability.
func GetModel(client *genai.Client, config TypeLeo.GenerativeModelConfig) (*genai.GenerativeModel, error) {
	if !utils.IsAllowedModel(config.ModelName) {
		return nil, fmt.Errorf("invalid model name: %s. Allowed models: %v",
			config.ModelName, utils.GetAllModels())
	}
	// Add validation before creating model
	err := utils.ValidateConfig(config)
	if err != nil {
		return nil, err
	}
	model := client.GenerativeModel(config.ModelName)
	model.SetTemperature(config.Temperature)
	model.SetTopP(config.TopP)
	model.SetTopK(config.TopK)
	model.SetMaxOutputTokens(config.MaxOutputTokens)
	return model, nil
}

// GenerateContent generates AI content using the Gemini model with the provided prompt.
//
// This function serves as the core content generation interface for the Gemini API,
// handling the direct interaction with the model while providing proper error handling
// and context management. It's designed to be a reliable foundation for all content
// generation operations across the application.
//
// Args:
//   - ctx: Context for managing request lifecycle, timeouts, and cancellation
//   - model: A properly configured Gemini model instance (see GetModel)
//   - prompt: The text prompt to send to the model for content generation
//
// Returns:
//   - *genai.GenerateContentResponse: The model's generated response containing candidates
//   - error: Any errors encountered during generation, including API errors
//
// Example usage:
//
//	resp, err := GenerateContent(ctx, model, "Explain quantum computing")
//	if err != nil {
//	    log.Fatalf("Content generation failed: %v", err)
//	}
//
// Important implementation notes:
// - The function uses genai.Text() to properly format the prompt for the API
// - Errors are propagated without modification to allow proper handling upstream
// - The context should include appropriate timeouts for production use
//
// Thread safety: This function is safe for concurrent use as the underlying
// Gemini client handles request synchronization.
func GenerateContent(ctx context.Context, model *genai.GenerativeModel, prompt string) (*genai.GenerateContentResponse, error) {
	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func GenerateContentStream(ctx context.Context, model *genai.GenerativeModel, prompt string) (*genai.GenerateContentResponse, error) {
	iter := model.GenerateContentStream(ctx, genai.Text(prompt))
	for {
		resp, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		utils.PrintResponse(resp)
	}
	return nil, nil
}

// GeneratePlan creates a structured learning plan using a generative AI model.
//
// This function leverages a generative model to create a customized learning plan
// based on a specified topic and difficulty level. It first sets up the model with
// appropriate system instructions, validates the input parameters, and then generates
// the content.
//
// Args:
//   - ctx: Context for the request, used for cancellation and timeouts
//   - model: A configured generative model instance to use for content generation
//   - topic: The subject matter for which to generate a learning plan
//   - level: The difficulty level of the plan (must be one of the allowed levels)
//
// Returns:
//   - *genai.GenerateContentResponse: The generated learning plan response
//   - error: An error if level validation fails or content generation fails
//
// Example usage:
//
//	resp, err := GeneratePlan(ctx, model, "Machine Learning", "intermediate")
//	if err != nil {
//	    log.Fatalf("Failed to generate plan: %v", err)
//	}
//
// Note: This function performs validation on the level parameter using the core.IsAllowedLevel
// function to ensure only supported difficulty levels are used. The system message used for
// generation is defined in prompts.GeneratePlanPrompt.
func GeneratePlan(ctx context.Context, model *genai.GenerativeModel, topic string, level string) (*genai.GenerateContentResponse, error) {
	systemMessage := prompts.GeneratePlanPrompt
	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(systemMessage)},
		Role:  "user",
	}
	if !utils.IsAllowedLevel(level) {
		return nil, fmt.Errorf("invalid level: %s. Allowed levels: %v",
			level, utils.GetAllLevels())
	}
	prompt := fmt.Sprintf("Topic: %s\nLevel: %s", topic, level)
	return GenerateContent(ctx, model, prompt)
}

// GetIfCourseCanBeMade evaluates whether a given topic is suitable for a university-level STEM course.
//
// This function leverages a generative AI model to analyze a topic and determine if it meets
// two key criteria:
// 1. The topic must be STEM-related
// 2. The topic must be specific enough to form the basis of a complete university course
//
// The function configures the model to return a structured JSON response with a boolean
// indicating whether the topic is suitable. The response schema is strictly enforced to
// ensure consistent output formatting.
//
// Args:
//   - ctx: Context for managing timeouts and cancellation
//   - model: A configured genai.GenerativeModel instance
//   - topic: The subject matter to evaluate (e.g., "Machine Learning", "CPU Architecture")
//
// Returns:
//   - *genai.GenerateContentResponse: Contains a JSON response in the format {"response": bool}
//   - error: Any error encountered during the evaluation process
//
// Example usage:
//
//	resp, err := GetIfCourseCanBeMade(ctx, model, "Large Language Models")
//	if err != nil {
//	    log.Fatalf("Failed to evaluate topic: %v", err)
//	}
//
// Note: The function uses a predefined system prompt (DetermineIfCourseCanBeMadeOnTopicPrompt)
// that contains specific evaluation criteria and example responses. The model is configured
// to return JSON to ensure consistent parsing of results.
func ValidateTopic(ctx context.Context, model *genai.GenerativeModel, topic string) (*genai.GenerateContentResponse, error) {
	systemMessage := prompts.ValidateTopicPrompt
	model.ResponseMIMEType = "application/json"
	model.ResponseSchema = TypeLeo.ValidateTopicSchema
	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(systemMessage)},
		Role:  "user",
	}
	prompt := fmt.Sprintf("Topic: %s\n", topic)
	return GenerateContent(ctx, model, prompt)
}

func RecommendAlternateTopics(ctx context.Context, model *genai.GenerativeModel, topic string) (*genai.GenerateContentResponse, error) {
	systemMessage := prompts.RecommendAlternateTopicsPrompt
	model.ResponseMIMEType = "application/json"
	model.ResponseSchema = TypeLeo.AlternateTopicsArraySchema
	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(systemMessage)},
		Role:  "user",
	}
	prompt := fmt.Sprintf("Topic: %s\n", topic)
	return GenerateContent(ctx, model, prompt)
}
