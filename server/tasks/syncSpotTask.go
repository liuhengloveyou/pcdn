package tasks

import (
	"arbitrage/common"
	"arbitrage/marts"
	"arbitrage/protos"
	"context"
	"fmt"
	"time"

	"github.com/bytedance/sonic"
	"go.uber.org/zap"
)

func syncSpotTask() {
	for {
		// 每分钟同步一次
		time.Sleep(time.Minute)

		// 每个网站的一个交易对只需要查一次Kline
		// syncedKline := make(map[string]bool)

		apiKeys, err := findAllApiKey()
		common.Logger.Debug("syncSpotTask: ", zap.Any("keys", len(apiKeys)), zap.Error(err))

		for i := 0; i < len(apiKeys); i++ {
			mart := marts.GetMartByName(apiKeys[i].MartDomain)
			if mart == nil {
				common.Logger.Error("syncSpotTask mart nil: ", zap.Any("key", apiKeys[i]))
				continue
			}

			if err := mart.Init(&protos.MartParamModel{
				AccessKey: apiKeys[i].AccessKey,
				SecretKey: apiKeys[i].SecretKey,
				Memo:      apiKeys[i].Memo,
			}); err != nil {
				common.Logger.Error("syncSpotTask mart init ERR: ", zap.Any("apikey", apiKeys[i]))
				continue
			}
			common.Logger.Info("syncSpotTask GetDomain: ", zap.Any("apikey", apiKeys[i]), zap.String("domain", mart.GetDomain()))

			// 查账号余额
			balanceData, err := mart.GetSpotAccountAssets()
			common.Logger.Info("syncSpotTask.balance: ",
				zap.Any("apikey", apiKeys[i]),
				zap.String("domain", mart.GetDomain()),
				zap.Any("data", len(balanceData)),
				zap.Error(err))
			if err != nil {
				common.Logger.Error("syncSpotTask: ",
					zap.Any("apikey", apiKeys[i]),
					zap.String("domain", mart.GetDomain()),
					zap.Any("data", balanceData),
					zap.Error(err))
				continue
			}

			// 存到緩存
			if len(balanceData) > 0 {
				balanceData[0].UpdateAt = time.Now().Format("2006-01-02 15:04")
				b, _ := sonic.Marshal(balanceData)
				err = common.RDB.Set(context.Background(), fmt.Sprintf("spotAccountWallet/%d/%s", apiKeys[i].UserId, apiKeys[i].MartDomain), b, time.Hour*100).Err()
				if err != nil {
					common.Logger.Error("syncSpotTask ERR: ",
						zap.Any("apikey", apiKeys[i]),
						zap.String("domain", mart.GetDomain()),
						zap.Any("data", balanceData),
						zap.Error(err))
				}
			}

			// // 每個交易對查一次就行了
			// klineCacheKey := fmt.Sprintf("latestKLine/%s/%s", apiKeys[i].MartDomain, apiKeys[i].Symbol)
			// if !syncedKline[klineCacheKey] {
			// 	// 查最新價格
			// 	syncLatestPrice(mart, apiKeys[i].MartDomain, apiKeys[i].Symbol)
			// 	// 查Kline
			// 	syncLatestKLine(mart, apiKeys[i].MartDomain, apiKeys[i].Symbol)
			// 	// 查depth
			// 	syncDepth(mart, apiKeys[i].MartDomain, apiKeys[i].Symbol)
			// 	syncedKline[klineCacheKey] = true
			// }
		}
	}

}

func syncLatestPrice(mart marts.Mart, martDomain, symbol string) {

	data, err := mart.GetPrice(symbol)
	if err != nil {
		common.Logger.Error("syncSpotTask.syncLatestPrice ERR: ", zap.String("symbol", symbol), zap.Any("mart", mart), zap.Error(err))
		return
	}

	// 存到緩存
	b, _ := sonic.Marshal(data)
	err = common.RDB.Set(context.Background(), fmt.Sprintf("latestPrice/%s/%s", martDomain, symbol), b, time.Hour*100).Err()
	common.Logger.Error("syncSpotTask.syncLatestPrice: ", zap.String("symbol", symbol), zap.Any("data", len(b)), zap.Error(err))
}

func syncLatestKLine(mart marts.Mart, martDomain, symbol string) {

	data, err := mart.LatestKLine(symbol)
	if err != nil {
		common.Logger.Error("syncSpotTask.latestKLine ERR: ", zap.String("symbol", symbol), zap.Any("mart", mart), zap.Error(err))
		return
	}

	// 存到緩存
	b, _ := sonic.Marshal(data)
	err = common.RDB.Set(context.Background(), fmt.Sprintf("latestKLine/%s/%s", martDomain, symbol), b, time.Hour*100).Err()
	common.Logger.Info("syncSpotTask.latestKLine: ", zap.String("mart", martDomain), zap.String("symbol", symbol), zap.Any("data", len(data)), zap.Error(err))

}

func syncDepth(mart marts.Mart, martDomain, symbol string) {
	data, err := mart.GetDepth(symbol, 10)
	if err != nil {
		common.Logger.Error("syncSpotTask.syncDepth ERR: ", zap.String("symbol", symbol), zap.Any("mart", mart), zap.Error(err))
		return
	}

	// 存到緩存
	b, _ := sonic.Marshal(data)
	err = common.RDB.Set(context.Background(), fmt.Sprintf("depth/%s/%s", martDomain, symbol), b, time.Hour*100).Err()
	common.Logger.Error("syncSpotTask.syncDepth: ", zap.String("symbol", symbol), zap.Any("data", len(b)), zap.Error(err))

}

// 从数据库里查询所有APIKEY
func findAllApiKey() (rr []protos.MartParamModel, err error) {
	err = common.OrmCli.Table("martparams").Find(&rr).Error

	return
}
