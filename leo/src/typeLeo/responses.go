package TypeLeo

// ValidateTopicResponse represents the structure of the response from ValidateTopic
type ValidateTopicResponse struct {
	IsValid bool   `json:"is_valid"`
	Reason  string `json:"reason"`
}

// AlternateTopicSuggestion represents a single alternative topic suggestion
type AlternateTopicSuggestionResponse struct {
	ID      int    `json:"id"`
	Subject string `json:"subject"`
}

// AlternateTopicsResponse represents an array of alternative topic suggestions
type AlternateTopicSuggestionArrayResponse []AlternateTopicSuggestionResponse
