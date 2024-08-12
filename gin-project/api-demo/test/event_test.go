package test

import (
	"testing"

	"api-demo/app/event/event"
	"api-demo/internal/global"
)

func TestEvent(t *testing.T) {
	global.EventDispatcher.Dispatch(&event.FooEvent{})
}
