package bus

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type DummyEvent struct {
	BaseEvent
}

type BusTest struct {
	suite.Suite
	bus *EventBus
}

func (s *BusTest) SetupTest() {
	s.bus = NewEventBus()
}

func getEvent() DummyEvent {
	return DummyEvent{
		BaseEvent{
			Payload: "hello",
			Name:    "Dummy",
		},
	}
}

func (s *BusTest) TestAddListeners() {
	listener := NewMockListener(s.T())
	s.Assert().Empty(s.bus.listeners)
	s.bus.AddListener(getEvent(), listener)
	s.Assert().Len(s.bus.listeners, 1)
}

func (s *BusTest) TestRemoveListeners() {
	listener := NewMockListener(s.T())
	event := getEvent()
	s.bus.AddListener(event, listener)
	s.bus.RemoveListener(event, listener)
	s.Assert().Empty(s.bus.listeners[event.ID()])
}

func (s *BusTest) TestDispatchSync() {
	listener := NewMockListener(s.T())
	event := getEvent()
	s.bus.AddListener(event, listener)
	listener.On("Handle", mock.Anything)

	s.bus.Dispatch(event)

	listener.AssertCalled(s.T(), "Handle", event)
}

func (s *BusTest) TestDispatchAsync() {
	listener := NewMockListener(s.T())
	event := getEvent()
	s.bus.AddListener(event, listener)
	listener.On("Handle", mock.Anything)

	s.bus.DispatchAsync(event)

	listener.AssertCalled(s.T(), "Handle", event)
}

func TestBusDispatcher(t *testing.T) {
	suite.Run(t, new(BusTest))
}
