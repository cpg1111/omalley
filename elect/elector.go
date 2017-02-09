package elect

import (
	"net"
	"sync"

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

func (e *Elector) Recv(conn *net.TCPConn, msg map[string]string) (bool, error) {
	if len(msg["candidate"]) > 0 {
		e.lock.Lock()
		defer e.lock.Unlock()
		e.Candidates[msg["candidate"]]++
	}
	return true, nil
}

func (e *Elector) Confirm(conn *net.TCPConn, msg map[string]string) (bool, error) {
	return true, nil
}

func (e *Elector) Vote(vDispatcher Dispatcher) error {
	act := <-e.Channel
	return vDispatcher.DispatchVote(act)
}
