package client

const (
	// Channel status

	// ChClosed channel closed
	ChClosed = "closed"
	// ChErrored channel errored
	ChErrored = "errored"
	// ChJoined channel joined
	ChJoined = "joined"
	// ChJoing channel joining
	ChJoing = "joining"

	// Channel events

	// ChanClose channel close event
	ChanClose = "phx_close"
	// ChanError channel error event
	ChanError = "phx_error"
	// ChanReply server reply event
	ChanReply = "phx_reply"
	// ChanJoin  client send join event
	ChanJoin = "phx_join"
	// ChanLeave client send leave event
	ChanLeave = "phx_leave"
)

// Wss linter
type Wss interface {
	Push(msg *Message) error
	Register(key string) *Puller
	MakeRef() string
}

// Chan linter
type Chan struct {
	conn   Wss
	topic  string
	status string
}

// OnMessage register a MsgCh to recv all msg on channel
func (ch *Chan) OnMessage() *Puller {
	key := toKey(ch.topic, "", "")
	return ch.conn.Register(key)
}

// Join channel join, return a MsgCh to receive join result
func (ch *Chan) Join() (*Puller, error) {
	return ch.Request(ChanJoin, "")
}

// Leave channel leave, return a MsgCh to receive leave result
func (ch *Chan) Leave() (*Puller, error) {
	return ch.Request(ChanLeave, "")
}

// OnEvent return a MsgCh to receive all msg on some event on this channel
func (ch *Chan) OnEvent(evt string) *Puller {
	key := toKey(ch.topic, evt, "")
	return ch.conn.Register(key)
}

// Request send a msg to channel and return a MsgCh to receive reply
func (ch *Chan) Request(evt string, payload interface{}) (*Puller, error) {
	msg := &Message{
		Topic:   ch.topic,
		Event:   evt,
		Ref:     ch.conn.MakeRef(),
		Payload: payload,
	}
	// If we're receiving a reply, all we care about is that message references
	// Match. Connection.dispatch has a dispatcher for this key pattern ("","",msg.Ref)
	key := toKey("", "", msg.Ref)
	puller := ch.conn.Register(key)
	if err := ch.conn.Push(msg); err != nil {
		puller.Close()
		return nil, err
	}
	return puller, nil
}

// Push send a msg to channel
func (ch *Chan) Push(evt string, payload interface{}) error {
	msg := &Message{
		Topic:   ch.topic,
		Event:   evt,
		Ref:     ch.conn.MakeRef(),
		Payload: payload,
	}
	return ch.conn.Push(msg)
}

// PushWithRef send a msg to channel
func (ch *Chan) PushWithRef(ref, evt *string, payload interface{}) error {
	msg := &Message{
		Topic:   ch.topic,
		Event:   *evt,
		Ref:     *ref,
		Payload: payload,
	}
	return ch.conn.Push(msg)
}

// GetRef linter
func (ch *Chan) GetRef() string {
	return ch.conn.MakeRef()
}
