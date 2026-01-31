package google

import (
	"fmt"
	"strings"

	"cloud.google.com/go/vertexai/genai"
	"cloud.google.com/go/vertexai/genai/tokenizer"
)

// Provider implements tokenization for Google Gemini models
type Provider struct{}

// New creates a new Google provider
func New() *Provider {
	return &Provider{}
}

// Model to Vertex AI model mapping
// Only models supported by the local tokenizer
var models = map[string]string{
	"gemini-1.5":       "gemini-1.5-flash",
	"gemini-1.5-flash": "gemini-1.5-flash",
	"gemini-1.5-pro":   "gemini-1.5-pro",
}

func (p *Provider) Name() string {
	return "Google"
}

func (p *Provider) CountTokens(text, model string) (int, error) {
	vertexModel, exists := models[strings.ToLower(model)]
	if !exists {
		return 0, fmt.Errorf("unsupported Google model: %s", model)
	}

	tok, err := tokenizer.New(vertexModel)
	if err != nil {
		return 0, fmt.Errorf("failed to create tokenizer for %s: %w", vertexModel, err)
	}

	resp, err := tok.CountTokens(genai.Text(text))
	if err != nil {
		return 0, fmt.Errorf("failed to count tokens: %w", err)
	}

	return int(resp.TotalTokens), nil
}

func (p *Provider) SupportsModel(model string) bool {
	_, exists := models[strings.ToLower(model)]
	return exists
}

func (p *Provider) IsExact() bool {
	return true
}

func (p *Provider) Models() []string {
	result := make([]string, 0, len(models))
	for name := range models {
		result = append(result, name)
	}
	return result
}
