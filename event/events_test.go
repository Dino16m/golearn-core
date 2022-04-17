package event

import (
	"testing"

	"github.com/dino16m/golearn-core/event/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type dummyPayload struct{}

type authDispatcherTestSuite struct {
	suite.Suite
	dispatcher *EventDispatcher
	eventName  EventName
}

func (s *authDispatcherTestSuite) SetupTest() {
	s.dispatcher = NewEventDispatcher()
	s.eventName = EventName("dummy")
}

func (s *authDispatcherTestSuite) TestRegisterListenerSetsIdOnListener() {
	listener := new(mocks.Listener)
	listener.On("SetId", mock.AnythingOfType("int"))
	s.dispatcher.AddListeners(s.eventName, listener)
	listener.AssertExpectations(s.T())
}

func (s *authDispatcherTestSuite) TestRegisteredListenerCalledOnDispatch() {
	listener := new(mocks.Listener)
	listener.On("SetId", mock.AnythingOfType("int"))
	listener.On("Handle", mock.AnythingOfType("dummyPayload"))
	s.dispatcher.AddListeners(s.eventName, listener)
	event := NewEvent(s.eventName, dummyPayload{})
	s.dispatcher.Dispatch(event)
	listener.AssertExpectations(s.T())
}

func (s *authDispatcherTestSuite) TestListenerNotCalledWhenRemoved() {
	payload := dummyPayload{}
	listener := new(mocks.Listener)
	var captured int
	listener.On("SetId", mock.MatchedBy(func(id int) bool {
		captured = id
		return true
	}))

	secondListener := new(mocks.Listener)
	var secondCapture int
	secondListener.On("SetId", mock.MatchedBy(func(id int) bool {
		secondCapture = id
		return true
	}))
	s.dispatcher.AddListeners(s.eventName, listener, secondListener)
	secondListener.On("GetId").Return(secondCapture)
	secondListener.On("Handle", mock.Anything)
	listener.On("GetId").Return(captured)
	s.dispatcher.RemoveListener(s.eventName, listener)
	s.dispatcher.Dispatch(NewEvent(s.eventName, payload))
	listener.AssertExpectations(s.T())
	listener.AssertNotCalled(s.T(), "Handle", payload)
}

func TestEventDispatcher(t *testing.T) {
	suite.Run(t, new(authDispatcherTestSuite))
}
