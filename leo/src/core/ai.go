package core

import (
	"context"
	"fmt"
	"leo/src/prompts"
	"leo/src/utils"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/iterator"
)

// GenerativeModelConfig defines the configuration parameters for a generative AI model.
//
// This struct encapsulates all the necessary parameters to configure and fine-tune
// the behavior of a large language model for text generation tasks. It provides
// control over both the model selection and its generation parameters.
//
// Fields:
//   - ModelName: The identifier of the model to use (e.g., "gemini-1.5-flash")
//   - Temperature: Controls randomness in generation (0.0-1.0). Higher values increase
//     creativity while lower values make output more focused and deterministic
//   - TopP: Nucleus sampling parameter (0.0-1.0). Filters out tokens below cumulative
//     probability threshold, helping balance diversity and quality
//   - TopK: Limits vocabulary to K most likely tokens. Lower values (e.g., 40)
//     increase output reliability while higher values allow more diversity
//   - MaxOutputTokens: Maximum length of generated response in tokens
//   - SystemMessage: Optional context or instructions to guide model behavior
//
// Example usage:
//
//	config := GenerativeModelConfig{
//	    ModelName: "gemini-1.5-flash",
//	    Temperature: 0.7,
//	    TopP: 0.95,
//	    TopK: 40,
//	    MaxOutputTokens: 1024,
//	    SystemMessage: "You are a helpful assistant",
//	}
//
// Note: The optimal values for temperature, topP, and topK depend on your specific
// use case. Lower values generally produce more focused and deterministic output,
// while higher values allow for more creativity and variation.
type GenerativeModelConfig struct {
	ModelName       string
	Temperature     float32
	TopP            float32
	TopK            int32
	MaxOutputTokens int32
}

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
func GetModel(client *genai.Client, config GenerativeModelConfig) (*genai.GenerativeModel, error) {
	if !Constants.IsAllowedModel(config.ModelName) {
		return nil, fmt.Errorf("invalid model name: %s. Allowed models: %v",
			config.ModelName, Constants.GetAllModels())
	}
	// Add validation before creating model
	if config.Temperature < Constants.MinTemperature || config.Temperature > Constants.MaxTemperature {
		return nil, fmt.Errorf("temperature must be between %.1f and %.1f", Constants.MinTemperature, Constants.MaxTemperature)
	}
	if config.TopP < Constants.MinTopP || config.TopP > Constants.MaxTopP {
		return nil, fmt.Errorf("topP must be between %.1f and %.1f", Constants.MinTopP, Constants.MaxTopP)
	}
	if config.TopK < Constants.MinTopK || config.TopK > Constants.MaxTopK {
		return nil, fmt.Errorf("topK must be between %d and %d", Constants.MinTopK, Constants.MaxTopK)
	}
	if config.MaxOutputTokens < Constants.MinMaxOutputTokens || config.MaxOutputTokens > Constants.MaxMaxOutputTokens {
		return nil, fmt.Errorf("maxOutputTokens must be between %d and %d", Constants.MinMaxOutputTokens, Constants.MaxMaxOutputTokens)
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
	if !Constants.IsAllowedLevel(level) {
		return nil, fmt.Errorf("invalid level: %s. Allowed levels: %v",
			level, Constants.GetAllLevels())
	}
	prompt := fmt.Sprintf("Topic: %s\nLevel: %s", topic, level)
	return GenerateContent(ctx, model, prompt)
}
