package providers

// Provider defines the interface for LLM tokenization providers
type Provider interface {
	// Name returns the provider name (e.g., "OpenAI", "Anthropic")
	Name() string

	// CountTokens counts tokens for the given text and model
	CountTokens(text, model string) (int, error)

	// SupportsModel returns true if this provider supports the given model
	SupportsModel(model string) bool

	// IsExact returns true if token counts are exact (not estimated)
	IsExact() bool

	// Models returns a list of supported model names
	Models() []string
}
