package graph

import (
	"fmt"
)

var (
	ErrListenerNotRegistred = fmt.Errorf("this listener is not registered")
)

var listenerRegistry = make(map[string]ListenerRegistration)

type HandleDeltaFunc func([]Delta) error

type ListenerRegistration struct {
	NewListenerFunc func(QuadStore) (Listener, error)
}

type Listener interface {
	Start()
	Stop()
}

func RegisterListener(name string, register ListenerRegistration) {
	if register.NewListenerFunc == nil {
		panic("NewListenerFunc must not be nil")
	}

	if _, found := listenerRegistry[name]; found {
		panic(fmt.Sprintf("Already registered listener %q.", name))
	}
	listenerRegistry[name] = register
}

func NewListener(name string, store QuadStore) (Listener, error) {
	r, registered := listenerRegistry[name]
	if !registered {
		return nil, ErrListenerNotRegistred
	}
	return r.NewListenerFunc(store)
}
