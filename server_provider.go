package event

import (
	"reflect"

	"github.com/gopi-frame/contract/container"
	"github.com/gopi-frame/contract/event"
	"github.com/gopi-frame/contract/support"
)

var _ support.ServerProvider = (*ServerProvider)(nil)

// ServerProvider dispatcher server provider
type ServerProvider struct{}

// Register register
func (s *ServerProvider) Register(c container.Container) {
	c.Bind(reflect.TypeFor[Dispatcher]().String(), func(c container.Container) any {
		return NewDispatcher()
	})
	c.Alias(reflect.TypeFor[Dispatcher]().String(), "events")
	c.Alias(reflect.TypeFor[Dispatcher]().String(), reflect.TypeFor[event.Dispatcher]().Name())
}
