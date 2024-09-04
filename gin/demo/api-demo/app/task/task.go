package task

import "api-demo/internal/crontab"

func Tasks() []crontab.TaskInterface {
	return []crontab.TaskInterface{
		&FooTask{},
	}
}
