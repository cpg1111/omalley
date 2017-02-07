package elect

import (
	"net"
	"sync"
	"time"

	"github.com/pullrequestrfb/omalley/action"
)

type Elector struct {
	Channel    chan *action.Action
	Candidates map[string]int
	lock       *sync.Mutex
}

func New(localChan chan *action.Action) *Elector {
	return &Elector{
		Channel:    localChan,
		Candidates: make(map[string]int),
	}
}

func (e *Elector) Recv(conn *net.Conn, msg map[string]string) (bool, error) {
	if msg["candidate"] != nil {
		e.lock.Lock()
		defer e.lock.Unlock()
		e.Candidates[msg["candidate"]]++
	}
}

func (e *Elector) Vote(vDispatcher Dispatcher) error {
	act := <-e.Channel
	return vDispatcher.DispatchVote(act)
}
