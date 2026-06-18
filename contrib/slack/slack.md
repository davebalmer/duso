# Slack Module

Slack API integration for duso. Provides webhooks and Socket Mode client.

## Usage

```duso
slack = require("slack")
```

## Functions

### post_webhook(url, payload)

Post a message to a Slack incoming webhook.

**Parameters:**
- `url` (string) - Webhook URL from Slack
- `payload` (object) - Message payload (text, blocks, attachments, etc.)

**Returns:** Boolean (true if successful)

**Example:**
```duso
slack.post_webhook("https://hooks.slack.com/services/...", {
  text = "Hello from Duso!"
})
```

### session(config)

Create a Socket Mode client connection.

**Parameters:**
- `config` (object) - Configuration:
  - `app_token` (string, required) - Slack app token (starts with `xapp_`)

**Returns:** Session object

**Example:**
```duso
bot = slack.session({
  app_token = "xapp_1234567890_..."
})
```

## Session Object

Returned by `session()`. Handles Socket Mode connection.

### Methods

- `read([timeout])` - Block until an event arrives. Returns event object or nil on timeout.
- `write(payload)` - Send a reply or response.
- `is_connected()` - Check if connection is still open.
- `close()` - Close the connection.

### Properties

- `id` (string) - Connection ID (UUID)

## Example: Echo Bot

```duso
slack = require("slack")

bot = slack.session({
  app_token = "xapp_..."
})

print("Bot connected: " + bot.id)

while bot.is_connected() do
  event = bot.read(timeout=30)
  
  if event == nil then
    continue
  end

  if event.type == "events_api" then
    payload = event.payload
    
    if payload.event and payload.event.type == "message" then
      msg = payload.event.text
      channel = payload.event.channel
      
      print("Got message: " + msg)
      
      // Send reply via API
      fetch("https://slack.com/api/chat.postMessage", {
        method = "POST",
        headers = {
          "Authorization" = "Bearer xoxb_...",
          "Content-Type" = "application/json"
        },
        body = format_json({
          channel = channel,
          text = "Echo: " + msg
        })
      })
    end
  end
end

bot.close()
```

## Example: Webhook Alert

```duso
slack = require("slack")

webhook_url = "https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXX"

slack.post_webhook(webhook_url, {
  text = "Alert: Database backup complete",
  blocks = [
    {
      type = "section",
      text = {
        type = "mrkdwn",
        text = "*Backup Status*\nCompleted at " + format_time(now())
      }
    }
  ]
})
```

## Notes

- Socket Mode requires a Slack app with appropriate permissions
- Events are auto-acknowledged by the client
- For blocking operations, use webhooks; for interactive bots, use Socket Mode
- Slack API calls still require bearer token (passed separately to fetch)

## See Also

- [websocket() - WebSocket client](/docs/reference/websocket.md)
- [fetch() - HTTP requests](/docs/reference/fetch.md)
