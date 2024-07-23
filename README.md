# EventBus
Package eventbus provides a simple event bus implementation.

## Installation

```shell
go get -u github.com/gopi-frame/eventbus
```

## Import
```go
import "github.com/gopi-frame/eventbus"
```

## Usage

```go
package main

import (
	"fmt"
	eventbuscontract "github.com/gopi-frame/contract/eventbus"
	"github.com/gopi-frame/eventbus"
)

type Event struct {
	Value string
}

func (e Event) Topic() string {
	return "event"
}

type AnotherEventWithSameTopic struct {}

func (e AnotherEventWithSameTopic) Topic() string {
	return "event"
}

// Listener is a simple listener can handle multiple events with same topic
type Listener struct{}

func (l *Listener) Handle(event eventbus.Event) error {
	// you shall assert the type of event firstly,
	// then do something
	e := event.(*Event)
	fmt.Println(e.Value)
	return nil
}

type Subscriber struct{}

func (s *Subscriber) Subscribe(bus eventbuscontract.Bus) {
	bus.Listen([]string{"event"}, eventbus.ListenFunc[*Event](s.OnEvent))
}

func (s *Subscriber) OnEvent(event *Event) error {
	fmt.Println(event.Value)
	return nil
}

func main() {
	bus := eventbus.New()
	// add an anonymous function as listener to topic "event" and specify the event type is *Event
	bus.Listen([]string{"event"}, bus.ListenFunc[*Event](func(event *Event) error {
		fmt.Println(event.Value)
		return nil
	})
	// add a listener to topic "event"
	bus.Listen([]string{"event"}, &Listener{})
	// add a subscriber
	bus.Subscribe(&Subscriber{})
	// publish an event to topic "event"
	if err := bus.Publish(&Event{Value: "Hello World"}); err != nil {
		panic(err)
	}
}
```