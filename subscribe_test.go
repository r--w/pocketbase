package pocketbase

import (
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/r--w/pocketbase/migrations"
	"github.com/stretchr/testify/assert"
)

func TestCollection_Subscribe(t *testing.T) {
	client := NewClient(defaultURL)
	defaultBody := map[string]interface{}{
		"field": "value_" + time.Now().Format(time.StampMilli),
	}
	collection := CollectionSet[map[string]any](client, migrations.PostsPublic)
	stream, err := collection.Subscribe()
	if err != nil {
		t.Error(err)
		return
	}
	defer stream.Unsubscribe()
	<-stream.Ready()

	ch := stream.Events()

	t.Run("subscribe event: create", func(t *testing.T) {
		resp, err := collection.Create(defaultBody)
		if err != nil {
			t.Error(err)
			return
		}
		e := <-ch
		assert.Equal(t, "create", e.Action)
		assert.Equal(t, resp.ID, e.Record["id"])
	})

	t.Run("subscribe event: update", func(t *testing.T) {
		resp, err := collection.Create(defaultBody)
		if err != nil {
			t.Error(err)
			return
		}
		<-ch // ignore create event
		body := map[string]interface{}{
			"field": "value_" + time.Now().Format(time.StampMilli),
		}
		err = collection.Update(resp.ID, body)
		if err != nil {
			t.Error(err)
			return
		}
		e := <-ch
		assert.Equal(t, "update", e.Action)
		assert.Equal(t, body["field"], e.Record["field"])
	})

	t.Run("subscribe event: delete", func(t *testing.T) {
		resp, err := collection.Create(defaultBody)
		if err != nil {
			t.Error(err)
			return
		}
		<-ch // ignore create event
		err = collection.Delete(resp.ID)
		if err != nil {
			t.Error(err)
			return
		}
		e := <-ch
		assert.Equal(t, "delete", e.Action)
		assert.Equal(t, resp.ID, e.Record["id"])
	})
}

func TestCollection_Unsubscribe(t *testing.T) {
	client := NewClient(defaultURL)
	defaultBody := map[string]interface{}{
		"field": "value_" + time.Now().Format(time.StampMilli),
	}
	collection := CollectionSet[map[string]any](client, migrations.PostsPublic)
	stream, err := collection.Subscribe()
	if err != nil {
		t.Error(err)
		return
	}
	<-stream.Ready()

	ch := stream.Events()

	resp, err := collection.Create(defaultBody)
	if err != nil {
		t.Error(err)
		return
	}
	e := <-ch
	assert.Equal(t, resp.ID, e.Record["id"])

	stream.Unsubscribe()

	if err := collection.Delete(resp.ID); err != nil {
		t.Error(err)
		return
	}

	if _, ok := <-ch; ok {
		t.Error("unsubscribe is not working.")
	}
}

func TestCollection_RealtimeReconnect(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping realtime reconnect in short mode")
		return
	}

	client := NewClient(defaultURL)
	transport := &http.Transport{
		Dial: func(network, addr string) (net.Conn, error) {
			conn, err := net.Dial(network, addr)
			if err == nil {
				// Simulate pocketbase closing realtime connection after 5m of inactivity
				time.AfterFunc(3*time.Second, func() { conn.Close() })
			}
			return conn, err
		},
	}
	client.client.SetTransport(transport)
	defaultBody := map[string]interface{}{
		"field": "value_" + time.Now().Format(time.StampMilli),
	}
	collection := CollectionSet[map[string]any](client, migrations.PostsPublic)
	stream, err := collection.Subscribe()
	if err != nil {
		t.Error(err)
		return
	}
	defer stream.Unsubscribe()
	<-stream.Ready()

	var got = false
	time.AfterFunc(13*time.Second, func() {
		if _, err := collection.Create(defaultBody); err != nil {
			t.Error(err)
		}
	})
	time.AfterFunc(15*time.Second, func() {
		defer stream.Unsubscribe()
		if !got {
			panic("stream reconnect is not working")
		}
	})

	for range stream.Events() {
		got = true
		break
	}
	assert.Equal(t, true, got)
}
