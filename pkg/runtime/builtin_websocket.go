package runtime

import (
	"fmt"
	"time"

	"github.com/duso-org/duso/pkg/script"
)

// builtinSendWebSocket sends a message to one or more WebSocket connections by ID
// Usage: send_websocket(conn_id, message) or send_websocket([conn_ids], message)
// Returns: bytes sent (number) for single connection, array of results for multiple, or nil if not found/queue full
func builtinSendWebSocket(evaluator *Evaluator, args map[string]any) (any, error) {
	// Get connection ID(s) - can be string or array of strings
	var connIDs []string

	if idArg, ok := args["0"]; ok {
		switch v := idArg.(type) {
		case string:
			if v != "" {
				connIDs = append(connIDs, v)
			}
		case []any:
			for _, id := range v {
				idStr := fmt.Sprintf("%v", id)
				if idStr != "" {
					connIDs = append(connIDs, idStr)
				}
			}
		default:
			idStr := fmt.Sprintf("%v", idArg)
			if idStr != "" {
				connIDs = append(connIDs, idStr)
			}
		}
	} else if idArg, ok := args["conn_id"]; ok {
		switch v := idArg.(type) {
		case string:
			if v != "" {
				connIDs = append(connIDs, v)
			}
		case []any:
			for _, id := range v {
				idStr := fmt.Sprintf("%v", id)
				if idStr != "" {
					connIDs = append(connIDs, idStr)
				}
			}
		default:
			idStr := fmt.Sprintf("%v", idArg)
			if idStr != "" {
				connIDs = append(connIDs, idStr)
			}
		}
	} else {
		return nil, fmt.Errorf("send_websocket() requires a connection ID or array of IDs")
	}

	if len(connIDs) == 0 {
		return nil, fmt.Errorf("send_websocket() connection ID(s) cannot be empty")
	}

	// Get message and stringify it
	var message string
	if msg, ok := args["1"]; ok {
		message = fmt.Sprintf("%v", msg)
	} else if msg, ok := args["message"]; ok {
		message = fmt.Sprintf("%v", msg)
	} else {
		return nil, fmt.Errorf("send_websocket() requires a message")
	}

	// Single connection
	if len(connIDs) == 1 {
		conn := GetConnection(connIDs[0])
		if conn == nil {
			return nil, nil // Connection not found, return nil
		}
		return conn.Write(message), nil
	}

	// Multiple connections - return array of results
	results := make([]any, len(connIDs))
	for i, connID := range connIDs {
		conn := GetConnection(connID)
		if conn == nil {
			results[i] = nil
		} else {
			results[i] = conn.Write(message)
		}
	}
	return results, nil
}

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
	wsConfig := DefaultWebSocketConfig()

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
			// Parse WebSocket config options
			if readQSize, ok := optsMap["read_queue_size"].(float64); ok {
				wsConfig.ReadQueueSize = int(readQSize)
			}
			if writeQSize, ok := optsMap["write_queue_size"].(float64); ok {
				wsConfig.WriteQueueSize = int(writeQSize)
			}
			if readTimeout, ok := optsMap["read_timeout"].(float64); ok {
				wsConfig.DefaultReadTimeout = time.Duration(readTimeout) * time.Second
			}
			if idleTimeout, ok := optsMap["idle_timeout"].(float64); ok {
				wsConfig.IdleTimeout = time.Duration(idleTimeout) * time.Second
			}
			if maxMsgSize, ok := optsMap["max_message_size"].(float64); ok {
				wsConfig.MaxMessageSize = int64(maxMsgSize)
			}
			if maxMsgPerSec, ok := optsMap["max_messages_per_second"].(float64); ok {
				wsConfig.MaxMessagesPerSecond = int(maxMsgPerSec)
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
			// Parse WebSocket config options
			if readQSize, ok := optsMap["read_queue_size"].(float64); ok {
				wsConfig.ReadQueueSize = int(readQSize)
			}
			if writeQSize, ok := optsMap["write_queue_size"].(float64); ok {
				wsConfig.WriteQueueSize = int(writeQSize)
			}
			if readTimeout, ok := optsMap["read_timeout"].(float64); ok {
				wsConfig.DefaultReadTimeout = time.Duration(readTimeout) * time.Second
			}
			if idleTimeout, ok := optsMap["idle_timeout"].(float64); ok {
				wsConfig.IdleTimeout = time.Duration(idleTimeout) * time.Second
			}
			if maxMsgSize, ok := optsMap["max_message_size"].(float64); ok {
				wsConfig.MaxMessageSize = int64(maxMsgSize)
			}
			if maxMsgPerSec, ok := optsMap["max_messages_per_second"].(float64); ok {
				wsConfig.MaxMessagesPerSecond = int(maxMsgPerSec)
			}
		}
	} else {
		headers = make(map[string]string)
	}

	// Connect to WebSocket with config
	conn, err := NewWebSocketClientConnectionWithConfig(url, headers, wsConfig)
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
				return nil, nil // Connection closed
			}
			return msg, nil // Return actual message (including empty string)
		}),
		"write": script.NewGoFunction(func(evaluator *Evaluator, args map[string]any) (any, error) {
			msg, ok := args["0"]
			if !ok {
				return nil, fmt.Errorf("write() requires a message argument")
			}
			msgStr := fmt.Sprintf("%v", msg)
			return conn.Write(msgStr), nil // Returns bytes (number) or nil on queue full
		}),
		"close": script.NewGoFunction(func(evaluator *Evaluator, args map[string]any) (any, error) {
			return nil, conn.Close()
		}),
		"is_connected": script.NewGoFunction(func(evaluator *Evaluator, args map[string]any) (any, error) {
			return conn.IsConnected(), nil
		}),
	}, nil
}
