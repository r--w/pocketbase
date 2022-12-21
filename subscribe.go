package pocketbase

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/SierraSoftworks/multicast/v2"
	"github.com/cenkalti/backoff/v4"
	"github.com/donovanhide/eventsource"
)

type Event[T any] struct {
	Action string `json:"action"`
	Record T      `json:"record"`
	Error  error  `json:"-"`
}

func (c Collection[T]) Subscribe(targets ...string) (*Stream[T], error) {
	opts := SubscribeOptions{
		ReconnectStrategy: backoff.NewConstantBackOff(time.Second),
	}
	return c.SubscribeWith(opts, targets...)
}

type SubscribeOptions struct {
	ReconnectStrategy backoff.BackOff
}

func (c Collection[T]) SubscribeWith(opts SubscribeOptions, targets ...string) (*Stream[T], error) {
	if err := c.Authorize(); err != nil {
		return nil, err
	}

	if len(targets) == 0 {
		targets = []string{c.Name}
	}

	stream := newStream[T]()
	ctx, cancel := context.WithCancel(context.Background())
	stream.unsubscribe = func() { cancel() }

	handleSSEEvent := func(ev eventsource.Event) {
		var e Event[T]
		e.Error = json.Unmarshal([]byte(ev.Data()), &e)
		stream.channel.C <- e
	}

	startStream := func(check bool) func() error {
		return func() (err error) {
			req := c.client.R().SetContext(ctx).SetDoNotParseResponse(true)
			resp, err := req.Get(c.url + "/api/realtime")
			defer resp.RawBody().Close()
			if err != nil {
				return
			}

			d := eventsource.NewDecoder(resp.RawBody())

			ev, err := d.Decode()
			if err != nil {
				return err
			}
			if event := ev.Event(); event != "PB_CONNECT" {
				return fmt.Errorf("first event must be PB_CONNECT, but got %s", event)
			}

			if err := c.authSubscribeStream([]byte(ev.Data()), targets); err != nil {
				return err
			}

			if !check {
				fmt.Println("loop")
				for {
					ev, err := d.Decode()
					if err != nil {
						return err
					}
					go handleSSEEvent(ev)
				}
			}

			return nil
		}
	}

	if err := startStream(true)(); err != nil {
		return nil, err
	}

	go func() {
		if err := backoff.Retry(startStream(false), backoff.WithContext(opts.ReconnectStrategy, ctx)); err != nil {
			log.Print(err)
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
}

func newStream[T any]() *Stream[T] {
	return &Stream[T]{
		channel:     multicast.New[Event[T]](),
		onceCleanup: &sync.Once{},
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

// Deprecated: now auth ready when Subscribe return nil error.
func (s *Stream[T]) WaitAuthReady() error { return nil }
