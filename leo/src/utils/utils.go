package utils

import (
	"fmt"
	"os"

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

func PrintResponse(resp *genai.GenerateContentResponse) {
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				fmt.Println(part)
			}
		}
	}
	fmt.Println("---")
}
