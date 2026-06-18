# websocket()

Establish a client WebSocket connection to a server.

`websocket(url [, config])`

## Parameters

- `url` (string) - WebSocket URL (ws://, wss://, http://, or https://)
- `config` (optional, object) - Configuration object:
  - `headers` (object) - Custom headers to send in the upgrade request (e.g., authentication tokens)

## Returns

WebSocket connection object with methods and properties

## Connection Object

The returned object has the following methods and properties:

- `id` (string) - Unique identifier for this connection (UUID v4)
- `read([timeout])` - Block until a message is received
  - `timeout` (optional, number) - Wait timeout in seconds. If omitted, blocks indefinitely.
  - Returns the message string, or `nil` on disconnect or timeout.
  - Supports both positional: `read(5)` and named: `read(timeout=5)` arguments.
- `write(message)` - Send a message to the server
- `close()` - Close the WebSocket connection
- `is_connected()` - Check if connection is still open (returns boolean)

## Examples

Basic WebSocket connection:

```duso
ws = websocket("ws://localhost:8080/chat")
print("Connected with ID: " + ws.id)

// Send a message
ws.write("Hello, server!")

// Receive message with timeout
msg = ws.read(timeout=10)
if msg then
  print("Received: " + msg)
end

// Close connection
ws.close()
```

WebSocket connection with custom headers:

```duso
ws = websocket("wss://api.example.com/stream", {
  headers = {
    "Authorization" = "Bearer token_123",
    "X-Client-ID" = "my-client"
  }
})

// Message loop
while ws.is_connected() do
  msg = ws.read(timeout=30)
  if msg == nil then
    print("Timeout or disconnected")
    break
  end
  print("Got: " + msg)
end

ws.close()
```

Automatic URL scheme conversion:

```duso
// http automatically converts to ws
ws = websocket("http://localhost:8080/ws")

// https automatically converts to wss
ws = websocket("https://api.example.com/ws")
```

Discord bot using WebSocket client:

```duso
ws = websocket("wss://gateway.discord.gg/", {
  headers = {
    "User-Agent" = "DiscordBot (my-bot/1.0)"
  }
})

ws.write(format_json({
  op = 0,
  d = {
    token = "bot-token",
    properties = {
      os = "linux",
      browser = "duso",
      device = "duso"
    }
  }
}))

// Listen for messages
while ws.is_connected() do
  msg = ws.read(timeout=60)
  if msg then
    event = parse_json(msg)
    print("Got event: " + tostring(event.op))
  end
end
```

## Error Handling

WebSocket operations throw errors on:
- Invalid URL format
- Connection failure (network error, timeout)
- Invalid scheme (must be ws, wss, http, or https)

Use try/catch to handle connection errors:

```duso
try
  ws = websocket("ws://unreachable.example.com")
catch (e)
  print("Failed to connect: " + e)
  exit(1)
end
```

## Message Format

Messages are always text (UTF-8 strings). Binary frames are not supported in this version.

## Timeout Behavior

- `read()` with no timeout blocks indefinitely until a message arrives
- `read(timeout)` returns `nil` if timeout expires before a message arrives
- `read()` returns `nil` if the server closes the connection
- Timeout is specified in seconds as a floating-point number

## Connection Lifecycle

1. **Connect** — `websocket()` dials the URL and returns immediately on success
2. **Read/Write** — Use `read()` and `write()` in a message loop
3. **Disconnect** — Server closes, client calls `close()`, or network error
4. **Cleanup** — After `close()`, further `read()`/`write()` calls fail gracefully

## Performance Notes

- `read()` is blocking—it blocks the script until a message arrives or timeout expires
- For long-running connections, use `read()` in a loop rather than polling
- Each script has its own connection; connections are not shared across spawned scripts
- Use `datastore()` to coordinate multiple WebSocket connections

## See Also

- [http_server() - WebSocket server](/docs/reference/http_server.md) - Create WebSocket servers
- [fetch() - HTTP client requests](/docs/reference/fetch.md) - Make regular HTTP requests
- [spawn() - Run scripts concurrently](/docs/reference/spawn.md) - Run message handlers in parallel
- [datastore() - Shared state](/docs/reference/datastore.md) - Coordinate between connections
