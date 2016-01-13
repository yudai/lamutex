package lamutex

import (
	"sync"
)

type LookAhead struct {
	locked   bool
	m        sync.Mutex
	requests chan chan bool
	reports  chan bool
}

func New(queue_length int) *LookAhead {
	return &LookAhead{
		locked:   false,
		requests: make(chan chan bool, queue_length),
		reports:  make(chan bool, 1),
	}
}

func (la *LookAhead) Lock() (ahead_report bool) {
	req := make(chan bool, 1)
	la.requests <- req
	la.tick()
	return <-req
}

func (la *LookAhead) Unlock(report bool) {
	la.reports <- report
	la.tick()
}

func (la *LookAhead) tick() {
	la.m.Lock()
	defer la.m.Unlock()

	if la.locked {
		select {
		case report := <-la.reports:
			if !report {
				num_flush := len(la.requests)
				for i := 0; i < num_flush; i++ {
					req := <-la.requests
					req <- false
				}
			}

			la.locked = false
		default:
			// do nothing
		}
	}

	if !la.locked {
		select {
		case req := <-la.requests:
			la.locked = true
			req <- true
		default:
			// do nothing
		}
	}
}
