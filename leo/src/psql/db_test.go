package psql

import (
	"testing"
)

func TestInitDB(t *testing.T) {
	// Create a temporary .env file for testing
	envPath := "../../.env"

	// Test successful connection
	t.Run("Successful connection", func(t *testing.T) {
		db, err := InitDB(envPath)
		if err != nil {
			t.Errorf("InitDB failed: %v", err)
		}
		if db != nil {
			defer db.Close()
		}
	})

	// Test non-existent env file
	t.Run("Non-existent env file", func(t *testing.T) {
		_, err := InitDB("nonexistent.env")
		if err == nil {
			t.Error("Expected error for non-existent env file, got nil")
		}
	})
}
