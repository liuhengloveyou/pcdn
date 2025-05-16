package tasks

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"arbitrage/common"
	"arbitrage/email"
	"arbitrage/protos"
	"arbitrage/repos"

	"github.com/bytedance/sonic"
	"go.uber.org/zap"
)

func riskBotCycle() {
	for {
		// 每分钟同步一次
		time.Sleep(time.Second * 55)

		riskBots, err := loadAllRunningRiskBot()
		if err != nil {
			common.Logger.Error("dealRiskBotLogic ", zap.Error(err))
		}
		common.Logger.Debug("riskBotCycle ", zap.Any("bots", len(riskBots)), zap.Error(err))

		RiskBotLogic(riskBots)
	}
}

func RiskBotLogic(bots []protos.RiskBotModel) {

	for i := 0; i < len(bots); i++ {
		common.Logger.Debug("riskBotCycle ", zap.Any("bot", bots[i]))

		if bots[i].BotType == protos.BotTypeBalanceMonitor {
			if bots[i].BalanceMonitorBot != nil {
				dealBalanceMonitorBot(&bots[i])
			}
		} else if bots[i].BotType == protos.BotTypePriceMonitor {
			if bots[i].PriceMonitorBot != nil {
				dealPriceMonitorBot(&bots[i])
			}
		} else if bots[i].BotType == protos.BotTypeBalanceScheduledReport {
			if bots[i].BalanceScheduleReportBot != nil {
				dealBalanceScheduleReportBot(&bots[i])
			}
		} else if bots[i].BotType == protos.BotTypeEmergencyBrake {
			if bots[i].EmergencyBrakeBot != nil {
				dealEmergencyBrakeBot(&bots[i])
			}
		}

	}
}

func dealBalanceMonitorBot(bot *protos.RiskBotModel) {
	now := time.Now().Unix()
	common.Logger.Info("dealBalanceMonitorBot ",
		zap.Int64("now", now),
		zap.Int64("updateTime", bot.UpdateTime/1000),
		zap.Int64("remindInterval", bot.BalanceMonitorBot.RemindInterval*60),
		zap.Any("bot", bot.BalanceMonitorBot))

	if now-bot.UpdateTime/1000 < bot.BalanceMonitorBot.RemindInterval*60 {
		return
	}

	symbols := strings.Split(bot.Symbol, "_")
	if len(symbols) != 2 {
		common.Logger.Error("dealBalanceMonitorBot symbol ERR",
			zap.Any("symbols", symbols),
			zap.Any("bot", bot.BalanceMonitorBot))
		return
	}

	tokenStr := strings.ToLower(symbols[0])
	balanceStr := strings.ToLower(symbols[1])
	var tokenVal float64
	var balanceVal float64
	var tokenIdx int
	var balanceIdx int

	// 查餘額
	balance, err := common.RDB.Get(context.Background(), fmt.Sprintf("spotAccountWallet/%d/%s", bot.UserId, bot.MartDomain)).Bytes()
	if err != nil {
		common.Logger.Error("dealBalanceMonitorBot redis ERR", zap.Error(err), zap.Any("balance", balance), zap.Any("bot", bot))
		return
	}

	var balanceJson []protos.SpotAssetsModel
	err = sonic.Unmarshal(balance, &balanceJson)
	if err != nil {
		common.Logger.Error("dealBalanceMonitorBot redis ERR", zap.Error(err), zap.Any("balance", balance), zap.Any("bot", bot))
		return
	}

	for i := 0; i < len(balanceJson); i++ {
		if strings.Compare(tokenStr, strings.ToLower(balanceJson[i].Currency)) == 0 {
			tokenVal, _ = strconv.ParseFloat(balanceJson[i].Available, 64)
			tokenIdx = i
			continue
		}
		if strings.Compare(balanceStr, strings.ToLower(balanceJson[i].Currency)) == 0 {
			balanceVal, _ = strconv.ParseFloat(balanceJson[i].Available, 64)
			balanceIdx = i
			continue
		}
	}
	if tokenIdx < 0 || tokenIdx >= len(balanceJson) || balanceIdx < 0 || balanceIdx >= len(balanceJson) {
		common.Logger.Error("dealBalanceMonitorBot redis ERR", zap.Error(err), zap.Any("balance", balance), zap.Any("bot", bot))
		return
	}

	tokenAlertThreshold, _ := strconv.ParseFloat(bot.BalanceMonitorBot.TokenAlertThreshold, 64)
	balanceAlertThreshold, _ := strconv.ParseFloat(bot.BalanceMonitorBot.BalanceAlertThreshold, 64)

	if tokenVal <= tokenAlertThreshold || balanceVal <= balanceAlertThreshold {
		if err := email.SendRiskManBalanceMail(bot.BalanceMonitorBot.Mail1, bot.MartDomain, bot.Symbol, &balanceJson[tokenIdx], &balanceJson[balanceIdx]); err != nil {
			common.Logger.Error("dealBalanceMonitorBot ", zap.Any("bot", bot.BalanceMonitorBot), zap.Error(err))
			return
		}
		if len(bot.BalanceMonitorBot.Mail2) > 0 {
			if err := email.SendRiskManBalanceMail(bot.BalanceMonitorBot.Mail2, bot.MartDomain, bot.Symbol, &balanceJson[tokenIdx], &balanceJson[balanceIdx]); err != nil {
				common.Logger.Error("dealBalanceMonitorBot ", zap.Any("bot", bot.BalanceMonitorBot), zap.Error(err))
				return
			}
		}

		if len(bot.BalanceMonitorBot.Mail3) > 0 {
			if err := email.SendRiskManBalanceMail(bot.BalanceMonitorBot.Mail3, bot.MartDomain, bot.Symbol, &balanceJson[tokenIdx], &balanceJson[balanceIdx]); err != nil {
				common.Logger.Error("dealBalanceMonitorBot ", zap.Any("bot", bot.BalanceMonitorBot), zap.Error(err))
				return
			}
		}

		// 更新数据库时间
		repos.RiskBotRepo.Update(bot.UserId, bot.Id, bot.BotType, bot.MartDomain, bot.Symbol)
	}
}

func dealPriceMonitorBot(bot *protos.RiskBotModel) {
	now := time.Now().Unix()
	common.Logger.Info("dealPriceMonitorBot ",
		zap.Int64("now", now),
		zap.Int64("updateTime", bot.UpdateTime/1000),
		zap.Int64("remindInterval", bot.PriceMonitorBot.RemindInterval*60),
		zap.Any("bot", bot.PriceMonitorBot))

	if now-bot.UpdateTime/1000 < bot.PriceMonitorBot.RemindInterval*60 {
		return
	}

	symbols := strings.Split(bot.Symbol, "_")
	if len(symbols) != 2 {
		common.Logger.Error("dealPriceMonitorBot symbol ERR",
			zap.Any("symbols", symbols),
			zap.Any("bot", bot.PriceMonitorBot))
		return
	}

	// 查價格
	pirce, err := common.RDB.Get(context.Background(), fmt.Sprintf("latestPrice/%s/%s", bot.MartDomain, bot.Symbol)).Bytes()
	if err != nil {
		common.Logger.Error("dealPriceMonitorBot redis ERR", zap.Error(err), zap.Any("pirce", pirce), zap.Any("bot", bot))
		return
	}

	var priceJson protos.PriceModel
	err = sonic.Unmarshal(pirce, &priceJson)
	if err != nil {
		common.Logger.Error("dealPriceMonitorBot redis ERR", zap.Error(err), zap.Any("pirce", pirce), zap.Any("bot", bot))
		return
	}

	priceVal, _ := strconv.ParseFloat(priceJson.Price, 64)
	minPrice, _ := strconv.ParseFloat(bot.PriceMonitorBot.MinPrice, 64)
	maxPrice, _ := strconv.ParseFloat(bot.PriceMonitorBot.MaxPrice, 64)
	common.Logger.Info("dealPriceMonitorBot ",
		zap.Int64("now", now),
		zap.Float64("priceVal", priceVal),
		zap.Float64("minPrice", minPrice),
		zap.Float64("maxPrice", maxPrice),
		zap.Int64("updateTime", bot.UpdateTime/1000),
		zap.Int64("remindInterval", bot.PriceMonitorBot.RemindInterval*60),
		zap.Any("bot", bot.PriceMonitorBot))

	if priceVal <= maxPrice && priceVal >= minPrice {
		if err := email.SendRiskManPriceMail(bot.PriceMonitorBot.Mail1, bot.MartDomain, bot.Symbol, &priceJson); err != nil {
			common.Logger.Error("PriceMonitorBot ", zap.Any("bot", bot.PriceMonitorBot), zap.Error(err))
			return
		}

		// 更新数据库时间
		repos.RiskBotRepo.Update(bot.UserId, bot.Id, bot.BotType, bot.MartDomain, bot.Symbol)
	}
}

func dealBalanceScheduleReportBot(bot *protos.RiskBotModel) {
	now := time.Now().Unix()
	common.Logger.Info("dealBalanceScheduleReportBot ",
		zap.Int64("now", now),
		zap.Int64("updateTime", bot.UpdateTime/1000),
		zap.Int64("remindInterval", bot.BalanceScheduleReportBot.RemindInterval),
		zap.Any("bot", bot.BalanceScheduleReportBot))

	if now-bot.UpdateTime/1000 < bot.BalanceScheduleReportBot.RemindInterval*60 {
		return
	}

	symbols := strings.Split(bot.Symbol, "_")
	if len(symbols) != 2 {
		common.Logger.Error("BalanceScheduleReportBot symbol ERR",
			zap.Any("symbols", symbols),
			zap.Any("bot", bot))
		return
	}

	tokenStr := strings.ToLower(symbols[0])
	balanceStr := strings.ToLower(symbols[1])
	var tokenIdx int
	var balanceIdx int

	// 查餘額
	balance, _ := common.RDB.Get(context.Background(), fmt.Sprintf("spotAccountWallet/%d/%s", bot.UserId, bot.MartDomain)).Bytes()
	var balanceJson []protos.SpotAssetsModel
	err := sonic.Unmarshal(balance, &balanceJson)
	if err != nil {
		common.Logger.Error("BalanceScheduleReportBot redis ERR", zap.Error(err), zap.Any("balance", balance), zap.Any("bot", bot))
		return
	}

	for i := 0; i < len(balanceJson); i++ {
		if strings.Compare(tokenStr, strings.ToLower(balanceJson[i].Currency)) == 0 {
			tokenIdx = i
			continue
		}
		if strings.Compare(balanceStr, strings.ToLower(balanceJson[i].Currency)) == 0 {
			balanceIdx = i
			continue
		}
	}
	if tokenIdx < 0 || tokenIdx >= len(balanceJson) || balanceIdx < 0 || balanceIdx >= len(balanceJson) {
		common.Logger.Error("BalanceScheduleReportBot redis ERR", zap.Error(err), zap.Any("balance", balance), zap.Any("bot", bot))
		return
	}

	if err := email.SendRiskManBalanceMail(bot.BalanceScheduleReportBot.Mail1, bot.MartDomain, bot.Symbol, &balanceJson[tokenIdx], &balanceJson[balanceIdx]); err != nil {
		common.Logger.Error("BalanceScheduleReportBot ", zap.Any("bot", bot), zap.Error(err))
		return
	}
	if len(bot.BalanceScheduleReportBot.Mail2) > 0 {
		if err := email.SendRiskManBalanceMail(bot.BalanceScheduleReportBot.Mail2, bot.MartDomain, bot.Symbol, &balanceJson[tokenIdx], &balanceJson[balanceIdx]); err != nil {
			common.Logger.Error("BalanceScheduleReportBot ", zap.Any("bot", bot), zap.Error(err))
		}
	}
	if len(bot.BalanceScheduleReportBot.Mail3) > 0 {
		if err := email.SendRiskManBalanceMail(bot.BalanceScheduleReportBot.Mail3, bot.MartDomain, bot.Symbol, &balanceJson[tokenIdx], &balanceJson[balanceIdx]); err != nil {
			common.Logger.Error("BalanceScheduleReportBot ", zap.Any("bot", bot), zap.Error(err))
		}
	}

	// 更新数据库时间
	repos.RiskBotRepo.Update(bot.UserId, bot.Id, bot.BotType, bot.MartDomain, bot.Symbol)
}

// Emergency Brake
func dealEmergencyBrakeBot(bot *protos.RiskBotModel) {
	now := time.Now().Unix()
	common.Logger.Info("dealEmergencyBrakeBot ",
		zap.Int64("now", now),
		zap.Int64("updateTime", bot.UpdateTime/1000),
		zap.Any("bot", bot))

	symbols := strings.Split(bot.Symbol, "_")
	if len(symbols) != 2 {
		common.Logger.Error("dealEmergencyBrakeBot symbol ERR",
			zap.Any("symbols", symbols),
			zap.Any("bot", bot.EmergencyBrakeBot))
		return
	}

	tokenStr := strings.ToLower(symbols[0])
	balanceStr := strings.ToLower(symbols[1])
	var tokenVal float64
	var balanceVal float64
	var tokenIdx int
	var balanceIdx int

	// 查餘額
	balance, _ := common.RDB.Get(context.Background(), fmt.Sprintf("spotAccountWallet/%d/%s", bot.UserId, bot.MartDomain)).Bytes()
	var balanceJson []protos.SpotAssetsModel
	err := sonic.Unmarshal(balance, &balanceJson)
	common.Logger.Error("dealEmergencyBrakeBot balance:", zap.Error(err), zap.Any("balance", balance), zap.Any("bot", bot))
	if err != nil {
		return
	}

	for i := 0; i < len(balanceJson); i++ {
		if strings.Compare(tokenStr, strings.ToLower(balanceJson[i].Currency)) == 0 {
			tokenVal, _ = strconv.ParseFloat(balanceJson[i].Available, 64)
			tokenIdx = i
			continue
		}
		if strings.Compare(balanceStr, strings.ToLower(balanceJson[i].Currency)) == 0 {
			balanceVal, _ = strconv.ParseFloat(balanceJson[i].Available, 64)
			balanceIdx = i
			continue
		}
	}
	if tokenIdx < 0 || tokenIdx >= len(balanceJson) || balanceIdx < 0 || balanceIdx >= len(balanceJson) {
		common.Logger.Error("dealEmergencyBrakeBot redis ERR", zap.Error(err), zap.Any("balance", balance), zap.Any("bot", bot))
		return
	}

	tokenAlertThreshold, _ := strconv.ParseFloat(bot.EmergencyBrakeBot.TokenAlertThreshold, 64)
	balanceAlertThreshold, _ := strconv.ParseFloat(bot.EmergencyBrakeBot.BalanceAlertThreshold, 64)

	if tokenVal <= tokenAlertThreshold || balanceVal <= balanceAlertThreshold {
		if err := email.SendRiskManBalanceMail(bot.EmergencyBrakeBot.Mail1, bot.MartDomain, bot.Symbol, &balanceJson[tokenIdx], &balanceJson[balanceIdx]); err != nil {
			common.Logger.Error("EmergencyBrakeBot mail2 ", zap.Any("bot", bot.EmergencyBrakeBot), zap.Error(err))
			return
		}
		if len(bot.EmergencyBrakeBot.Mail2) > 0 {
			if err := email.SendRiskManBalanceMail(bot.EmergencyBrakeBot.Mail2, bot.MartDomain, bot.Symbol, &balanceJson[tokenIdx], &balanceJson[balanceIdx]); err != nil {
				common.Logger.Error("EmergencyBrakeBot mail2", zap.Any("bot", bot.EmergencyBrakeBot), zap.Error(err))
				return
			}
		}

		if len(bot.EmergencyBrakeBot.Mail3) > 0 {
			if err := email.SendRiskManBalanceMail(bot.EmergencyBrakeBot.Mail3, bot.MartDomain, bot.Symbol, &balanceJson[tokenIdx], &balanceJson[balanceIdx]); err != nil {
				common.Logger.Error("EmergencyBrakeBot mail3", zap.Any("bot", bot.EmergencyBrakeBot), zap.Error(err))
				return
			}
		}

		// 更新数据库时间
		repos.RiskBotRepo.Update(bot.UserId, bot.Id, bot.BotType, bot.MartDomain, bot.Symbol)
	}
}

// 从数据库里查询所有运行状态的风控机器人
func loadAllRunningRiskBot() (rr []protos.RiskBotModel, err error) {
	tx := common.OrmCli.Table("risk_bots")
	tx = tx.Where("is_running = ?", protos.BotIsRunning)
	err = tx.Find(&rr).Order("uid, id").Error

	return
}
