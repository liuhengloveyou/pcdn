package tasks

import (
	"sync"
	"time"
)

var (
	// 系统中正在运行的bot
	// uid-botType-site 每个账号-每个类型-每个网站
	runningBots sync.Map
)

func RunTasks() {
	// go PayTaskCycle()
	// go taskCycle()
	// go syncSpotTask()
	// go riskBotCycle()
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
