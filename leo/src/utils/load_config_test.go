package utils

import (
	"log"
	"os"
	"testing"
)

func createTempEnvFile(t *testing.T) (string, func()) {
	content := []byte("GEMINI_API_KEY=test_key\nPORT=3000\n")
	tmpfile, err := os.CreateTemp("", "test.env")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := tmpfile.Write(content); err != nil {
		t.Fatal(err)
	}
	return tmpfile.Name(), func() { os.Remove(tmpfile.Name()) }
}

func TestLoadConfig_GeminiKey(t *testing.T) {
	envFile, cleanup := createTempEnvFile(t)
	defer cleanup()
	apiKey, err := LoadConfig(envFile, "GEMINI_API_KEY")
	if err != nil {
		t.Fatal(err)
	}
	log.Println(apiKey)
}

func TestLoadConfig_FileNotFound(t *testing.T) {
	apiKey, err := LoadConfig(".env.nonexistent", "GEMINI_API_KEY")
	if err == nil {
		t.Fatal("Expected error for non-existent file, got nil")
	}
	if apiKey != "" {
		t.Errorf("Expected empty string for non-existent file, got %q", apiKey)
	}
}

func TestLoadConfig_APIKeyNotFound(t *testing.T) {
	envFile, cleanup := createTempEnvFile(t)
	defer cleanup()
	apiKey, err := LoadConfig(envFile, "RANDOM_KEY")
	if err != nil {
		t.Fatal("Expected no error, since file exists but key is not present")
	}
	if apiKey != "" {
		t.Fatalf("Expected empty string for non-existent file, got %q", apiKey)
	}
}
