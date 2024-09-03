package eventbus

import (
	"errors"
	"github.com/gopi-frame/contract/eventbus"
	"github.com/stretchr/testify/assert"
	"testing"
)

type testEvent struct {
	data any
}

func (e *testEvent) Topic() string {
	return "testEvent"
}

type testEventWithSameTopic struct {
	data any
}

func (e *testEventWithSameTopic) Topic() string {
	return "testEvent"
}

type testEventWillCauseError struct {
}

func (e *testEventWillCauseError) Topic() string {
	return "testEvent"
}

type testListener struct{}

func (l *testListener) Handle(event eventbus.Event) error {
	if e, ok := event.(*testEvent); ok {
		e.data = "testEvent handled"
	} else if e, ok := event.(*testEventWithSameTopic); ok {
		e.data = "testEventWithSameTopic handled"
	} else {
		return errors.New("unknown event")
	}
	return nil
}

type testSubscriber struct {
}

func (s *testSubscriber) Subscribe(bus eventbus.Bus) {
	bus.Listen([]string{"testEvent"}, ListenFunc[*testEvent](s.OnEvent))
}

func (s *testSubscriber) OnEvent(event *testEvent) error {
	event.data = "testEvent handled"
	return nil
}

func TestBus_Listen(t *testing.T) {
	bus := NewBus()
	assert.False(t, bus.HasListener("testEvent"))
	bus.Listen([]string{"testEvent"}, ListenFunc[*testEvent](func(event *testEvent) error {
		event.data = "testEvent handled"
		return nil
	}))
	assert.True(t, bus.HasListener("testEvent"))

	t.Run("append listeners", func(t *testing.T) {
		bus.Listen([]string{"testEvent"}, ListenFunc[*testEvent](func(event *testEvent) error {
			event.data = "testEvent handled"
			return nil
		}))
		assert.True(t, bus.HasListener("testEvent"))
	})
}

func TestBus_Subscribe(t *testing.T) {
	bus := NewBus()
	assert.False(t, bus.HasListener("testEvent"))
	bus.Subscribe(&testSubscriber{})
	assert.True(t, bus.HasListener("testEvent"))
}

func TestBus_Dispatch(t *testing.T) {
	t.Run("dispatch to listener", func(t *testing.T) {
		bus := NewBus()
		bus.Listen([]string{"testEvent"}, new(testListener))
		t.Run("success", func(t *testing.T) {
			e := new(testEvent)
			if err := bus.Dispatch(e); err != nil {
				assert.FailNow(t, err.Error())
			}
			assert.Equal(t, "testEvent handled", e.data)
			e2 := new(testEventWithSameTopic)
			if err := bus.Dispatch(e2); err != nil {
				assert.FailNow(t, err.Error())
			}
			assert.Equal(t, "testEventWithSameTopic handled", e2.data)
		})

		t.Run("exception", func(t *testing.T) {
			e := new(testEventWillCauseError)
			if err := bus.Dispatch(e); err == nil {
				assert.FailNow(t, "should have returned exception")
			} else {
				assert.Equal(t, "unknown event", err.Error())
			}
		})
	})

	t.Run("dispatch to subscriber", func(t *testing.T) {
		bus := NewBus()
		bus.Subscribe(&testSubscriber{})
		e := new(testEvent)
		if err := bus.Dispatch(e); err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, "testEvent handled", e.data)

		e2 := new(testEventWithSameTopic)
		if err := bus.Dispatch(e2); err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Nil(t, e2.data)
	})
}
