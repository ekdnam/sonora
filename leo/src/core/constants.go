package core

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

func (c *constants) IsAllowedModel(modelName string) bool {
	return c.AllowedModels[modelName]
}

func (c *constants) IsAllowedLevel(level string) bool {
	return c.AllowedLevels[level]
}

func (c *constants) GetAllModels() []string {
	models := make([]string, 0, len(c.AllowedModels))
	for model := range c.AllowedModels {
		models = append(models, model)
	}
	return models
}

func (c *constants) GetAllLevels() []string {
	levels := make([]string, 0, len(c.AllowedLevels))
	for level := range c.AllowedLevels {
		levels = append(levels, level)
	}
	return levels
}
