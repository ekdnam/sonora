package TypeLeo

import "github.com/google/generative-ai-go/genai"

// ValidateTopicSchema defines the schema for validating if a topic is suitable for course creation
var ValidateTopicSchema = &genai.Schema{
	Type: genai.TypeObject,
	Properties: map[string]*genai.Schema{
		"is_valid": {Type: genai.TypeBoolean},
		"reason":   {Type: genai.TypeString},
	},
	Required: []string{"is_valid", "reason"},
}

// AlternateTopicSingletonSchema wraps the properties into a single object schema
var AlternateTopicSingletonSchema = &genai.Schema{
	Type: genai.TypeObject,
	Properties: map[string]*genai.Schema{
		"id":      {Type: genai.TypeInteger},
		"subject": {Type: genai.TypeString},
	},
	Required: []string{"id", "subject"},
}

// AlternateTopicsArraySchema defines an array of alternative topic suggestions
var AlternateTopicsArraySchema = &genai.Schema{
	Type:  genai.TypeArray,
	Items: AlternateTopicSingletonSchema,
}
