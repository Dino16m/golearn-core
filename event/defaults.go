package event

const (
	// expects map[string]interface{} with keys firstName, lastName, email,
	// and id which should have an int value
	UserCreated = EventName("usercreated")
)
