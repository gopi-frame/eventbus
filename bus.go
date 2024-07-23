package eventbus

import (
	"errors"
	"github.com/gopi-frame/collection/kv"
	"github.com/gopi-frame/collection/list"
	"github.com/gopi-frame/contract/eventbus"
)

// NewBus creates a new [Bus] instance
func NewBus() *Bus {
	return &Bus{
		listeners: kv.NewMap[string, *list.List[eventbus.Listener]](),
	}
}

// Bus event bus.
type Bus struct {
	listeners *kv.Map[string, *list.List[eventbus.Listener]]
}

// Listen adds listeners for given topics.
func (d *Bus) Listen(topics []string, listeners ...eventbus.Listener) {
	d.listeners.Lock()
	defer d.listeners.Unlock()
	for _, topic := range topics {
		if l, ok := d.listeners.Get(topic); ok {
			l.Push(listeners...)
		} else {
			d.listeners.Set(topic, list.NewList(listeners...))
		}
	}
}

// HasListener returns true if there is a listener for given topic else false.
func (d *Bus) HasListener(topic string) (exists bool) {
	d.listeners.RLock()
	defer d.listeners.RUnlock()
	listeners, ok := d.listeners.Get(topic)
	if !ok {
		return false
	}
	return listeners.IsNotEmpty()
}

// Subscribe adds a subscriber
func (d *Bus) Subscribe(subscriber eventbus.Subscriber) {
	subscriber.Subscribe(d)
}

// Dispatch dispatches an event and returns an error if any.
func (d *Bus) Dispatch(e eventbus.Event) error {
	d.listeners.RLock()
	defer d.listeners.RUnlock()
	listenerSet, ok := d.listeners.Get(e.Topic())
	var errs []error
	if ok {
		listenerSet.Each(func(index int, listener eventbus.Listener) bool {
			if err := listener.Handle(e); err != nil {
				errs = append(errs, err)
			}
			return true
		})
	}
	return errors.Join(errs...)
}
