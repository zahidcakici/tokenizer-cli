# Tokenizer CLI

A fast CLI tool to count tokens for various Large Language Model (LLM) providers.

## Features

- 🚀 **Multi-model support**: OpenAI, Google Gemini
- 📊 **Token counts**: OpenAI, Gemini 1.5
- 📄 **File input**: Analyze tokens from files directly
- ⚡ **Fast**: Local tokenization, no API calls required

## Installation

### Homebrew (macOS/Linux)

```bash
brew install zahidcakici/tap/tokenizer
```

### From Source

```bash
go install github.com/zahidcakici/tokenizer-cli@latest
```

### Alternative: Bash script (curl)

```bash
curl -sSL https://raw.githubusercontent.com/zahidcakici/tokenizer-cli/main/install.sh | sudo bash
```


## Usage

```bash
# Direct text input
tokenizer "Hello world, this is a test"
# Output: Model: gpt-4.1 (OpenAI) | Tokens: 7

# From file
tokenizer -f context.txt

# Specify model
tokenizer -m gpt-4 "Hello world"
# Output: Model: llama-3.3 (Meta) | Tokens: 2

tokenizer -m gemini-1.5-flash "Hello world"
# Output: Model: gemini-1.5-flash (Google) | Tokens: 2


# List all supported models
tokenizer --list
```

## Supported Models

| Provider | Models | Token Counting |
|----------|--------|----------------|
| OpenAI | gpt-4.1, gpt-4o, o1, o3, etc. | ✅ Exact |
| Google | gemini-1.5, gemini-1.5-flash, gemini-1.5-pro | ✅ Exact |

Default model: `gpt-4.1`

## License

MIT
