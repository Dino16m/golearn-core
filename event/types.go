package event

// Listener interface must be implemented by all listeners
type Listener interface {
	Handle(interface{})
	SetId(int)
	GetId() int
}

type Event struct {
	name EventName
	data interface{}
}

// EventName custom type of events
type EventName string

// Dispatcher interface
type Dispatcher interface {
	AddListeners(EventName, ...Listener)
	// Dispatch calls event handlers in the same order as they were registered
	Dispatch(EventName, interface{})
	RemoveListener(EventName, Listener)
}

// BaseHandler fulfils the identification behavioour of an event handler
// leaving the user to just create the main handle method
type BaseHandler struct {
	id int
}

// SetId sets the identification id on the handler,
// this id allows the id to be trackable and removable
func (b *BaseHandler) SetId(id int) {
	b.id = id
}

// GetId returns the id of the handler
func (b *BaseHandler) GetId() int {
	return b.id
}
