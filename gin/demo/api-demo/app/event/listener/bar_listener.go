package listener

import (
	"fmt"

	appEvent "api-demo/app/event/event"
	"api-demo/internal/event"
)

type BarListener struct{}

func (f BarListener) Listen() []event.EventInterface {
	return []event.EventInterface{
		&appEvent.FooEvent{},
	}
}

func (f BarListener) Process(e event.EventInterface) {
	fmt.Println("bar listener process event:", e, e.Name())
}
