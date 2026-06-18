package runtime

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/websocket"
)

// WebSocketConnection represents an active WebSocket connection in Duso
type WebSocketConnection struct {
	ws     *websocket.Conn
	closed bool
	mutex  sync.Mutex
	id     string
}

// NewWebSocketConnection creates a new WebSocket connection wrapper (server-side)
func NewWebSocketConnection(ws *websocket.Conn) *WebSocketConnection {
	return &WebSocketConnection{
		ws:     ws,
		closed: false,
		id:     generateUUIDv4(),
	}
}

// NewWebSocketClientConnection creates a client WebSocket connection
func NewWebSocketClientConnection(urlStr string, headers map[string]string) (*WebSocketConnection, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("invalid WebSocket URL: %w", err)
	}

	// Convert http/https to ws/wss
	switch u.Scheme {
	case "http":
		u.Scheme = "ws"
	case "https":
		u.Scheme = "wss"
	case "ws", "wss":
		// Already correct
	default:
		return nil, fmt.Errorf("unsupported scheme: %s (use http, https, ws, or wss)", u.Scheme)
	}

	// Dial the WebSocket
	config, err := websocket.NewConfig(u.String(), u.String())
	if err != nil {
		return nil, fmt.Errorf("invalid WebSocket config: %w", err)
	}

	// Add custom headers
	for k, v := range headers {
		config.Header.Set(k, v)
	}

	ws, err := websocket.DialConfig(config)
	if err != nil {
		return nil, fmt.Errorf("WebSocket connection failed: %w", err)
	}

	return &WebSocketConnection{
		ws:     ws,
		closed: false,
		id:     generateUUIDv4(),
	}, nil
}

// Accept accepts the WebSocket connection (protocol handshake already done by upgrade)
func (wsc *WebSocketConnection) Accept() error {
	wsc.mutex.Lock()
	defer wsc.mutex.Unlock()

	if wsc.closed {
		return fmt.Errorf("connection already closed")
	}

	// Connection is already accepted by the HTTP upgrade
	return nil
}

// Read blocks until a message is received or connection closes
// Returns the message string, or empty string on disconnect/timeout
func (wsc *WebSocketConnection) Read(timeout *time.Duration) (string, error) {
	wsc.mutex.Lock()
	if wsc.closed {
		wsc.mutex.Unlock()
		return "", fmt.Errorf("connection closed")
	}
	wsc.mutex.Unlock()

	if timeout != nil {
		wsc.ws.SetReadDeadline(time.Now().Add(*timeout))
	}

	var msg string
	err := websocket.Message.Receive(wsc.ws, &msg)

	if err != nil {
		wsc.mutex.Lock()
		wsc.closed = true
		wsc.mutex.Unlock()

		// Return error on EOF/disconnect
		if strings.Contains(err.Error(), "EOF") || strings.Contains(err.Error(), "closed") {
			return "", nil // Nil equivalent for disconnect
		}
		// Check for timeout
		if strings.Contains(err.Error(), "deadline exceeded") || strings.Contains(err.Error(), "timeout") {
			return "", nil // Nil equivalent for timeout
		}
		return "", fmt.Errorf("websocket receive error: %w", err)
	}

	// Clear read deadline after successful receive
	if timeout != nil {
		wsc.ws.SetReadDeadline(time.Time{})
	}

	return msg, nil
}

// Write sends a message to the WebSocket client
func (wsc *WebSocketConnection) Write(message string) error {
	wsc.mutex.Lock()
	defer wsc.mutex.Unlock()

	if wsc.closed {
		return fmt.Errorf("connection closed")
	}

	return websocket.Message.Send(wsc.ws, message)
}

// Close closes the WebSocket connection
func (wsc *WebSocketConnection) Close() error {
	wsc.mutex.Lock()
	defer wsc.mutex.Unlock()

	if wsc.closed {
		return nil
	}

	wsc.closed = true
	return wsc.ws.Close()
}

// IsConnected returns whether the connection is still open
func (wsc *WebSocketConnection) IsConnected() bool {
	wsc.mutex.Lock()
	defer wsc.mutex.Unlock()
	return !wsc.closed
}

// ID returns the unique identifier for this connection
func (wsc *WebSocketConnection) ID() string {
	return wsc.id
}

// generateUUIDv4 generates a UUID v4 (random) string
func generateUUIDv4() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("ws_%d", time.Now().UnixNano())
	}

	// Set version 4 (random) and variant bits
	b[6] = (b[6] & 0x0f) | 0x40 // version 4
	b[8] = (b[8] & 0x3f) | 0x80 // variant

	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

// IsWebSocketUpgrade checks if the request is a WebSocket upgrade request
func IsWebSocketUpgrade(r *http.Request) bool {
	return strings.ToLower(r.Header.Get("Upgrade")) == "websocket" &&
		strings.Contains(strings.ToLower(r.Header.Get("Connection")), "upgrade")
}

// WebSocketHandler creates an http.Handler that performs WebSocket upgrade
// and calls the provided upgrade handler
func WebSocketHandler(upgradeHandler func(*WebSocketConnection, *http.Request) error) http.Handler {
	return websocket.Handler(func(ws *websocket.Conn) {
		conn := NewWebSocketConnection(ws)
		// Get the underlying HTTP request
		upgradeHandler(conn, ws.Request())
	})
}
