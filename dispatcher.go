package event

import (
	"github.com/gopi-frame/contract/event"
	"github.com/gopi-frame/support/lists"
	"github.com/gopi-frame/support/maps"
)

var _ event.Dispatcher = (*Dispatcher)(nil)

// NewDispatcher creates a new [Dispatcher] instance
func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		listeners: maps.NewMap[string, *lists.List[event.Listener]](),
	}
}

// Dispatcher dispatcher
type Dispatcher struct {
	listeners *maps.Map[string, *lists.List[event.Listener]]
}

// Listen listen
func (d *Dispatcher) Listen(events []event.Event, listeners ...event.Listener) {
	d.listeners.Lock()
	defer d.listeners.Unlock()
	for _, event := range events {
		if list, ok := d.listeners.Get(event.Topic()); ok {
			list.Push(listeners...)
		} else {
			d.listeners.Set(event.Topic(), lists.NewList(listeners...))
		}
	}
}

// HasListener has listener
func (d *Dispatcher) HasListener(event event.Event) (exists bool) {
	d.listeners.Lock()
	defer d.listeners.Unlock()
	listeners, ok := d.listeners.Get(event.Topic())
	if !ok {
		exists = false
		return
	}
	exists = listeners.IsNotEmpty()
	return
}

// Subscribe adds an subscriber
func (d *Dispatcher) Subscribe(subscriber event.Subscriber) {
	subscriber.Subscribe(d)
}

// Dispatch dispatches an event
func (d *Dispatcher) Dispatch(e event.Event) {
	d.listeners.Lock()
	defer d.listeners.Unlock()
	listenerSet := d.listeners.GetOr(e.Topic(), lists.NewList[event.Listener]())
	listenerSet.Each(func(index int, listener event.Listener) bool {
		return listener.Handle(e)
	})
}
