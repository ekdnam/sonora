package utils

import (
	TypeLeo "leo/src/typeLeo"
	"log"
	"os"
	"testing"

	"github.com/google/generative-ai-go/genai"
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
	log.Printf("Error: %v", err)
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

func TestConvertFromStringToType(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		datatype string
		want     interface{}
		wantErr  bool
	}{
		{"string conversion", "hello", "string", "hello", false},
		{"int conversion", "123", "int", 123, false},
		{"float conversion", "3.14", "float", 3.14, false},
		{"bool conversion true", "true", "bool", true, false},
		{"bool conversion false", "false", "bool", false, false},
		{"invalid int", "abc", "int", nil, true},
		{"invalid float", "xyz", "float", nil, true},
		{"invalid bool", "notbool", "bool", nil, true},
		{"unsupported type", "test", "array", nil, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ConvertFromStringToType(tc.content, tc.datatype)
			if (err != nil) != tc.wantErr {
				t.Errorf("ConvertFromStringToType() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if !tc.wantErr && got != tc.want {
				t.Errorf("ConvertFromStringToType() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestConvertFromResponseToString(t *testing.T) {
	tests := []struct {
		name     string
		response *genai.GenerateContentResponse
		want     []string
	}{
		{
			name: "normal response with single candidate and part",
			response: &genai.GenerateContentResponse{
				Candidates: []*genai.Candidate{
					{
						Content: &genai.Content{
							Parts: []genai.Part{
								genai.Text("Hello, world!"),
							},
						},
					},
				},
			},
			want: []string{"Hello, world!"},
		},
		{
			name: "multiple candidates with multiple parts",
			response: &genai.GenerateContentResponse{
				Candidates: []*genai.Candidate{
					{
						Content: &genai.Content{
							Parts: []genai.Part{
								genai.Text("First part"),
								genai.Text("Second part"),
							},
						},
					},
					{
						Content: &genai.Content{
							Parts: []genai.Part{
								genai.Text("Third part"),
							},
						},
					},
				},
			},
			want: []string{"First part Second part", "Third part"},
		},
		{
			name:     "nil response",
			response: nil,
			want:     []string{},
		},
		{
			name: "empty candidates",
			response: &genai.GenerateContentResponse{
				Candidates: []*genai.Candidate{},
			},
			want: []string{},
		},
		{
			name: "nil content in candidate",
			response: &genai.GenerateContentResponse{
				Candidates: []*genai.Candidate{
					{Content: nil},
				},
			},
			want: []string{},
		},
		{
			name: "mixed valid and nil content",
			response: &genai.GenerateContentResponse{
				Candidates: []*genai.Candidate{
					{Content: nil},
					{
						Content: &genai.Content{
							Parts: []genai.Part{
								genai.Text("Valid content"),
							},
						},
					},
					{Content: nil},
				},
			},
			want: []string{"Valid content"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := ConvertFromResponseToString(tc.response)
			if len(got) != len(tc.want) {
				t.Errorf("ConvertFromResponseToString() got %d elements, want %d elements", len(got), len(tc.want))
				return
			}
			for i := range got {
				if got[i] != tc.want[i] {
					t.Errorf("ConvertFromResponseToString()[%d] = %v, want %v", i, got[i], tc.want[i])
				}
			}
		})
	}
}

func TestParseValidateTopicResponseJSON(t *testing.T) {
	tests := []struct {
		name    string
		jsonStr string
		want    *TypeLeo.ValidateTopicResponse
		wantErr bool
	}{
		{
			name:    "valid response - true",
			jsonStr: `{"is_valid": true, "reason": "Valid topic for a course"}`,
			want: &TypeLeo.ValidateTopicResponse{
				IsValid: true,
				Reason:  "Valid topic for a course",
			},
			wantErr: false,
		},
		{
			name:    "valid response - false",
			jsonStr: `{"is_valid": false, "reason": "Topic too broad"}`,
			want: &TypeLeo.ValidateTopicResponse{
				IsValid: false,
				Reason:  "Topic too broad",
			},
			wantErr: false,
		},
		{
			name:    "invalid json",
			jsonStr: `{invalid json}`,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "missing fields",
			jsonStr: `{"is_valid": true}`,
			want: &TypeLeo.ValidateTopicResponse{
				IsValid: true,
				Reason:  "",
			},
			wantErr: false,
		},
		{
			name:    "empty string",
			jsonStr: "",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ParseStringJsonResponse[TypeLeo.ValidateTopicResponse](tc.jsonStr)
			if (err != nil) != tc.wantErr {
				t.Errorf("ConvertFromStringToValidateTopicResponse() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if tc.wantErr {
				return
			}
			if got.IsValid != tc.want.IsValid {
				t.Errorf("IsValid = %v, want %v", got.IsValid, tc.want.IsValid)
			}
			if got.Reason != tc.want.Reason {
				t.Errorf("Reason = %v, want %v", got.Reason, tc.want.Reason)
			}
		})
	}
}
