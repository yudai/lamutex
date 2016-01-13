package lamutex

import (
	"sync"
	"testing"
	"time"
)

func TestLockTwoTimes(t *testing.T) {
	var r bool

	m := New(10)

	r = m.Lock()
	if r != true {
		t.Errorf("Failed to get the lock")
	}
	m.Unlock(true)

	r = m.Lock()
	if r != true {
		t.Errorf("Failed to get the lock")
	}
	m.Unlock(true)
}

func TestLockConcurrentlyTightly(t *testing.T) {
	m := New(3)

	num := 10
	wg := new(sync.WaitGroup)
	wg.Add(num)
	for i := 0; i < num; i++ {
		go func() {
			defer wg.Done()
			r := m.Lock()
			if r != true {
				t.Errorf("Got an unexpected false report")
			}
		}()
	}

	time.Sleep(time.Second * 1)

	for i := 0; i < num; i++ {
		m.Unlock(true)
	}

	wg.Wait()
}

func TestLockConcurrently(t *testing.T) {
	m := New(10)

	wg := new(sync.WaitGroup)
	wg.Add(3)
	for i := 0; i < 3; i++ {
		go func() {
			defer wg.Done()
			r := m.Lock()
			if r != true {
				t.Errorf("Got an unexpected false report")
			}
		}()
	}

	time.Sleep(time.Second * 1)
	m.Unlock(true)
	m.Unlock(true)
	m.Unlock(true)

	wg.Wait()

	m.Lock()
	wg.Add(3)
	for i := 0; i < 3; i++ {
		go func() {
			defer wg.Done()
			r := m.Lock()
			if r != false {
				t.Errorf("Got an unexpected true report")
			}
		}()
	}
	time.Sleep(time.Second * 1)
	m.Unlock(false)

	wg.Wait()
}
