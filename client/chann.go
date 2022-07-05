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

// Chan linter
type Chan struct {
	conn   Connection
	topic  string
	status string
}
