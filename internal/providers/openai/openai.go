package openai

import (
	"fmt"
	"strings"

	"github.com/pkoukk/tiktoken-go"
)

// Provider implements tokenization for OpenAI models
type Provider struct{}

// New creates a new OpenAI provider
func New() *Provider {
	return &Provider{}
}

// Model configuration
type modelConfig struct {
	encoding string
}

var models = map[string]modelConfig{
	// GPT-4o family (o200k_base encoding)
	"gpt-4o":       {encoding: "o200k_base"},
	"gpt-4o-mini":  {encoding: "o200k_base"},
	"gpt-4.1":      {encoding: "o200k_base"},
	"gpt-4.1-mini": {encoding: "o200k_base"},
	"gpt-4.1-nano": {encoding: "o200k_base"},

	// Reasoning models (o200k_base encoding)
	"o1":      {encoding: "o200k_base"},
	"o1-mini": {encoding: "o200k_base"},
	"o1-pro":  {encoding: "o200k_base"},
	"o3":      {encoding: "o200k_base"},
	"o3-mini": {encoding: "o200k_base"},
	"o4-mini": {encoding: "o200k_base"},

	// GPT-4 and GPT-3.5 (cl100k_base encoding)
	"gpt-4":         {encoding: "cl100k_base"},
	"gpt-4-turbo":   {encoding: "cl100k_base"},
	"gpt-3.5-turbo": {encoding: "cl100k_base"},
}

func (p *Provider) Name() string {
	return "OpenAI"
}

func (p *Provider) CountTokens(text, model string) (int, error) {
	cfg, exists := models[strings.ToLower(model)]
	if !exists {
		return 0, fmt.Errorf("unsupported OpenAI model: %s", model)
	}

	tke, err := tiktoken.GetEncoding(cfg.encoding)
	if err != nil {
		return 0, fmt.Errorf("failed to get encoding %s: %w", cfg.encoding, err)
	}

	tokens := tke.Encode(text, nil, nil)
	return len(tokens), nil
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
