# Grok AI API Module for Duso

Access Grok AI's API directly from Duso scripts with an options-based, idiomatic interface.

Grok is fully OpenAI-compatible, so this module provides the same familiar interface as the OpenAI module.

## Setup

Set your API key as an environment variable:

```bash
export XAI_API_KEY=xai-xxxxx
duso script.du
```

Or pass it explicitly in your script:

```duso
grok = require("grok")
response = grok.prompt("Hello", {key = "xai-xxxxx"})
```

Get your API key at [https://console.x.ai/](https://console.x.ai/).

## Quick Start

### One-shot query

```duso
grok = require("grok")
response = grok.prompt("What is Duso?")
print(response)
```

### Multi-turn conversation

```duso
grok = require("grok")

chat = grok.session({
  system = "You are a helpful assistant"
})

response1 = chat.prompt("What is a closure?")
response2 = chat.prompt("Can you give me an example?")

print(chat.usage)  // Check token usage
```

### With temperature control

```duso
grok = require("grok")

// Lower temperature = more deterministic
response = grok.prompt("Solve this math problem: 2 + 2", {
  temperature = 0.5
})

// Higher temperature = more creative
response = grok.prompt("Write a poem about code", {
  temperature = 1.0
})
```

### With tools (Agent patterns)

```duso
grok = require("grok")

// Define a tool using standard format
var calculator = {
  name = "calculator",
  description = "Performs basic math operations",
  parameters = {
    operation = {type = "string"},
    a = {type = "number"},
    b = {type = "number"}
  },
  required = ["operation", "a", "b"],
  handler = function(args)
    if args.operation == "add" then return args.a + args.b end
    if args.operation == "multiply" then return args.a * args.b end
  end
}

// Create agent - handler is automatically extracted!
agent = grok.session({
  tools = [calculator]
})

// Ask the agent - it will automatically call tools
response = agent.prompt("What is 15 * 27?")
print(response)  // "405"
```

## API Reference

### `grok.prompt(message, [options])`

Send a one-off message to Grok and get a response.

**Parameters:**
- `message` (string, required) — Your message
- `options` (object, optional) — Configuration options

**Options:**
- `model` — Model name (default: `"grok-4-0709"`)
- `key` — API key (uses `XAI_API_KEY` env var if not provided)
- `temperature` — Sampling temperature 0–2 (default: 1.0)
- `max_tokens` — Max response length (default: 2048)
- `top_p` — Nucleus sampling parameter
- `system` — System prompt
- `tools` — Array of tool definitions
- `tool_choice` — Tool selection strategy ("auto", "none", or tool name)
- `timeout` — Request timeout in seconds (default: 30)

**Returns:** Response text as a string

### `grok.session([options])`

Create a multi-turn conversation session.

**Options:** Same as `grok.prompt()`

**Session object methods:**

- `prompt(message)` — Send a message, get response
- `continue_conversation()` — Generate the next response without adding user input
- `add_tool_result(tool_call_id, result)` — Provide tool execution result
- `clear()` — Reset conversation history
- `set(key, value)` — Update configuration
- `usage` — Object with `{input_tokens, output_tokens}`

**Example:**

```duso
grok = require("grok")
chat = grok.session()

r1 = chat.prompt("What's 2+2?")
r2 = chat.prompt("And 3*4?")

print("Input tokens: " + chat.usage.input_tokens)
print("Output tokens: " + chat.usage.output_tokens)
```

### `grok.models([key])`

List available Grok models.

**Parameters:**
- `key` (string, optional) — API key (uses `XAI_API_KEY` env var if not provided)

**Returns:** Array of model objects

**Example:**

```duso
grok = require("grok")
models = grok.models()
for model in models do
  print(model.id)
end
```

## Available Models

Current default: `grok-4-0709`

See [xAI's documentation](https://x.ai/) for the full list of available models.

## See Also

- [OpenAI module](/contrib/openai/openai.md) — Uses the same interface
- [`fetch()`](/docs/reference/fetch.md) — Make custom HTTP requests
