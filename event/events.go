package event

import (
	"math/rand"
	"sync"
)

func NewEvent(id EventName, data interface{}) Event {
	return Event{
		name: EventName(id),
		data: data,
	}
}

// EventDispatcher ...
type EventDispatcher struct {
	listeners map[EventName][]Listener
	ids       map[int]bool
	mutex     *sync.Mutex
}

// NewEventDispatcher construct the event dispatcher
func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		mutex:     &sync.Mutex{},
		ids:       make(map[int]bool),
		listeners: make(map[EventName][]Listener),
	}
}

// AddListeners ...
func (a *EventDispatcher) AddListeners(name EventName, listeners ...Listener) {
	a.listeners[name] = append(a.listeners[name], listeners...)
	for _, listener := range listeners {
		id := a.getID()
		listener.SetId(id)
	}
}

func (a *EventDispatcher) getID() int {

	var id int
	a.mutex.Lock()
	for {
		id = rand.Int()
		exists := a.ids[id]
		if exists == false {
			a.ids[id] = true
			break
		}
	}
	a.mutex.Unlock()
	return id
}

// Dispatch emits an event
func (a *EventDispatcher) Dispatch(event Event) {
	listeners := a.listeners[event.name]
	for _, listener := range listeners {
		listener.Handle(event.data)
	}
}

// RemoveListener unsubscribes a listener from a particular events
func (a *EventDispatcher) RemoveListener(
	name EventName, listenerToRemove Listener,
) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	listeners := a.listeners[name]
	if len(listeners) < 1 {
		return
	}
	listenersCopy := listeners
	for i, listener := range listeners {
		if listener.GetId() == listenerToRemove.GetId() {
			listenersCopy = spliceSlice(listenersCopy, i)
		}
	}
	a.listeners[name] = listenersCopy
	delete(a.ids, listenerToRemove.GetId())
}

func spliceSlice(slice []Listener, index int) []Listener {
	if len(slice) == 1 {
		return slice[0:0]
	}
	slice[index] = slice[len(slice)-1]
	return slice[:len(slice)-1]
}
