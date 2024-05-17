package event

import "github.com/gopi-frame/contract/event"

// Listener create listener from an anonymous function
func Listener(callback func(event event.Event) bool) event.Listener {
	return &listener{callback}
}

type listener struct {
	callback func(event event.Event) bool
}

func (l *listener) Handle(event event.Event) bool {
	return l.callback(event)
}
