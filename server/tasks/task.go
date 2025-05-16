package tasks

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"arbitrage/bots"
	"arbitrage/channels"
	"arbitrage/common"
	martsws "arbitrage/marts-ws"
	"arbitrage/protos"
	"arbitrage/service"

	"go.uber.org/zap"
)

var (
	// 系统中正在运行的bot
	// uid-botType-site 每个账号-每个类型-每个网站
	runningBots sync.Map
)

func RunTasks() {
	go PayTaskCycle()
	go taskCycle()
	// go syncSpotTask()
	// go riskBotCycle()
}

func taskCycle() {
	var lastSyncTime int64 // 上次同步数据库机器人的时间
	time.Sleep(2 * time.Second)

	for {
		if time.Now().Unix()-lastSyncTime > 60 {
			syncTaskFromDb() // 每分钟从数据同步Bot状态
			lastSyncTime = time.Now().Unix()
		}

		var ev *protos.BotTaskEvent
		select {
		case ev = <-channels.TaskCannel:
			common.Logger.Debug("task/BotTaskEvent: ", zap.Any("ev", ev))
		case <-time.After(time.Minute):
			// nothing
		}

		common.Logger.Debug("task.ev ", zap.Any("ev: ", ev))
		if ev != nil {
			if ev.Method == protos.BotTaskStop {
				stopOneBotTask(fmt.Sprintf("%d-%d", ev.Uid, ev.BotID))
			}
		}

	}

}

func syncTaskFromDb() {
	// 查询运行状态的bot
	allRunStatBots, err := findAllRunningBot()
	common.Logger.Debug("task/findAllRunningBot: ", zap.Int("bots", len(allRunStatBots)))
	if err != nil {
		common.Logger.Error("task/findAllRunningBot ERR: ", zap.Error(err))
		return
	}

	// 停掉没有RUN的 和更新过的
	runningBots.Range(func(key, value any) bool {
		has := false
		bot := value.(bots.Bot)
		for i := 0; i < len(allRunStatBots); i++ {
			key1 := allRunStatBots[i].Key()
			if strings.Compare(key.(string), key1) == 0 && allRunStatBots[i].UpdateTime == bot.Prop().UpdateTime {
				has = true
			}
		}

		if !has {
			stopOneBotTask(key.(string))
		}

		return true
	})

	for i := 0; i < len(allRunStatBots); i++ {
		key := allRunStatBots[i].Key()
		bot, ok := runningBots.Load(key)
		if ok && bot != nil {
			continue // 已经在运行
		}

		// 需要新建一个
		runBot := bots.GetBotByType(allRunStatBots[i].BotType)
		if runBot == nil {
			common.Logger.Warn("task botType:", zap.Int64("botType", int64(allRunStatBots[i].BotType)))
			continue
		}

		var initErr error
		switch allRunStatBots[i].BotType {
		case protos.BotTypeArbirageDepth1:
			initErr = runBot.Init(allRunStatBots[i])
		case protos.BotTypeObserveAndReport:
			initErr = runBot.Init(allRunStatBots[i])
		case protos.BotTypeTimerReport:
			initErr = runBot.Init(allRunStatBots[i])
		default:
			common.Logger.Warn("task no bot type:", zap.Any("bot", allRunStatBots[i]))
			continue
		}
		common.Logger.Info("task.init: ", zap.Any("botType", runBot.Type()), zap.Any("name", runBot.Name()), zap.Any("param", allRunStatBots[i]), zap.Error(initErr))
		if initErr != nil {
			log := &protos.BusinessLog{
				BusinessType: protos.BUSINESS_TYPE_ERR,
				Payload:      fmt.Sprintf("BOT init ERR: [%v]; bot: [%v]", err, allRunStatBots[i]),
			}
			log.UserId = allRunStatBots[i].UserId
			log.TenantId = allRunStatBots[i].TenantId
			service.BusinessLogService.Add(log)

			common.Logger.Warn("task init ERR:", zap.Error(err), zap.Any("bot", allRunStatBots[i]))
			continue
		}

		go runBot.Run()
		runningBots.Store(key, runBot)
	}

	// 请求最优价
	runningBots.Range(func(key, value any) bool {
		bot := value.(bots.Bot)

		switch bot.Type() {
		case protos.BotTypeArbirageDepth1:
			if bot.Prop() != nil && bot.Prop().ArbirageDepth1Bot != nil {
				martwsA := martsws.GetMartByName(bot.Prop().ArbirageDepth1Bot.MartA)
				if martwsA != nil {
					martwsA.AddSymbol(bot.Prop().ArbirageDepth1Bot.SymbolA)
				}
				martwsB := martsws.GetMartByName(bot.Prop().ArbirageDepth1Bot.MartB)
				if martwsB != nil {
					martwsB.AddSymbol(bot.Prop().ArbirageDepth1Bot.SymbolB)
				}
			}
		case protos.BotTypeObserveAndReport:
			if bot.Prop() != nil && bot.Prop().ObserveAndReportBot != nil {
				martwsA := martsws.GetMartByName(bot.Prop().ObserveAndReportBot.MartA)
				if martwsA != nil {
					martwsA.AddSymbol(bot.Prop().ObserveAndReportBot.SymbolA)
				}
				martwsB := martsws.GetMartByName(bot.Prop().ObserveAndReportBot.MartB)
				if martwsB != nil {
					martwsB.AddSymbol(bot.Prop().ObserveAndReportBot.SymbolB)
				}
			}
		case protos.BotTypeTimerReport:
			// fmt.Println("syncTaskFromDb: ", bot.Name(), bot.Prop().TimerReportBot)
			if bot.Prop() != nil && bot.Prop().TimerReportBot != nil {
				martwsA := martsws.GetMartByName(bot.Prop().TimerReportBot.MartA)
				if martwsA != nil {
					martwsA.AddSymbol(bot.Prop().TimerReportBot.SymbolA)
				}
				martwsB := martsws.GetMartByName(bot.Prop().TimerReportBot.MartB)
				if martwsB != nil {
					martwsB.AddSymbol(bot.Prop().TimerReportBot.SymbolB)
				}
			}
		}
		return true
	})
}

func stopOneBotTask(key string) {
	bot, _ := runningBots.Load(key)
	if bot != nil {
		common.Logger.Warn("stopOneBotTask: ", zap.String("key", key), zap.Any("bot", bot.(bots.Bot).Name()))

		if _, ok := bot.(bots.Bot); ok {
			bot.(bots.Bot).Stop()
			runningBots.Delete(key)
		}
	}
}

// 从数据库里查询所有运行状态的机器人
func findAllRunningBot() (rr []protos.BotModel, err error) {
	tx := common.OrmCli.Table("bots")
	tx = tx.Where("is_running = ?", protos.BotIsRunning)
	err = tx.Find(&rr).Error

	return
}
