package utils

import (
	"os"

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

// IsAllowedModel checks if a given model name is in the allowed list of Gemini models.
//
// This function validates whether a model name is officially supported by checking against
// a predefined allowlist of Gemini model versions. The allowlist is maintained to ensure
// only stable and supported model versions are used in production.
//
// Parameters:
//   - modelName: The name of the Gemini model to validate (e.g., "gemini-1.5-pro")
//
// Returns:
//   - bool: true if the model is in the allowed list, false otherwise
//
// Example usage:
//
//	if !IsAllowedModel(modelName) {
//	    return fmt.Errorf("unsupported model: %s", modelName)
//	}
//
// Note: The list of allowed models is periodically updated as new versions are released
// and older versions are deprecated. Always use this validation before making API calls.
func IsAllowedModel(modelName string) bool {
	return allowedModels[modelName]
}

var allowedModels = map[string]bool{
	"gemini-2.0-flash-exp":               true,
	"gemini-1.5-flash":                   true,
	"gemini-1.5-flash-8b":                true,
	"gemini-1.5-pro":                     true,
	"gemini-2.0-flash-thinking-exp-1219": true,
}

// GetAllowedModels returns a slice of all supported Gemini model names.
//
// This function provides access to the complete list of officially supported Gemini models
// that are allowed for use in production environments. It converts the internal allowlist
// map into a slice of model names for easier iteration and display purposes.
//
// Returns:
//   - []string: A slice containing all allowed model names
//
// Example usage:
//
//	models := GetAllowedModels()
//	for _, model := range models {
//	    fmt.Printf("Supported model: %s\n", model)
//	}
//
// Note: The returned slice is generated dynamically from the allowedModels map,
// ensuring it always reflects the current set of supported models. The order of
// models in the returned slice is not guaranteed to be stable between calls.

func GetAllowedModels() []string {
	models := make([]string, 0, len(allowedModels))
	for model := range allowedModels {
		models = append(models, model+"\n")
	}
	return models
}

// IsAllowedLevel checks if a given course difficulty level is supported by the system.
//
// This function validates whether a specified course difficulty level is among the
// officially supported levels in the platform. It performs a simple lookup in the
// allowedLevels map to determine if the level is valid.
//
// Args:
//   - level: A string representing the course difficulty level to validate
//
// Returns:
//   - bool: true if the level is in the allowed list, false otherwise
//
// Example usage:
//
//	if !IsAllowedLevel(userLevel) {
//	    return fmt.Errorf("unsupported difficulty level: %s", userLevel)
//	}
//
// Note: The supported levels are fixed to beginner, intermediate, and advanced to
// maintain consistent course categorization across the platform.
func IsAllowedLevel(level string) bool {
	return allowedLevels[level]
}

// allowedLevels defines the set of valid course difficulty levels supported by the system.
// This map acts as a source of truth for level validation throughout the application.
var allowedLevels = map[string]bool{
	"beginner":     true,
	"intermediate": true,
	"advanced":     true,
}

// GetAllowedLevels returns a slice of all supported course difficulty levels.
//
// This function provides access to the complete list of officially supported
// difficulty levels that can be assigned to courses. It converts the internal
// allowlist map into a slice of level names for easier iteration and display
// purposes.
//
// Returns:
//   - []string: A slice containing all allowed difficulty level names
//
// Example usage:
//
//	levels := GetAllowedLevels()
//	for _, level := range levels {
//	    fmt.Printf("Supported level: %s\n", level)
//	}
//
// Note: The returned slice is generated dynamically from the allowedLevels map,
// ensuring it always reflects the current set of supported levels. The order of
// levels in the returned slice is not guaranteed to be stable between calls.
func GetAllowedLevels() []string {
	levels := make([]string, 0, len(allowedLevels))
	for level := range allowedLevels {
		levels = append(levels, level+"\n")
	}
	return levels
}
