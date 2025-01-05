package TypeLeo

// ValidateTopicResponse represents the structure of the response from ValidateTopic
type ValidateTopicResponse struct {
	IsValid bool   `json:"is_valid"`
	Reason  string `json:"reason"`
}
