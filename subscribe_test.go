package pocketbase

import (
	"testing"
	"time"

	"github.com/r--w/pocketbase/migrations"
	"github.com/stretchr/testify/assert"
)

func TestCollection_Subcribe(t *testing.T) {
	client := NewClient(defaultURL)
	defaultBody := map[string]interface{}{
		"field": "value_" + time.Now().Format(time.StampMilli),
	}
	collection := Collection[map[string]any]{client, migrations.PostsPublic}
	stream, err := collection.Subscribe()
	if err != nil {
		t.Error(err)
		return
	}
	if err := stream.WaitAuthReady(); err != nil {
		t.Error(err)
		return
	}

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

func TestCollection_Unsubcribe(t *testing.T) {
	client := NewClient(defaultURL)
	defaultBody := map[string]interface{}{
		"field": "value_" + time.Now().Format(time.StampMilli),
	}
	collection := Collection[map[string]any]{client, migrations.PostsPublic}
	stream, err := collection.Subscribe()
	if err != nil {
		t.Error(err)
		return
	}
	if err := stream.WaitAuthReady(); err != nil {
		t.Error(err)
		return
	}
	timer := time.NewTimer(3 * time.Second)
	defer timer.Stop()
	breaked := false
	go func() {
		<-timer.C
		if !breaked {
			panic("unsubscribe is not working, timeout")
		}
	}()

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

	for range ch {
	}
	breaked = true
}
