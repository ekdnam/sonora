# Leo - AI Course Generator

Leo is a Go-based application that leverages Google's Gemini AI to generate structured course content and learning plans. It's designed to create customized educational content for different learning levels, from beginner to advanced.

## Features

- Generate structured course plans for any STEM topic
- Support for multiple difficulty levels (beginner, intermediate, advanced)
- Topic validation and alternative topic suggestions
- Streaming and non-streaming content generation
- Configurable AI model parameters

## Prerequisites

- Go 1.23.4 or higher
- Google Cloud Gemini API key

## Installation

1. Clone the repository
2. Install dependencies:
```bash
go mod download
```

3. Create a `.env` file in the root directory with your Gemini API key:
```bash
GEMINI_API_KEY=your_api_key_here
```

## Configuration

The application supports various Gemini models and configuration parameters:

### Supported Models
- gemini-2.0-flash-exp
- gemini-1.5-flash
- gemini-1.5-flash-8b
- gemini-1.5-pro
- gemini-2.0-flash-thinking-exp-1219

### Model Configuration Parameters
- Temperature: 0.0 - 1.0
- TopP: 0.0 - 1.0
- TopK: 1 - 40
- MaxOutputTokens: 1 - 2048

## Usage

Basic usage example:

```go
package main

import (
    "context"
    "leo/src/core"
    TypeLeo "leo/src/typeLeo"
)

func main() {
    // Initialize model with configuration
    model, err := core.GetModel(client, TypeLeo.GenerativeModelConfig{
        ModelName:       "gemini-1.5-pro",
        Temperature:     0.5,
        TopP:           0.95,
        TopK:           40,
        MaxOutputTokens: 8192,
    })

    // Generate a course plan
    topic := "your_topic"
    level := "advanced"
    resp, err := core.GeneratePlan(ctx, model, topic, level)
}
```

## Testing

The project includes both unit tests and integration tests. To run integration tests:

```bash
go test -tags=integration ./...
```

## Project Structure

- `/src`
  - `/core` - Core AI integration and content generation
  - `/prompts` - System prompts and templates
  - `/typeLeo` - Type definitions and schemas
  - `/utils` - Utility functions
  - `/psql` - Database integration (WIP)

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.