package TypeLeo

type GenerativeModelConfig struct {
	ModelName       string
	Temperature     float32
	TopP            float32
	TopK            int32
	MaxOutputTokens int32
}

type constants struct {
	AllowedModels      map[string]bool
	AllowedLevels      map[string]bool
	MinTemperature     float32
	MaxTemperature     float32
	MinTopP            float32
	MaxTopP            float32
	MinTopK            int32
	MaxTopK            int32
	MinMaxOutputTokens int32
	MaxMaxOutputTokens int32
}

var Constants = constants{
	AllowedModels: map[string]bool{
		"gemini-2.0-flash-exp":               true,
		"gemini-1.5-flash":                   true,
		"gemini-1.5-flash-8b":                true,
		"gemini-1.5-pro":                     true,
		"gemini-2.0-flash-thinking-exp-1219": true,
	},
	AllowedLevels: map[string]bool{
		"beginner":     true,
		"intermediate": true,
		"advanced":     true,
	},
	MinTemperature:     0.0,
	MaxTemperature:     1.0,
	MinTopP:            0.0,
	MaxTopP:            1.0,
	MinTopK:            1,
	MaxTopK:            40,
	MinMaxOutputTokens: 1,
	MaxMaxOutputTokens: 2048,
}
