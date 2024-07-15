package eventbus

import "github.com/gopi-frame/contract/eventbus"

// Listener create listener from an anonymous function
func Listener(callback func(event eventbus.Event) error) eventbus.Listener {
	return &listener{callback}
}

type listener struct {
	callback func(event eventbus.Event) error
}

func (l *listener) Handle(event eventbus.Event) error {
	return l.callback(event)
}
