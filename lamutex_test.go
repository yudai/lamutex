package lamutex

import (
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestLockTwoTimes(t *testing.T) {
	var r bool

	m := New()

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
	m := New()

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
	m := New()

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

func TestLockConcurrently100000(t *testing.T) {
	m := New()

	n := 100000
	rand.Seed(time.Now().Unix())
	wg := new(sync.WaitGroup)
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(i int) {
			defer wg.Done()
			time.Sleep(time.Millisecond * time.Duration(rand.Int31n(10000)))
			t.Logf("Taking lock #%d\n", i)
			r := m.Lock()
			if r != true {
				t.Logf("Ahead failed #%d\n", i)
			} else {
				t.Logf("Locking #%d\n", i)
				time.Sleep(time.Millisecond * time.Duration(rand.Int31n(1000)))
				report := rand.Int31n(10) < 8
				t.Logf("Releasing with %t #%d\n", report, i)
				m.Unlock(report)
			}
		}(i)
	}

	time.Sleep(time.Second * 1)

	wg.Wait()

	if len(m.m) != 1 {
		t.Errorf("Inconsistent count")
	}
	if m.num != 0 {
		t.Errorf("Inconsistent count")
	}
}
