package eventbus

import "github.com/gopi-frame/contract/eventbus"

// ListenFunc constructs a listener from a callback function with the event type.
// If the event type can not be asserted to the type of the callback function, then the callback will not be called.
func ListenFunc[T eventbus.Event](callback func(event T) error) eventbus.Listener {
	return &listener{func(event eventbus.Event) error {
		if event, ok := event.(T); ok {
			return callback(event)
		}
		return nil
	}}
}

func Listener[T eventbus.Event](eventListener eventbus.EventListener[T]) eventbus.Listener {
	return &listener{func(event eventbus.Event) error {
		return eventListener.Handle(event.(T))
	}}
}

type listener struct {
	callback func(event eventbus.Event) error
}

func (l *listener) Handle(event eventbus.Event) error {
	return l.callback(event)
}
