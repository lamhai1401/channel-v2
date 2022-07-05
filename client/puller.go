package client

import "sync"

// Puller channel to receive result
type Puller struct {
	key    string
	ch     chan *Message // channel to handle all msg receive
	closed bool          // check close puller
	mutex  sync.Mutex
}

// Close return puller to regCenter
func (puller *Puller) Close() {
	puller.mutex.Lock()
	if puller.closed {
		return
	}

	puller.closed = true
	// puller.center.Unregister(puller)
	close(puller.ch) // maybe crash here

	puller.mutex.Unlock()
}

// Pull pull message from chan
func (puller *Puller) Pull() <-chan *Message {
	puller.mutex.Lock()
	defer puller.mutex.Unlock()
	return puller.ch
}

// isClosed to check puller close or not
func (puller *Puller) isClosed() bool {
	puller.mutex.Lock()
	defer puller.mutex.Unlock()
	return puller.closed
}

// Push linter
func (puller *Puller) Push(msg *Message) bool {
	if puller.isClosed() {
		return false
	}
	puller.ch <- msg
	return true
}
