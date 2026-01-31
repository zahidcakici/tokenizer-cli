package tokenizer

import (
	"fmt"
	"sort"

	"github.com/zahidcakici/tokenizer-cli/internal/providers"
	"github.com/zahidcakici/tokenizer-cli/internal/providers/google"
	"github.com/zahidcakici/tokenizer-cli/internal/providers/openai"
)

// DefaultModel is the default model when none is specified
const DefaultModel = "gpt-4.1"

// Result holds the token counting result
type Result struct {
	Model      string
	Provider   string
	TokenCount int
	IsEstimate bool
}

// registry holds all registered providers
var registry []providers.Provider

func init() {
	// Register all providers
	registry = []providers.Provider{
		openai.New(),
		google.New(),
	}
}

// CountTokens counts tokens for the given text and model
func CountTokens(text, modelName string) (*Result, error) {
	for _, provider := range registry {
		if provider.SupportsModel(modelName) {
			count, err := provider.CountTokens(text, modelName)
			if err != nil {
				return nil, err
			}

			return &Result{
				Model:      modelName,
				Provider:   provider.Name(),
				TokenCount: count,
				IsEstimate: !provider.IsExact(),
			}, nil
		}
	}

	return nil, fmt.Errorf("unsupported model: %s. Use --list to see supported models", modelName)
}

// ListModels returns all supported models grouped by provider
func ListModels() map[string][]string {
	result := make(map[string][]string)

	for _, provider := range registry {
		models := provider.Models()
		sort.Strings(models)
		result[provider.Name()] = models
	}

	return result
}

// GetDefaultModel returns the default model name
func GetDefaultModel() string {
	return DefaultModel
}

// GetProviderInfo returns provider info for a model
func GetProviderInfo(modelName string) (providerName string, isExact bool, found bool) {
	for _, provider := range registry {
		if provider.SupportsModel(modelName) {
			return provider.Name(), provider.IsExact(), true
		}
	}
	return "", false, false
}
