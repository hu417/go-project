package main

import (
	"fmt"

	"github.com/robfig/cron/v3"
)

func main() {
	c := cron.New()

	c.Start()
	defer c.Stop()

	_, err := c.AddFunc("0 30 * * * *", func() {
		fmt.Println("Every hour on the half hour")
	})
	_, err = c.AddFunc("@hourly", func() {
		fmt.Println("Every hour")
	})
	_, err = c.AddFunc("@every 1h30m", func() {
		fmt.Println("Every hour thirty")
	})

	_, err = c.AddFunc("@every 3s", func() {
		fmt.Println("@every 3s 执行开始")
	})

	_, err = c.AddFunc("@every 4s", func() {
		fmt.Println("@every 4s 执行开始")
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	select {}
}
