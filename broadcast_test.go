package broadcast

import (
	"sync"
	"testing"
	"time"
)

const (
	N       = 3
	testStr = "Test"
	timeout = time.Second
)

type ListenFunc func(int, *Broadcaster, *sync.WaitGroup)

func setupN(f ListenFunc) (*Broadcaster, *sync.WaitGroup) {
	var b Broadcaster
	var wg sync.WaitGroup
	wg.Add(N)
	for i := 0; i < N; i++ {
		go f(i, &b, &wg)
	}
	wg.Wait()
	return &b, &wg
}

func TestSend(t *testing.T) {
	b, wg := setupN(func(i int, b *Broadcaster, wg *sync.WaitGroup) {
		l := b.Listen()
		wg.Done()
		select {
		case v := <-l.Ch:
			if v.(string) != testStr {
				t.Error("bad value received")
			}
		case <-time.After(timeout):
			t.Error("receive timed out")
		}
		wg.Done()
	})
	wg.Add(N)
	b.Send(testStr)
	wg.Wait()
}

func TestBroadcasterClose(t *testing.T) {
	b, wg := setupN(func(i int, b *Broadcaster, wg *sync.WaitGroup) {
		l := b.Listen()
		wg.Done()
		select {
		case _, ok := (<-l.Ch):
			if ok {
				t.Error("receive after close")
			}
		case <-time.After(timeout):
			t.Error("receive timed out")
		}
		wg.Done()
	})
	wg.Add(N)
	b.Close()
	wg.Wait()
}

func TestListenerClose(t *testing.T) {
	b, wg := setupN(func(i int, b *Broadcaster, wg *sync.WaitGroup) {
		l := b.Listen()
		if i == 0 {
			l.Close()
		}
		wg.Done()
		select {
		case <-l.Ch:
			if i == 0 {
				t.Error("receive after close")
			}
		case <-time.After(timeout):
			if i != 0 {
				t.Error("receive timed out")
			}
		}
		wg.Done()
	})
	wg.Add(N)
	b.Send(testStr)
	wg.Wait()
}
