package tasks

import (
	"pcdn-server/service"
	"time"

	"github.com/robfig/cron/v3"
)

func RunTasks() {
	c := cron.New()

	// 每个小时检查一下限速规则
	c.AddFunc("3 * * * *", func() {
		service.TcService.SyncAllTrifficLimitToDevice()
	})

	c.Start()
}

func taskCycle() {
	var lastSyncTime int64 // 上次同步数据库机器人的时间
	time.Sleep(2 * time.Second)

	for {
		if time.Now().Unix()-lastSyncTime > 60 {
			// syncTaskFromDb() // 每分钟从数据同步Bot状态
			lastSyncTime = time.Now().Unix()
		}

		// var ev *protos.BotTaskEvent
		// select {
		// case ev = <-channels.TaskCannel:
		// 	common.Logger.Debug("task/BotTaskEvent: ", zap.Any("ev", ev))
		// case <-time.After(time.Minute):
		// 	// nothing
		// }

		// common.Logger.Debug("task.ev ", zap.Any("ev: ", ev))
		// if ev != nil {
		// 	if ev.Method == protos.BotTaskStop {
		// 		// stopOneBotTask(fmt.Sprintf("%d-%d", ev.Uid, ev.BotID))
		// 	}
		// }

	}

}
