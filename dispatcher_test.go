package event

import (
	"testing"

	"github.com/gopi-frame/contract/event"
	"github.com/stretchr/testify/assert"
)

type _event struct{ message string }

func (e *_event) Topic() string { return "topic1" }

type _subscriber struct{}

func (s *_subscriber) onTopic1(event event.Event) bool {
	return true
}

func (s *_subscriber) Subscribe(d event.Dispatcher) {
	d.Listen([]event.Event{new(_event)}, Listener(s.onTopic1))
}

func TestDispatcher(t *testing.T) {
	_dispatcher := NewDispatcher()
	var message string
	_dispatcher.Listen([]event.Event{new(_event)}, Listener(func(event event.Event) bool {
		message = event.(*_event).message
		return true
	}))
	_dispatcher.Subscribe(new(_subscriber))
	assert.True(t, _dispatcher.HasListener(new(_event)))
	_dispatcher.Dispatch(&_event{message: "hello"})
	assert.Equal(t, "hello", message)
}
