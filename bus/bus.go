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
	listeners := bus.listeners[event.ID()]
	for _, listener := range listeners {
		listener.Handle(event)
	}
}

func (bus *EventBus) DispatchAsync(event Event) {

	listeners := bus.listeners[event.ID()]
	var wg sync.WaitGroup
	for _, listener := range listeners {
		listener := listener.(Listener)
		wg.Add(1)
		go func(listener Listener) {
			listener.Handle(event)
			wg.Done()
		}(listener)
	}
	wg.Wait()

}

func (bus *EventBus) AddListener(event Event, listener Listener) {
	id := event.ID()
	bus.listeners[id] = append(bus.listeners[id], listener)
}

func (bus *EventBus) RemoveListener(event Event, listener Listener) {
	eventId := event.ID()
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
