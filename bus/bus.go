package bus

import (
	"reflect"
	"sync"
)

type EventBus struct {
	listeners map[string][]Listener
}

func NewEventBus() *EventBus {
	return &EventBus{
		listeners: make(map[string][]Listener),
	}
}

func (bus *EventBus) Dispatch(event Event) {
	listeners := bus.listeners[reflect.TypeOf(event).Name()]
	for _, listener := range listeners {
		listener.Handle(event)
	}
}

func (bus *EventBus) DispatchAsync(event Event) {

	listeners := bus.listeners[reflect.TypeOf(event).Name()]
	var wg sync.WaitGroup
	for _, listener := range listeners {
		wg.Add(1)
		go func(listener Listener) {
			listener.Handle(event)
			wg.Done()
		}(listener)
	}
	wg.Wait()

}

// AddListener will not work if the added event is a pointer.
// It only works with value types
func (bus *EventBus) AddListener(event Event, listener Listener) {
	id := reflect.TypeOf(event).Name()

	bus.listeners[id] = append(bus.listeners[id], listener)
}

func (bus *EventBus) RemoveListener(event Event, listener Listener) {
	eventId := reflect.TypeOf(event).Name()
	listeners := bus.listeners[eventId]
	match := -1
	for index, existingListener := range listeners {
		if reflect.DeepEqual(existingListener, listener) {
			match = index
			break
		}
	}
	if match == -1 {
		return
	}
	bus.listeners[eventId] = spliceSlice(listeners, match)

}

func spliceSlice(slice []Listener, index int) []Listener {
	if len(slice) == 1 {
		return slice[0:0]
	}
	slice[index] = slice[len(slice)-1]
	return slice[:len(slice)-1]
}
