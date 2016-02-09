package lamutex

import (
	"sync/atomic"
)

type LookAhead struct {
	m      chan bool
	num    int64
	report bool
}

func New() *LookAhead {
	m := make(chan bool, 1)
	m <- true
	return &LookAhead{
		m:      m,
		num:    0,
		report: true,
	}
}

func (la *LookAhead) Lock() (ahead_report bool) {
	atomic.AddInt64(&la.num, 1)
	return <-la.m
}

func (la *LookAhead) Unlock(report bool) {
	if !report {
		l := la.num - 1
		var i int64
		for i = 0; i < l; i++ {
			la.m <- false
		}
		atomic.AddInt64(&la.num, -l)
	}
	atomic.AddInt64(&la.num, -1)
	la.m <- true
}
