package eventbus

import (
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

// Bus event bus
type Bus struct {
	listeners *kv.Map[string, *list.List[eventbus.Listener]]
}

// Listen listen
func (d *Bus) Listen(events []eventbus.Event, listeners ...eventbus.Listener) {
	d.listeners.Lock()
	defer d.listeners.Unlock()
	for _, event := range events {
		if l, ok := d.listeners.Get(event.Topic()); ok {
			l.Push(listeners...)
		} else {
			d.listeners.Set(event.Topic(), list.NewList(listeners...))
		}
	}
}

// HasListener has listener
func (d *Bus) HasListener(event eventbus.Event) (exists bool) {
	d.listeners.RLock()
	defer d.listeners.RUnlock()
	listeners, ok := d.listeners.Get(event.Topic())
	if !ok {
		return false
	}
	return listeners.IsNotEmpty()
}

// Subscribe adds an subscriber
func (d *Bus) Subscribe(subscriber eventbus.Subscriber) {
	subscriber.Subscribe(d)
}

// Dispatch dispatches an event
func (d *Bus) Dispatch(e eventbus.Event) {
	d.listeners.RLock()
	defer d.listeners.RUnlock()
	listenerSet, ok := d.listeners.Get(e.Topic())
	if ok {
		listenerSet.Each(func(index int, listener eventbus.Listener) bool {
			if err := listener.Handle(e); err != nil {
				return false
			}
			return true
		})
	}
}