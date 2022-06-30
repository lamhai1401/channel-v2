package client

import (
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Wss constant
const (
	// Time allowed to write a message to the peer.
	writeWait = 30 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second
)

// WSocket Socket implementation for web socket
type WSocket struct {
	conn  *websocket.Conn
	mutex sync.Mutex
}

// Close implements Socket.Close
func (ws *WSocket) Close() error {
	return ws.conn.Close()
}

// Recv implements Socket.Recv
func (ws *WSocket) Recv() (*Message, error) {
	var msg *Message
	// _, resp, err := ws.conn.ReadMessage()
	err := ws.conn.ReadJSON(&msg)
	if err != nil {
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
			return nil, fmt.Errorf("websocket closeGoingAway error: %v", err)
		}
		return nil, fmt.Errorf("recv err: %v", err)
	}

	if err != nil {
		return nil, fmt.Errorf("signaler Unmarshal err: %v", err)
	}
	return msg, nil
}

// Send implements Socket.Send
func (ws *WSocket) Send(msg *Message) error {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()
	if err := ws.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
		return err
	}

	if err := ws.conn.WriteJSON(msg); err != nil {
		return err
	}
	return nil
}
