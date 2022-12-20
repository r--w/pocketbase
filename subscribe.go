package pocketbase

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/SierraSoftworks/multicast/v2"
	"github.com/r3labs/sse/v2"
)

type Event[T any] struct {
	Action string `json:"action"`
	Record T      `json:"record"`
	Error  error  `json:"-"`
}

func (c Collection[T]) Subscribe(targets ...string) (*Stream[T], error) {
	if err := c.Authorize(); err != nil {
		return nil, err
	}

	if len(targets) == 0 {
		targets = []string{c.Name}
	}

	client := sse.NewClient(c.url + "/api/realtime")
	sch := make(chan *sse.Event)
	if err := client.SubscribeChanRaw(sch); err != nil {
		return nil, err
	}

	stream := newStream[T]()
	stream.unsubscribe = func() { client.Unsubscribe(sch) }

	handleSSEEvent := func(ev *sse.Event) {
		var e Event[T]
		e.Error = json.Unmarshal(ev.Data, &e)
		stream.channel.C <- e
	}

	once := &sync.Once{}
	stream.firstAuthResultLocker.Lock()
	go func() {
		for ev := range sch {
			switch string(ev.Event) {
			case "PB_CONNECT":
				err := c.authSubscribeStream(ev.Data, targets)
				once.Do(func() {
					defer stream.firstAuthResultLocker.Unlock()
					stream.firstAuthResult = err
				})
			default:
				go handleSSEEvent(ev)
			}
		}
	}()

	return stream, nil
}

type SubscriptionsSet struct {
	ClientID      string   `json:"clientId"`
	Subscriptions []string `json:"subscriptions"`
}

func (c Collection[T]) authSubscribeStream(data []byte, targets []string) (err error) {
	var s SubscriptionsSet
	if err = json.Unmarshal(data, &s); err != nil {
		return
	}
	s.Subscriptions = targets
	resp, err := c.client.R().SetBody(s).Post(c.url + "/api/realtime")
	if err != nil {
		return
	}
	if code := resp.StatusCode(); code != http.StatusNoContent {
		return fmt.Errorf("auth subscribe stream failed. resp status code is %v", code)
	}
	return
}

type Stream[T any] struct {
	channel     *multicast.Channel[Event[T]]
	unsubscribe func()

	onceCleanup *sync.Once

	firstAuthResultLocker *sync.RWMutex
	firstAuthResult       error
}

func newStream[T any]() *Stream[T] {
	return &Stream[T]{
		channel:     multicast.New[Event[T]](),
		onceCleanup: &sync.Once{},

		firstAuthResultLocker: &sync.RWMutex{},
	}
}

func (s *Stream[T]) Events() <-chan Event[T] {
	return s.channel.Listen().C
}

func (s *Stream[T]) Unsubscribe() {
	s.onceCleanup.Do(func() {
		s.unsubscribe()
		s.channel.Close()
	})
}

func (s *Stream[T]) WaitAuthReady() error {
	s.firstAuthResultLocker.RLock()
	defer s.firstAuthResultLocker.RUnlock()
	return s.firstAuthResult
}
