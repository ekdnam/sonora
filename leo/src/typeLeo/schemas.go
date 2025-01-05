package TypeLeo

import "github.com/google/generative-ai-go/genai"

// ValidateTopicSchema defines the schema for validating if a topic is suitable for course creation
var ValidateTopicSchema = &genai.Schema{
	Type: genai.TypeBoolean,
}

// AlternateTopicSingletonObject defines the properties for a single alternative topic suggestion
var AlternateTopicSingletonObject = map[string]*genai.Schema{
	"id":      {Type: genai.TypeInteger},
	"subject": {Type: genai.TypeString},
}

// AlternateTopicSingletonSchema wraps the properties into a single object schema
var AlternateTopicSingletonSchema = &genai.Schema{
	Type:       genai.TypeObject,
	Properties: AlternateTopicSingletonObject,
}

// AlternateTopicsArraySchema defines an array of alternative topic suggestions
var AlternateTopicsArraySchema = &genai.Schema{
	Type:  genai.TypeArray,
	Items: AlternateTopicSingletonSchema,
}
