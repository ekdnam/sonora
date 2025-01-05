package utils

import (
	"encoding/json"
	"fmt"
	TypeLeo "leo/src/typeLeo"
	"os"
	"strings"

	"strconv"

	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
)

// LoadConfig loads environment variables from a .env file and returns the value for a specific key.
//
// This function attempts to load environment variables from a file specified by the path parameter
// using godotenv. If the file cannot be loaded, the function will log a fatal error and terminate
// the program. After loading, it retrieves and returns the value associated with the specified key_name.
//
// Parameters:
//   - path: The filesystem path to the .env file (e.g., ".env" or "/path/to/.env")
//   - key_name: The name of the environment variable to retrieve
//
// Returns:
//   - string: The value of the environment variable specified by key_name
//
// Example usage:
//
//	apiKey := LoadConfig(".env", "API_KEY")
//
// Note: This function will terminate the program if the .env file cannot be loaded.
func LoadConfig(path string, key_name string) (string, error) {
	err := godotenv.Load(path)
	if err != nil {
		return "", err
	}
	return os.Getenv(key_name), nil
}

// ValidateConfig ensures that the provided GenerativeModelConfig parameters are within acceptable ranges.
//
// This function validates temperature, topP, topK, and maxOutputTokens against predefined constants.
// Each parameter must fall within its respective min/max range to be considered valid.
//
// Args:
//
//	config: A GenerativeModelConfig struct containing the model parameters to validate.
//
// Returns:
//
//	error: nil if all parameters are valid, or an error describing which parameter is out of range.
func ValidateConfig(config TypeLeo.GenerativeModelConfig) error {
	if config.Temperature < TypeLeo.Constants.MinTemperature || config.Temperature > TypeLeo.Constants.MaxTemperature {
		return fmt.Errorf("temperature must be between %.1f and %.1f", TypeLeo.Constants.MinTemperature, TypeLeo.Constants.MaxTemperature)
	}
	if config.TopP < TypeLeo.Constants.MinTopP || config.TopP > TypeLeo.Constants.MaxTopP {
		return fmt.Errorf("topP must be between %.1f and %.1f", TypeLeo.Constants.MinTopP, TypeLeo.Constants.MaxTopP)
	}
	if config.TopK < TypeLeo.Constants.MinTopK || config.TopK > TypeLeo.Constants.MaxTopK {
		return fmt.Errorf("topK must be between %d and %d", TypeLeo.Constants.MinTopK, TypeLeo.Constants.MaxTopK)
	}
	if config.MaxOutputTokens < TypeLeo.Constants.MinMaxOutputTokens || config.MaxOutputTokens > TypeLeo.Constants.MaxMaxOutputTokens {
		return fmt.Errorf("maxOutputTokens must be between %d and %d", TypeLeo.Constants.MinMaxOutputTokens, TypeLeo.Constants.MaxMaxOutputTokens)
	}
	return nil
}

// PrintResponse prints the content of a GenerateContentResponse to standard output.
//
// Iterates through all candidates in the response and prints their content parts.
// A separator ("---") is printed after all content has been output.
//
// Args:
//
//	resp: A pointer to a GenerateContentResponse containing the content to print.
func PrintResponse(resp *genai.GenerateContentResponse) {
	outputs := ConvertFromResponseToString(resp)
	for idx, output := range outputs {
		fmt.Printf("[%d] %s\n", idx, output)
		fmt.Println("---")
	}
}

// IsAllowedModel checks if a given model name is in the list of allowed models.
//
// Args:
//
//	modelName: The name of the model to check.
//
// Returns:
//
//	bool: true if the model is allowed, false otherwise.
func IsAllowedModel(modelName string) bool {
	return TypeLeo.Constants.AllowedModels[modelName]
}

// IsAllowedLevel checks if a given difficulty level is in the list of allowed levels.
//
// Args:
//
//	level: The difficulty level to check.
//
// Returns:
//
//	bool: true if the level is allowed, false otherwise.
func IsAllowedLevel(level string) bool {
	return TypeLeo.Constants.AllowedLevels[level]
}

// GetAllModels returns a slice containing all allowed model names.
//
// Returns:
//
//	[]string: A slice containing all allowed model names from the Constants.AllowedModels map.
func GetAllModels() []string {
	models := make([]string, 0, len(TypeLeo.Constants.AllowedModels))
	for model := range TypeLeo.Constants.AllowedModels {
		models = append(models, model)
	}
	return models
}

// GetAllLevels returns a slice containing all allowed difficulty levels.
//
// Returns:
//
//	[]string: A slice containing all allowed difficulty levels from the Constants.AllowedLevels map.
func GetAllLevels() []string {
	levels := make([]string, 0, len(TypeLeo.Constants.AllowedLevels))
	for level := range TypeLeo.Constants.AllowedLevels {
		levels = append(levels, level)
	}
	return levels
}

// ConvertFromStringToType converts a string to the specified data type.
//
// Args:
//   - content: The string to convert
//   - datatype: The target data type ("string", "int", "float", "bool")
//
// Returns:
//   - interface{}: The converted value
//   - error: An error if the conversion fails
func ConvertFromStringToType(content string, datatype string) (interface{}, error) {
	switch datatype {
	case "string":
		return content, nil
	case "int":
		return strconv.Atoi(content)
	case "float":
		return strconv.ParseFloat(content, 64)
	case "bool":
		return strconv.ParseBool(content)
	default:
		return nil, fmt.Errorf("unsupported data type: %s", datatype)
	}
}

// ConvertFromResponseToString converts a Gemini API response into a slice of strings.
//
// This utility function extracts the text content from a GenerateContentResponse object,
// which contains candidates with potentially multiple content parts. It flattens the
// hierarchical response structure into a simple string slice for easier consumption
// by client code.
//
// Args:
//   - resp: A pointer to genai.GenerateContentResponse containing the model's response
//
// Returns:
//   - []string: A slice containing the string representation of each content part
//     from all valid candidates. Returns an empty slice if the response is invalid
//     or contains no content.
//
// Example usage:
//
//	resp, _ := model.GenerateContent(ctx, prompt)
//	textParts := ConvertFromResponseToString(resp)
//	for _, text := range textParts {
//	    fmt.Println(text)
//	}
//
// Note: This function safely handles nil checks for both candidates and their content,
// making it robust for production use. The output maintains the original order of
// candidates and their parts.
func ConvertFromResponseToString(resp *genai.GenerateContentResponse) []string {
	if resp == nil {
		return []string{}
	}
	output := []string{}
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			var current string
			for _, part := range cand.Content.Parts {
				current += fmt.Sprint(part) + " "
			}
			current = strings.TrimSpace(current)
			output = append(output, current)
		}
	}
	return output
}

// ParseStringJsonResponse parses a JSON string into a typed struct.
//
// This generic function provides type-safe JSON deserialization for any struct type T.
// It takes a JSON string as input and returns a pointer to the deserialized struct
// of type T. The function uses Go's built-in json.Unmarshal under the hood and provides
// proper error handling.
//
// Type Parameters:
//   - T: The target struct type to deserialize the JSON into. Must satisfy the empty
//     interface (any).
//
// Args:
//   - jsonStr: A string containing valid JSON data that matches the structure of T.
//
// Returns:
//   - *T: A pointer to the deserialized struct of type T.
//   - error: An error if JSON parsing fails, or nil on success.
//
// Example usage:
//
//	type Person struct {
//	    Name string `json:"name"`
//	    Age  int    `json:"age"`
//	}
//
//	jsonStr := `{"name": "Alice", "age": 30}`
//	person, err := ParseStringJsonResponse[Person](jsonStr)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Name: %s, Age: %d\n", person.Name, person.Age)
//
// Note: This function assumes the input JSON string is well-formed and matches
// the target type T. If the JSON structure doesn't match T, an error will be returned.
func ParseStringJsonResponse[T any](jsonStr string) (*T, error) {
	var response T
	err := json.Unmarshal([]byte(jsonStr), &response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %v", err)
	}
	return &response, nil
}

// ParseJsonArrayResponse parses a JSON string into a slice of typed structs.
//
// This generic function provides type-safe JSON deserialization for arrays/slices of any struct type T.
// It takes a JSON string representing an array as input and returns a slice of the deserialized type T.
// The function uses Go's built-in json.Unmarshal under the hood and provides proper error handling.
//
// Type Parameters:
//   - T: The target struct type to deserialize each JSON array element into. Must satisfy the empty
//     interface (any).
//
// Args:
//   - jsonStr: A string containing valid JSON array data where each element matches the structure of T.
//
// Returns:
//   - []T: A slice containing the deserialized structs of type T.
//   - error: An error if JSON parsing fails, or nil on success.
//
// Example usage:
//
//	type Topic struct {
//	    Name string `json:"name"`
//	    Description string `json:"description"`
//	}
//
//	jsonStr := `[{"name": "Physics", "description": "Study of matter"},
//	             {"name": "Chemistry", "description": "Study of substances"}]`
//	topics, err := ParseJsonArrayResponse[Topic](jsonStr)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for _, topic := range topics {
//	    fmt.Printf("Topic: %s - %s\n", topic.Name, topic.Description)
//	}
//
// Note: This function assumes the input JSON string represents a well-formed array where each element
// matches the target type T. If the JSON structure doesn't match, an error will be returned.
func ParseJsonArrayResponse[T any](jsonStr string) ([]T, error) {
	var response []T
	err := json.Unmarshal([]byte(jsonStr), &response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON array response: %v", err)
	}
	return response, nil
}
