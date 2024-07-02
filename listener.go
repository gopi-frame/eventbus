package event

import "github.com/gopi-frame/event/contract"

// Listener create listener from an anonymous function
func Listener(callback func(event contract.Event) bool) contract.Listener {
	return &listener{callback}
}

type listener struct {
	callback func(event contract.Event) bool
}

func (l *listener) Handle(event contract.Event) bool {
	return l.callback(event)
}
