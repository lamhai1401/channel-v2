package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"

	"github.com/gorilla/websocket"
	"github.com/lamhai1401/gologs/logs"
)

// Connection linter
type Connection interface {
	Chan(topic string) (*Chan, error) //  Chan create a new channel on connection
	Push(msg *Message) error
	OnMessage() *Puller
	Close() error
	Start()
}

// RegisterCenter regis and unregis puller
type RegisterCenter interface {
	Register(key string) *Puller
	Unregister(puller *Puller)
	CloseAllPullers()
	GetPullers(key string) []*Puller
}

const maxMimumReadBuffer = 1024 * 1024 * 2
const maxMimumWriteBuffer = 1024 * 1024 * 2

// GetElixirURL linter
func GetElixirURL(_url string, args url.Values) (string, error) {
	surl, err := url.Parse(_url)

	if err != nil {
		return "", err
	}

	if !surl.IsAbs() {
		return "", errors.New("url should be absolute")
	}

	oscheme := surl.Scheme
	switch oscheme {
	case "ws":
		break
	case "wss":
		break
	case "http":
		surl.Scheme = "ws"
	case "https":
		surl.Scheme = "wss"
	default:
		return "", errors.New("schema should be http or https")
	}

	surl.Path = path.Join(surl.Path, "websocket")
	surl.RawQuery = args.Encode()

	// originURL := fmt.Sprintf("%s://%s", surl.Scheme, surl.Host)
	originURL := fmt.Sprintf("%s://%s%s?%s", surl.Scheme, surl.Host, surl.Path, surl.RawQuery)

	return originURL, nil
}

func openConnection(originURL string, dialer *websocket.Dialer) (*websocket.Conn, error) {
	// var wsConn *websocket.Conn

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(getTimeout())*time.Second)
	defer cancel()

	wsConn, _, err := dialer.DialContext(ctx, originURL, nil)
	if err != nil {
		if isTimeoutError(err) {
			logs.Warn("*** Connection timeout. Try to reconnect")
		}
		return nil, err
	}

	return wsConn, nil
}

// ConnectWss linter
func ConnectWss(_url string, args url.Values) (*websocket.Conn, error) {
	originURL, err := GetElixirURL(_url, args)
	if err != nil {
		return nil, err
	}

	ticket := time.NewTicker(500 * time.Millisecond)
	limit := 0

	var wsConn *websocket.Conn

	dialer := websocket.Dialer{
		Proxy:             http.ProxyFromEnvironment,
		HandshakeTimeout:  30 * time.Second,
		EnableCompression: true,
		ReadBufferSize:    maxMimumReadBuffer,
		WriteBufferSize:   maxMimumWriteBuffer,
	}

	for range ticket.C {
		if limit == 100 {
			fmt.Printf("Cannot connect %s 100 times. Call os.exit\n", originURL)
			os.Exit(0)
		}

		wsConn, err = openConnection(originURL, &dialer)
		if err != nil {
			limit++
			continue
		}
		break
	}

	err = wsConn.SetCompressionLevel(6)
	if err != nil {
		return nil, err
	}

	wsConn.EnableWriteCompression(true)
	return wsConn, nil
}
