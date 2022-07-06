package client

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/lamhai1401/gologs/logs"
)

const all = ""

const (
	// ConnConnecting linter
	ConnConnecting = "connecting"
	// ConnOpen linter
	ConnOpen = "open"
	// ConnClosing linter
	ConnClosing = "closing"
	// Connclosed linter
	Connclosed = "closed"
)

// Connect linter
func Connect(_url string, args url.Values) (Connection, error) {
	conn, err := ConnectWss(_url, args)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	cl := &Client{
		socket: &WSocket{
			conn: conn,
		},
		ctx:       ctx,
		cancel:    cancel,
		status:    ConnOpen,
		msgs:      make(chan *Message, maxMsgChannSize),
		ref:       new(refMaker),
		refCenter: newRegCenter(),
	}

	cl.Start()
	return cl, nil
}

// Socket linter
type Socket interface {
	Send(*Message) error
	Recv() (*Message, error)
	Close() error
}

// Client linter
type Client struct {
	ref       *refMaker
	refCenter RegisterCenter // register and unregister puller
	socket    Socket         // send and recv msg
	msgs      chan *Message
	ctx       context.Context
	cancel    func()
	status    string

	mutex sync.Mutex
}

// OnMessage receive all message on connection
func (conn *Client) OnMessage() *Puller {
	return conn.refCenter.Register(all)
}

// Close linter
func (conn *Client) Close() error {
	conn.mutex.Lock()
	conn.status = ConnClosing
	conn.mutex.Unlock()
	conn.cancel()
	return conn.socket.Close()
}

// Push linter
func (conn *Client) Push(msg *Message) error {
	return conn.socket.Send(msg)
}

// Start linter
func (conn *Client) Start() {
	go conn.pullLoop()
	go conn.coreLoop()
	go conn.heartbeatLoop()
}

func (conn *Client) heartbeatLoop() {
	msg := &Message{
		Topic:   "phoenix",
		Event:   "heartbeat",
		Payload: "",
	}
	for {
		select {
		case <-time.After(30000 * time.Millisecond):
			// Set the message reference right before we send it:
			msg.Ref = conn.ref.makeRef()
			conn.Push(msg)
		case <-conn.ctx.Done():
			return
		}
	}
}

func (conn *Client) pullLoop() {
	for {
		msg, err := conn.socket.Recv()
		if err != nil {
			conn.msgs <- nil
			logs.Error(fmt.Sprintf("%s\n", err))
			close(conn.msgs)
			conn.refCenter.CloseAllPullers() // TODO fix here
			return
		}
		select {
		case <-conn.ctx.Done():
			return
		case conn.msgs <- msg:
		}
	}
}

func (conn *Client) coreLoop() {
	for {
		select {
		case <-conn.ctx.Done():
			return
		case msg, ok := <-conn.msgs:
			if !ok || msg == nil {
				return
			}
			conn.dispatch(msg)
		}
	}
}

func (conn *Client) dispatch(msg *Message) {
	var wg sync.WaitGroup
	wg.Add(4)
	go conn.pushToChans(&wg, conn.refCenter.GetPullers(all), msg)
	go conn.pushToChans(&wg, conn.refCenter.GetPullers(toKey(msg.Topic, "", "")), msg)
	go conn.pushToChans(&wg, conn.refCenter.GetPullers(toKey(msg.Topic, msg.Event, "")), msg)
	go conn.pushToChans(&wg, conn.refCenter.GetPullers(toKey("", "", msg.Ref)), msg)
	wg.Wait()
}

func (conn *Client) pushToChans(wg *sync.WaitGroup, pullers []*Puller, msg *Message) {
	defer wg.Done()
	if len(pullers) == 0 {
		return
	}
	for _, puller := range pullers {
		puller.Push(msg)
	}
}

// MakeRef linter
func (conn *Client) MakeRef() string {
	return conn.ref.makeRef()
}

// Register linter
func (conn *Client) Register(key string) *Puller {
	return conn.refCenter.Register(key)
}

// Chan linter
func (conn *Client) Chan(topic string) (*Chan, error) {
	if conn.status != ConnOpen {
		return nil, errors.New("Connection is " + conn.status)
	}

	ch := &Chan{
		conn:   conn,
		topic:  topic,
		status: ChJoing,
	}

	return ch, nil
}
