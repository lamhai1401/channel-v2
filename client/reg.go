package client

import (
	"fmt"
	"sync"
)

func toKey(topic, event, ref string) string {
	return fmt.Sprintf("KEY:%s:%s:%s:", topic, event, ref)
}

func newRegCenter() *RegCenter {
	return &RegCenter{
		regs: make(map[string][]*Puller),
	}
}

// RegCenter linter
type RegCenter struct {
	sync.Mutex
	regs map[string][]*Puller
}

// Register linter
func (center *RegCenter) Register(key string) *Puller {
	center.Lock()
	defer center.Unlock()

	puller := &Puller{
		key: key,
		ch:  make(chan *Message, maxMsgChannSize),
	}
	center.regs[key] = append(center.regs[key], puller)

	return puller
}

// Unregister linter
func (center *RegCenter) Unregister(puller *Puller) {
	center.Lock()
	defer center.Unlock()

	pullers := center.regs[puller.key]
	for i, _puller := range pullers {
		if _puller != puller {
			continue
		}
		pullers[i] = pullers[len(pullers)-1]
		center.regs[puller.key] = pullers[:len(pullers)-1]
		return
	}
}

// GetPullers linter
func (center *RegCenter) GetPullers(key string) []*Puller {
	center.Lock()
	defer center.Unlock()

	pullers := center.regs[key]
	copied := make([]*Puller, len(pullers))
	copy(copied, pullers)

	return copied
}

// CloseAllPullers linter
func (center *RegCenter) CloseAllPullers() {
	center.Lock()
	for id, pullers := range center.regs {
		delete(center.regs, id)
		center.closePuller(pullers)
	}
	center.Unlock()
}

func (center *RegCenter) closePuller(lst []*Puller) {
	for _, pull := range lst {
		center.Unregister(pull)
		pull.Close()
	}
}
