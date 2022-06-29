package client

import "sync"

// Puller channel to receive result
type Puller struct {
	key    string
	ch     chan *Message // channel to handle all msg receive
	center RefCenter     // to remove and regis
	closed bool          // check close puller
	mutex  sync.Mutex
}

// Close return puller to regCenter
func (puller *Puller) Close() {
	puller.mutex.Lock()
	if puller.center == nil || puller.closed {
		return
	}

	puller.closed = true
	puller.center.Unregister(puller)
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
