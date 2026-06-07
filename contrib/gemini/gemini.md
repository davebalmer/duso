# Gemini API Module for Duso

Access Google's Gemini API directly from Duso scripts using the OpenAI-compatible interface.

## Setup

Set your API key as an environment variable:

```bash
export GEMINI_API_KEY=your_api_key_here
duso script.du
```

Or pass it explicitly:

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

## Available Models

- `gemini-2.5-pro` (default) - Latest high-capability model
- `gemini-2.5-flash` - Fast and efficient
- `gemini-2.5-flash-lite` - Budget-friendly option
- `gemini-1.5-pro` - Previous generation pro model
- `gemini-1.5-flash` - Previous generation flash model

See Google's [Gemini API documentation](https://ai.google.dev/gemini-api) for the latest model list.

## Configuration Options

Same as OpenAI module - see [openai.md](/contrib/openai/openai.md) for full reference.

Key differences:
- API key environment variable: `GEMINI_API_KEY` (not `OPENAI_API_KEY`)
- Default model: `gemini-2.5-pro`
- Endpoint: `https://generativelanguage.googleapis.com/v1beta/openai/chat/completions`

## Environment Variables

- `GEMINI_API_KEY` - Your API key (required if not passed in config)

## See Also

- [openai.md](/contrib/openai/openai.md) - Full API documentation (identical interface)
- [Google AI Studio](https://aistudio.google.com) - Get your API key free
- [Gemini API Documentation](https://ai.google.dev/gemini-api) - Complete API reference
