package event

import (
	"api-demo/app/event/listener"
	"api-demo/internal/event"
	"api-demo/internal/global"
)

var listenerList = []event.ListenerInterface{
	listener.FooListener{},
	listener.BarListener{},
}

func init() {
	for _, l := range listenerList {
		global.EventDispatcher.Register(l)
	}
}
