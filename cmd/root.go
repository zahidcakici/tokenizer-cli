package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/zahidcakici/tokenizer-cli/internal/tokenizer"
)

var (
	filePath  string
	modelName string
	listFlag  bool
)

var rootCmd = &cobra.Command{
	Use:   "tokenizer [text]",
	Short: "Count tokens for LLM models",
	Long: `A CLI tool to count tokens for various Large Language Model (LLM) providers.

Supports OpenAI, Meta Llama, Google Gemini, and Anthropic Claude models.
Token counts are exact for OpenAI, Llama 3+, and Gemini. Claude uses estimation.

Examples:
  tokenizer "Hello world, this is a test"
  tokenizer -f context.txt
  tokenizer -m gpt-4.1 -f prompt.md
  tokenizer --list`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// If the first positional argument is a help request, show help
		if len(args) > 0 {
			if args[0] == "help" || args[0] == "-h" || args[0] == "--help" {
				_ = cmd.Help()
				return nil
			}
		}

		// Handle --list flag
		if listFlag {
			printModelList()
			return nil
		}

		// Get input text
		var text string
		if filePath != "" {
			content, err := os.ReadFile(filePath)
			if err != nil {
				return fmt.Errorf("failed to read file: %w", err)
			}
			text = string(content)
		} else if len(args) > 0 {
			text = strings.Join(args, " ")
		} else {
			return fmt.Errorf("no input provided. Use text argument or -f flag")
		}

		// Count tokens
		result, err := tokenizer.CountTokens(text, modelName)
		if err != nil {
			return err
		}

		// Print result
		printResult(result)
		return nil
	},
}

func printModelList() {
	// Color definitions
	headerColor := color.New(color.FgCyan, color.Bold)
	providerColor := color.New(color.FgYellow, color.Bold)
	defaultColor := color.New(color.FgGreen)
	estimateColor := color.New(color.FgRed, color.Faint)

	// Print header
	fmt.Println()
	headerColor.Println("╔══════════════════════════════════════════════════════╗")
	headerColor.Println("║          📋 SUPPORTED LLM MODELS                    ║")
	headerColor.Println("╚══════════════════════════════════════════════════════╝")
	fmt.Println()

	models := tokenizer.ListModels()

	// Sort providers for consistent output
	providers := make([]string, 0, len(models))
	for provider := range models {
		providers = append(providers, provider)
	}
	sort.Strings(providers)

	for _, provider := range providers {
		// Check if provider is exact
		_, isExact, _ := tokenizer.GetProviderInfo(models[provider][0])
		exactIcon := "✓"
		exactLabel := "exact"
		if !isExact {
			exactIcon = "≈"
			exactLabel = "estimated"
		}

		// Print provider header
		providerColor.Printf("  %s %s ", exactIcon, provider)
		estimateColor.Printf("(%s)\n", exactLabel)

		for _, model := range models[provider] {
			if model == tokenizer.DefaultModel {
				defaultColor.Printf("    ▸ %s ", model)
				color.New(color.FgGreen, color.Bold).Println("★ default")
			} else {
				fmt.Printf("    ▸ %s\n", model)
			}
		}
		fmt.Println()
	}

	// Print legend
	color.New(color.Faint).Println("  Legend: ★ default model  |  ✓ exact count  |  ≈ estimated count")
	fmt.Println()
}

func printResult(result *tokenizer.Result) {
	// Color definitions
	headerColor := color.New(color.FgCyan, color.Bold)
	modelColor := color.New(color.FgMagenta, color.Bold)
	providerColor := color.New(color.FgYellow)
	tokenColor := color.New(color.FgGreen, color.Bold)
	estimateColor := color.New(color.FgRed)

	// Determine accuracy indicator
	accuracyIcon := "✓"
	accuracyLabel := "exact"
	accuracyColorizer := color.New(color.FgGreen)
	if result.IsEstimate {
		accuracyIcon = "≈"
		accuracyLabel = "estimated"
		accuracyColorizer = estimateColor
	}

	// Print result box
	fmt.Println()
	headerColor.Println("╔══════════════════════════════════════════════════════╗")
	headerColor.Println("║          🔢 TOKEN COUNT RESULT                      ║")
	headerColor.Println("╚══════════════════════════════════════════════════════╝")
	fmt.Println()

	// Model information
	fmt.Print("  Model:    ")
	modelColor.Println(result.Model)

	// Provider information
	fmt.Print("  Provider: ")
	providerColor.Println(result.Provider)

	// Token count
	fmt.Print("  Tokens:   ")
	tokenColor.Printf("%d", result.TokenCount)
	fmt.Print(" ")
	accuracyColorizer.Printf("(%s %s)\n", accuracyIcon, accuracyLabel)

	fmt.Println()
}

func init() {
	rootCmd.Flags().StringVarP(&filePath, "file", "f", "", "Path to file to analyze")
	rootCmd.Flags().StringVarP(&modelName, "model", "m", tokenizer.DefaultModel, "Model to use for tokenization")
	rootCmd.Flags().BoolVarP(&listFlag, "list", "l", false, "List all supported models")
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
