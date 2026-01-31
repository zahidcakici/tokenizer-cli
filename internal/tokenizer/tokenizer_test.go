package tokenizer

import (
	"strings"
	"testing"
)

func TestCountTokensOpenAI(t *testing.T) {
	result, err := CountTokens("Hello world", "gpt-4.1")
	if err != nil {
		t.Fatalf("CountTokens() error = %v", err)
	}

	if result.Provider != "OpenAI" {
		t.Errorf("Expected provider OpenAI, got %s", result.Provider)
	}

	if result.IsEstimate {
		t.Error("OpenAI should be exact, not estimate")
	}

	if result.TokenCount < 2 || result.TokenCount > 3 {
		t.Errorf("Expected 2-3 tokens, got %d", result.TokenCount)
	}
}

func TestCountTokensOpenAIMultipleModels(t *testing.T) {
	models := []string{"gpt-4", "gpt-3.5-turbo", "gpt-4.1"}

	for _, model := range models {
		result, err := CountTokens("Hello world", model)
		if err != nil {
			t.Fatalf("CountTokens() error = %v for model %s", err, model)
		}

		if result.Provider != "OpenAI" {
			t.Errorf("Expected provider OpenAI, got %s for model %s", result.Provider, model)
		}

		if result.IsEstimate {
			t.Error("OpenAI should be exact, not estimate")
		}
	}
}

func TestCountTokensGoogle(t *testing.T) {
	result, err := CountTokens("Hello world", "gemini-1.5-flash")
	if err != nil {
		t.Fatalf("CountTokens() error = %v", err)
	}

	if result.Provider != "Google" {
		t.Errorf("Expected provider Google, got %s", result.Provider)
	}

	if result.IsEstimate {
		t.Error("Google Gemini should be exact, not estimate")
	}
}

func TestCountTokensGoogleMultipleModels(t *testing.T) {
	models := []string{"gemini-1.5-flash", "gemini-1.5-pro"}

	for _, model := range models {
		result, err := CountTokens("Hello world", model)
		if err != nil {
			t.Fatalf("CountTokens() error = %v for model %s", err, model)
		}

		if result.Provider != "Google" {
			t.Errorf("Expected provider Google, got %s for model %s", result.Provider, model)
		}

		if result.IsEstimate {
			t.Error("Google should be exact, not estimate")
		}
	}
}

func TestCountTokensInvalidModel(t *testing.T) {
	_, err := CountTokens("test", "invalid-model-xyz")
	if err == nil {
		t.Error("Expected error for invalid model")
	}
}

func TestListModels(t *testing.T) {
	models := ListModels()

	// Check that we have at least OpenAI and Google providers
	expectedProviders := []string{"OpenAI", "Google"}
	for _, provider := range expectedProviders {
		if _, exists := models[provider]; !exists {
			t.Errorf("Missing provider %s", provider)
		}
	}

	// Verify each provider has models
	for provider, modelList := range models {
		if len(modelList) == 0 {
			t.Errorf("Provider %s has no models", provider)
		}
	}
}

func TestGetDefaultModel(t *testing.T) {
	if GetDefaultModel() != "gpt-4.1" {
		t.Errorf("Default model should be gpt-4.1, got %s", GetDefaultModel())
	}
}

func TestCountTokensEmptyString(t *testing.T) {
	result, err := CountTokens("", "gpt-4.1")
	if err != nil {
		t.Fatalf("CountTokens() error = %v", err)
	}

	if result.TokenCount != 0 {
		t.Errorf("Expected 0 tokens for empty string, got %d", result.TokenCount)
	}
}

func TestCountTokensLongText(t *testing.T) {
	// Create a longer text
	longText := strings.Repeat("This is a test sentence. ", 100)

	result, err := CountTokens(longText, "gpt-4.1")
	if err != nil {
		t.Fatalf("CountTokens() error = %v", err)
	}

	if result.TokenCount == 0 {
		t.Error("Expected non-zero token count for long text")
	}

	// Token count should be reasonable (approximately 500+ tokens)
	if result.TokenCount < 400 {
		t.Errorf("Expected at least 400 tokens for long text, got %d", result.TokenCount)
	}
}

func TestCountTokensAllProviders(t *testing.T) {
	tests := []struct {
		name     string
		model    string
		provider string
		exact    bool
	}{
		{"OpenAI GPT-4.1", "gpt-4.1", "OpenAI", true},
		{"OpenAI GPT-4", "gpt-4", "OpenAI", true},
		{"OpenAI GPT-3.5 Turbo", "gpt-3.5-turbo", "OpenAI", true},
		{"Google Gemini 1.5 Flash", "gemini-1.5-flash", "Google", true},
		{"Google Gemini 1.5 Pro", "gemini-1.5-pro", "Google", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := CountTokens("Hello world", tt.model)
			if err != nil {
				t.Fatalf("CountTokens() error = %v", err)
			}

			if result.Provider != tt.provider {
				t.Errorf("Expected provider %s, got %s", tt.provider, result.Provider)
			}

			if result.IsEstimate == tt.exact {
				t.Errorf("Expected IsEstimate=%v, got %v", !tt.exact, result.IsEstimate)
			}

			if result.Model != tt.model {
				t.Errorf("Expected model %s, got %s", tt.model, result.Model)
			}

			if result.TokenCount <= 0 {
				t.Errorf("Expected positive token count, got %d", result.TokenCount)
			}
		})
	}
}

func TestGetProviderInfo(t *testing.T) {
	tests := []struct {
		name        string
		model       string
		expProvider string
		expExact    bool
		expFound    bool
	}{
		{"OpenAI Model", "gpt-4.1", "OpenAI", true, true},
		{"Google Model", "gemini-1.5-flash", "Google", true, true},
		{"Google Pro Model", "gemini-1.5-pro", "Google", true, true},
		{"Invalid Model", "nonexistent-model", "", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider, exact, found := GetProviderInfo(tt.model)

			if found != tt.expFound {
				t.Errorf("Expected found=%v, got %v", tt.expFound, found)
			}

			if found {
				if provider != tt.expProvider {
					t.Errorf("Expected provider %s, got %s", tt.expProvider, provider)
				}
				if exact != tt.expExact {
					t.Errorf("Expected exact=%v, got %v", tt.expExact, exact)
				}
			}
		})
	}
}
