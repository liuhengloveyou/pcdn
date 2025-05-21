package main

import (
	"fmt"

	"github.com/robfig/cron/v3"
)

func InitTasks() {
	c := cron.New(cron.WithSeconds())

	// _, err := c.AddFunc("*/1 * * * * *", func() {
	// 	fmt.Println("Every second")
	// })
	// if err != nil {
	// 	return
	// }

	// 定时更新
	if _, err := c.AddFunc("1 1 1 * * *", func() {
		if err := updater.BackgroundRun(); err != nil {
			fmt.Println("Failed to update app:", err)
		}
	}); err != nil {
		return
	}

	c.Start()

	fmt.Println("Starting cron")
}
