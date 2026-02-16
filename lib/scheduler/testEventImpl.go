package scheduler

import (
	"fmt"
	"log"
	"time"
)

type TestEvent struct {
	BaseEvent
}

func NewTestEvent() *TestEvent {
	testEvent := &TestEvent{
		BaseEvent: BaseEvent{
			name:         "test_event",
			version:      "0.1",
			description:  "test event",
			capabilities: []string{"log"},
		},
	}

	if err := Register(testEvent); err != nil {
		fmt.Printf("NewTestEvent: cant register event: %s", err)
	}
	return testEvent
}

func (te *TestEvent) Run(data interface{}) (interface{}, error) {
	start := time.Now()

	if data == nil {
		return nil, fmt.Errorf("data cannot be nil")
	}
	log.Printf("log test event start at: %s", start)
	return data, nil
}
