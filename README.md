[![Apache 2.0 License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0) ![Go 1.25](https://img.shields.io/badge/Go-1.25-darkcyan?logo=go) ![GitHub Release](https://img.shields.io/github/v/release/duso-org/duso)

![Duso logo which is a stylized ASL hand sign for the letter "D"](/docs/duso-logo.png)

# Duso

**A complete server runtime. One 10MB binary with everything built in. No npm, no dependencies, no bloat. A simple scripting language powered by Go, with native goroutine concurrency and instant hot reload.**

## Download & Install

**Pre-built binaries:** [duso.rocks/download](https://duso.rocks/download)

Or install with Homebrew:

```bash
brew install duso-org/tap/duso
```

Or build from source:

**Linux & macOS:**
```bash
git clone https://github.com/duso-org/duso.git
cd duso
./build.sh
```

**Windows PowerShell:**
```powershell
git clone https://github.com/duso-org/duso.git
cd duso
.\build.ps1
```

Then optionally symlink it (Linux & macOS):
```bash
ln -s $(pwd)/bin/duso /usr/local/bin/duso
```

## Quick Start

### Run a script

```bash
duso script.du
```

### Run inline code

```bash
duso -c 'print("Hello, World!")'
```

### Interactive REPL mode

```bash
duso -repl
```

### One-liner web server

```bash
duso -c 'http_server().start()'
```

Open `http://localhost:8080`. You have a working server.

### Basic AI prompt

```duso
ai = require("claude")
print(ai.prompt("What is 2+2?"))
```

### Interactive chatbot

```duso
ai = require("openai")
chat = ai.session()

while true do
  prompt = input("\n\nYou: ")
  if lower(prompt) == "exit" then break end

  write("\n\nOpenAI: ")
  busy("thinking...")
  write(chat.prompt(prompt))
end
```

### AI workflow with parallel experts

```duso
ai = require("claude")

prompt = input("Ask the panel: ")
busy("asking...")

experts = ["Astronomer", "Astrologer", "Biologist", "Accountant"]
responses = parallel(map(experts, function(expert)
  return function()
    return ai.prompt(prompt, {
      system = """
        You are an expert {{expert}}. Always reason and
        interact from this mindset. Limit your field of
        knowledge to this expertise.
      """,
      max_tokens = 500
    })
  end
end))

busy("summarizing...")
summary = ai.prompt("""
  Summarize these responses:
  {{join(responses, "\n\n---\n\n")}}
  
  List 3 things they have in common.
  Then list the 3 things that are the most different.
""")

print(markdown_ansi(summary))
```

### REST API with database

```duso
// server.du
server = http_server({port = 3000})

server.route("GET", "/users/:id", "get-user.du")
server.route("POST", "/users", "create-user.du")

print("Running on :3000")
server.start()
```

```duso
// get-user.du
ctx = context()
user = datastore("users").get(ctx.request().params.id)
ctx.response().json(user)
```

```duso
// create-user.du
ctx = context()
user = ctx.request().json()
datastore("users").set(user.id, user)
ctx.response().json(user)
```

### Concurrent workers with shared state

```duso
// bees.du - Spawn workers, wait for all to finish
bees = 10
swarm = datastore("swarm")
swarm.set("done", 0)
swarm.set("buzzes", 0)

for i = 1, bees do
  spawn("worker.du", {bee_id = i})
end

swarm.wait("done", bees)
print("All done! Total buzzes: " + swarm.get("buzzes"))
```

```duso
// worker.du - Spawned worker increments shared counters
ctx = context()
swarm = datastore("swarm")

buzzes = ceil(random() * 10)
for i = 1, buzzes do
  sleep(random() * 0.5)
  swarm.increment("buzzes")
end

swarm.increment("done")
```

## Learn More

```bash
duso -read              # Interactive guided tour
duso -doc claude        # Look up any function
duso -repl              # Test code in real time
duso -init myproject    # Create a starter project
duso -extract examples  # Extract all examples
```

## Why I Made Duso

**Duso is intentionally simple and predictable.** No magic. No multiple ways to do the same thing. Every pattern is consistent so AI can reason about code reliably and write better scripts faster.

**Duso is a joy to use.** Everything including the runtime, libs, and docs is bundled in a single 10MB binary. No package management. No version conflicts. No stack building. Duso makes coding fun again.

[Dave Balmer](https://balmer.dev), creator of Duso

## Key Features

- **Hot Reload**: Edit, test, deploy. Same binary.
- **One Binary**: Everything included. No npm, pip, cargo.
- **Simple Concurrency**: Goroutines without the complexity. `spawn()` and `parallel()`.
- **Full Web Stack**: HTTP, routing, WebSockets, SSL, CORS, JWT, templates built in.
- **General Purpose**: Build APIs, web servers, scripts, tools, background jobs, batch processors.
- **Built-in Datastore**: ACID NoSQL key-value store, in-memory or persisted, perfect for caching and state.
- **SQL Support**: MySQL, MariaDB, TiDB, CouchDB drivers built in.
- **AI Integrations**: Claude, OpenAI, Gemini, Groq, Ollama, Azure AI. All built in.
- **Docs in Binary**: No internet needed. `duso -doc` for any function.
- **Integrated Debugger**: Breakpoints, stack traces, concurrent-aware.
- **Linter & LSP Server**: Built-in static analysis and language server for IDE integration.
- **Editor Extensions**: Syntax highlighting and code completion for VS Code, JetBrains, Vim.
- **AI-Friendly Design**: Simple, readable syntax that LLM agents understand and extend naturally.
- **Single-Binary Deployment**: Linux, macOS, Windows. Same build, every platform.
- **Open Source**: Apache 2.0. Community-driven.

## Full Documentation

- **Website:** [duso.rocks](https://duso.rocks)
- **Learning Guide:** [docs/learning-duso.md](/docs/learning-duso.md)
- **Built-in:** `duso -read` or `duso -doc <topic>`

## Community Libraries

Built-in integrations for:

- **AI:** Claude, OpenAI, Gemini, Groq, Ollama, Azure AI, DeepSeek
- **Databases:** MySQL, MariaDB, TiDB, CouchDB
- **Payments:** Stripe
- **Testing & Utils:** Icons (Phospher), Zero Language Model (zlm)

See [contrib/](/contrib/) for full list and docs.

## Contributing

Duso is open source under Apache 2.0. Contributions welcome:

1. Fork [github.com/duso-org/duso](https://github.com/duso-org/duso)
2. Create a branch (`git checkout -b feature/thing`)
3. Commit changes (`git commit -am 'add thing'`)
4. Push to branch (`git push origin feature/thing`)
5. Open a Pull Request

See [CONTRIBUTING.md](/CONTRIBUTING.md) and [COMMUNITY.md](/COMMUNITY.md) for guidelines.

## Contributors

- [Dave Balmer](https://balmer.dev): design, development, documentation, dedication

## Sponsors

- **[Shannan.dev](https://shannan.dev)**: Provides AI-driven business intelligence solutions
- **[Ludonode](https://ludonode.com)**: Provides agentic development and consulting

## FAQ

**Q: Is Duso production-ready?**
A: Yes.

**Q: How do I deploy Duso?**
A: It's one binary. `scp` it to a server and run it. Alternatively, containerize it in Docker or deploy to Fly.io, Railway, or Heroku.

**Q: Can I bundle scripts into the binary?**
A: Yes. Use `duso -bundle` to embed scripts, configs, and static files into a single executable.

**Q: Can I extend Duso with Go?**
A: Yes. The language is designed to be extended. You can write custom builtins in Go.

## License

Apache 2.0. See [LICENSE](/LICENSE).
