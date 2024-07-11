package eventbus

import "github.com/gopi-frame/eventbus/contract"

// Listener create listener from an anonymous function
func Listener(callback func(event contract.Event) error) contract.Listener {
	return &listener{callback}
}

type listener struct {
	callback func(event contract.Event) error
}

func (l *listener) Handle(event contract.Event) error {
	return l.callback(event)
}
