package runtime

import (
	"fmt"
	"time"

	"github.com/duso-org/duso/pkg/script"
)

// builtinWebSocket establishes a WebSocket client connection
// Usage: websocket(url [, config])
// Returns: WebSocket connection object with read(), write(), close(), is_connected(), and id methods
func builtinWebSocket(evaluator *Evaluator, args map[string]any) (any, error) {
	// Get URL from first positional or named argument
	var url string

	if u, ok := args["0"]; ok {
		url = fmt.Sprintf("%v", u)
	} else if u, ok := args["url"]; ok {
		url = fmt.Sprintf("%v", u)
	} else {
		return nil, fmt.Errorf("websocket() requires a URL")
	}

	if url == "" {
		return nil, fmt.Errorf("websocket() URL cannot be empty")
	}

	// Get options from second positional or named argument
	var headers map[string]string
	if opts, ok := args["1"]; ok {
		if optsMap, ok := opts.(map[string]any); ok {
			headers = make(map[string]string)
			if headersOpt, ok := optsMap["headers"]; ok {
				if headerMap, ok := headersOpt.(map[string]any); ok {
					for k, v := range headerMap {
						headers[k] = fmt.Sprintf("%v", v)
					}
				}
			}
		}
	} else if opts, ok := args["config"]; ok {
		if optsMap, ok := opts.(map[string]any); ok {
			headers = make(map[string]string)
			if headersOpt, ok := optsMap["headers"]; ok {
				if headerMap, ok := headersOpt.(map[string]any); ok {
					for k, v := range headerMap {
						headers[k] = fmt.Sprintf("%v", v)
					}
				}
			}
		}
	} else {
		headers = make(map[string]string)
	}

	// Connect to WebSocket
	conn, err := NewWebSocketClientConnection(url, headers)
	if err != nil {
		return nil, err
	}

	// Return object with id and methods
	return map[string]any{
		"id": conn.ID(),
		"read": script.NewGoFunction(func(evaluator *Evaluator, args map[string]any) (any, error) {
			// Get optional timeout (positional or named)
			var timeout *time.Duration
			if t, ok := args["0"]; ok {
				if timeoutSec, ok := t.(float64); ok && timeoutSec > 0 {
					d := time.Duration(timeoutSec * float64(time.Second))
					timeout = &d
				}
			} else if t, ok := args["timeout"]; ok {
				if timeoutSec, ok := t.(float64); ok && timeoutSec > 0 {
					d := time.Duration(timeoutSec * float64(time.Second))
					timeout = &d
				}
			}

			msg, err := conn.Read(timeout)
			if err != nil {
				return nil, err
			}
			if msg == "" {
				return nil, nil // Return nil (not empty string) on disconnect/timeout
			}
			return msg, nil
		}),
		"write": script.NewGoFunction(func(evaluator *Evaluator, args map[string]any) (any, error) {
			msg, ok := args["0"]
			if !ok {
				return nil, fmt.Errorf("write() requires a message argument")
			}
			msgStr := fmt.Sprintf("%v", msg)
			return nil, conn.Write(msgStr)
		}),
		"close": script.NewGoFunction(func(evaluator *Evaluator, args map[string]any) (any, error) {
			return nil, conn.Close()
		}),
		"is_connected": script.NewGoFunction(func(evaluator *Evaluator, args map[string]any) (any, error) {
			return conn.IsConnected(), nil
		}),
	}, nil
}
