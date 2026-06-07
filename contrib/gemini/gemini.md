# Gemini API Module for Duso

Access Google's Gemini API directly from Duso scripts using the OpenAI-compatible interface.

## Setup

1. Get a free API key from [Google AI Studio](https://aistudio.google.com)
2. Set your API key as an environment variable:

```bash
export GEMINI_API_KEY=your_api_key_here
duso script.du
```

Or pass it explicitly in code:

```duso
gemini = require("gemini")
response = gemini.prompt("Hello", {key = "your_api_key_here"})
```

## Quick Start

```duso
gemini = require("gemini")

// One-shot query
response = gemini.prompt("What is Gemini?")
print(response)

// Multi-turn conversation
chat = gemini.session({
  system = "You are a helpful assistant",
  model = "gemini-2.5-pro"
})

response1 = chat.prompt("What is Google Gemini known for?")
response2 = chat.prompt("What are its advantages?")
print(chat.usage)
```

## Endpoint

Default: `https://generativelanguage.googleapis.com/v1beta/openai/`

Cloud-hosted API (no local setup required, unlike Ollama).

## Available Models

- `gemini-2.5-pro` (default) - Latest high-capability model with 2M context
- `gemini-2.5-flash` - Fast and efficient
- `gemini-2.5-flash-lite` - Budget-friendly option
- `gemini-1.5-pro` - Previous generation pro model
- `gemini-1.5-flash` - Previous generation flash model

See [Gemini API documentation](https://ai.google.dev/gemini-api) for the latest model list.

## API Key Required

Gemini is a cloud API. You need a valid Google API key from [AI Studio](https://aistudio.google.com).

## Configuration Options

Same as OpenAI module - see [openai.md](/contrib/openai/openai.md) for full reference.

Key differences:
- API key environment variable: `GEMINI_API_KEY` (not `OPENAI_API_KEY`)
- Default model: `gemini-2.5-pro`
- Cloud endpoint (no local installation needed)

## See Also

- [openai.md](/contrib/openai/openai.md) - Full API documentation (identical interface)
- [Google AI Studio](https://aistudio.google.com) - Get your free API key
- [Gemini API Documentation](https://ai.google.dev/gemini-api) - Complete API reference
