// https://api.bscscan.com/api?module=account&action=tokentx&contractaddress=0x766e09665a8128D9F7Fd9d90BCD3d9cdc50b067F&address=0x32Acf1A0DE3AE6a96582d5c5AEB171ABD6e609F3&page=1&offset=5&startblock=0&endblock=999999999&sort=asc&apikey=YMNQFTQ9R89ZEQWF23TD3M72ZNCY1D3KGW
package tasks

import (
	"arbitrage/bsc"
	"arbitrage/channels"
	"arbitrage/common"
	"arbitrage/protos"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/nanmu42/etherscan-api"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var bscScanClient *etherscan.Client

func PayTaskCycle() {
	time.Sleep(2 * time.Second)

	bscScanClient = etherscan.NewCustomized(etherscan.Customization{
		Timeout: 15 * time.Second,
		Key:     "YMNQFTQ9R89ZEQWF23TD3M72ZNCY1D3KGW",
		BaseURL: "https://api.bscscan.com/api?",
		Verbose: false,
	})

	for {
		var ev *protos.VIPOrderStruct
		select {
		case ev = <-channels.PayTaskCannel:
			common.Logger.Debug("PayTaskCycle: ", zap.Any("ev", ev))
		case <-time.After(time.Minute):
			// nothing
		}

		// 简单处理接口限速问题5/s
		time.Sleep(time.Millisecond * 300)

		DealPayTask()
	}

}

func DealPayTask() {
	orders, err := findAllPayingOrder()
	if err != nil {
		common.Logger.Debug("PayTaskCycle: ", zap.Error(err))
		return
	}

	if len(orders) > 0 {
		// checkPay(orders)
		if err := checkPayByBscscan(orders); err != nil {
			common.Logger.Error("checkPayByBscscan ERR: ", zap.Error(err))
		}
	}

	// 处理超时的
	dealTimeOutPayingOrder()
}

// https://api.bscscan.com/api?module=account&action=tokentx&contractaddress=0x766e09665a8128D9F7Fd9d90BCD3d9cdc50b067F&address=0x32Acf1A0DE3AE6a96582d5c5AEB171ABD6e609F3&page=1&offset=5&startblock=0&endblock=999999999&sort=asc&apikey=YMNQFTQ9R89ZEQWF23TD3M72ZNCY1D3KGW
func checkPayByBscscan(orders []protos.VIPOrderStruct) error {
	// check ERC20 transactions from/to a specified address
	contractAddress := "0x766e09665a8128D9F7Fd9d90BCD3d9cdc50b067F" // MOLI 合约地址
	address := "0x32Acf1A0DE3AE6a96582d5c5AEB171ABD6e609F3"         // 收款地址
	startblock := 0
	endblock := 999999999
	transfers, err := bscScanClient.ERC20Transfers(&contractAddress, &address, &startblock, &endblock, 1, 100, true)
	if err != nil {
		return err
	}

	payok := 0
	for _, tran := range transfers {
		common.Logger.Sugar().Infoln("tran:", tran.ContractAddress, time.Since(tran.TimeStamp.Time()).Minutes(), tran)
		if time.Since(tran.TimeStamp.Time()) > 10*time.Minute {
			continue
		}
		if strings.Compare(strings.ToUpper(contractAddress), strings.ToUpper(tran.ContractAddress)) != 0 {
			continue
		}

		for _, order := range orders {
			vipPrice := strings.ReplaceAll(order.VipPrice, ".", "")

			common.Logger.Sugar().Infoln("From: ", order.WalletAddr, tran.From, strings.Compare(order.WalletAddr, tran.From))
			common.Logger.Sugar().Infoln("To: ", order.Receiver, tran.To, strings.Compare(order.Receiver, tran.To))
			common.Logger.Sugar().Infoln("Value: ", vipPrice, tran.Value, tran.Value.Int().String())

			if order.Status == protos.ORDERSTAT_PAYING &&
				strings.Compare(strings.ToUpper(order.Receiver), strings.ToUpper(tran.To)) == 0 &&
				strings.Compare(strings.ToUpper(order.WalletAddr), strings.ToUpper(tran.From)) == 0 &&
				strings.HasPrefix(tran.Value.Int().String(), vipPrice) {

				// 查到支付记录了
				order.TranHash = tran.Hash
				PayingOrderOK(&order)
			}
			if order.Status == protos.ORDERSTAT_PAY_OK {
				payok += 1
			}

		}
		if payok >= len(orders) {
			break
		}
	}

	return nil
}

func checkPay(orders []protos.VIPOrderStruct) {
	client, err := bsc.NewBscClient()
	if err != nil {
		common.Logger.Error("NewBscClient ERR: ", zap.Error(err))
		return
	}

	latestBlockNumber, _ := client.LatestBlockNumber()
	if latestBlockNumber <= 0 {
		return
	}
	// latestBlockNumber = 45990262 // TODO

	// 最大查i个块
	for i := 0; i < 10; i++ {
		blockNumber := latestBlockNumber - uint64(i)
		block := client.BlockByNumber(blockNumber)
		common.Logger.Debug("BlockByNumber: ", zap.Any("i", i), zap.Any("blockNumber", blockNumber), zap.Any("blockTime", block.Time()))
		if block == nil {
			common.Logger.Error("BlockByNumber ERR")
			return
		}
		if (uint64(time.Now().Unix()) - block.Time()) > 300 {
			return // 看查最近5分钟内数据
		}

		payok := 0
		trans := client.TokenTransfer(blockNumber)
		for _, tran := range trans {
			for _, order := range orders {
				// common.Logger.Debug()
				fmt.Println("From: ", order.WalletAddr, tran.From.String(), strings.Compare(order.WalletAddr, tran.From.String()))
				fmt.Println("To: ", order.Receiver, tran.To.String(), strings.Compare(order.Receiver, tran.To.String()))
				fmt.Println("Value: ", order.VipPrice, tran.Value)
				vipPrice := strings.ReplaceAll(order.VipPrice, ".", "")

				if order.Status == protos.ORDERSTAT_PAYING &&
					strings.Compare(order.Receiver, tran.To.String()) == 0 &&
					strings.Compare(order.WalletAddr, tran.From.String()) == 0 &&
					strings.HasPrefix(tran.Value.String(), vipPrice) {
					// 查到支付记录了
					PayingOrderOK(&order)
				}
				if order.Status == protos.ORDERSTAT_PAY_OK {
					payok += 1
				}

			}
			if payok >= len(orders) {
				break
			}
		}
	}

}

// 从数据库里查询所有正在支付的订单
func findAllPayingOrder() (rr []protos.VIPOrderStruct, err error) {
	tx := common.OrmCli.Table("vip_orders")
	tx = tx.Where("status = ?", protos.ORDERSTAT_PAYING)
	err = tx.Order("id desc").Find(&rr).Error

	return
}

// 超时的正在支付的订单
func dealTimeOutPayingOrder() {
	tx := common.OrmCli.Table("vip_orders")
	tx = tx.Where("status = ? AND update_time < ?", protos.ORDERSTAT_PAYING, time.Now().UnixMilli()-600*1000)
	err := tx.UpdateColumns(map[string]interface{}{"update_time": time.Now().UnixMilli(), "status": protos.ORDERSTAT_PAYERROR}).Error
	if err != nil {
		common.Logger.Error("dealTimeOutPayingOrder ERR: ", zap.Error(err))
	}
}

// 正在支付的订单成功
func PayingOrderOK(order *protos.VIPOrderStruct) {

	vipPrice, _ := strconv.ParseFloat(order.VipPrice, 64)
	if vipPrice <= 0 {
		return
	}

	// 查会员数据
	var member protos.VIPMemberStruct
	err := common.OrmCli.Where("uid = ?", order.UserId).Take(&member).Error
	if err != nil {
		common.Logger.Error("PayingOrderOK ERR: ", zap.Error(err))
	}

	common.OrmCli.Transaction(func(tx *gorm.DB) error {
		err := tx.Table("vip_orders").Where("id = ?", order.Id).UpdateColumns(map[string]interface{}{"update_time": time.Now().UnixMilli(), "status": protos.ORDERSTAT_PAY_OK, "tran_hash": order.TranHash}).Error
		if err != nil {
			common.Logger.Error("dealTimeOutPayingOrder ERR: ", zap.Error(err))
			return err
		}

		if member.Id > 0 {
			// 更新会员时间
			member.UpdateTime = time.Now().UnixMilli()
			endTime := time.UnixMilli(member.EndTime)
			if endTime.Before(time.Now()) {
				endTime = time.Now()
			}

			if vipPrice >= 1200 {
				member.EndTime = endTime.AddDate(1, 0, 0).UnixMilli()
			} else if vipPrice >= 150 {
				member.EndTime = endTime.AddDate(0, 1, 0).UnixMilli()
			} else if vipPrice >= 10 {
				member.EndTime = endTime.AddDate(0, 0, 1).UnixMilli()
			}
			if err := tx.Updates(&member).Error; err != nil {
				return err
			}
		} else {
			// 新建会员记录
			member.CreateTime = time.Now().UnixMilli()
			member.UpdateTime = member.CreateTime
			member.UserId = order.UserId
			if vipPrice >= 1200 {
				member.EndTime = time.Now().AddDate(1, 0, 0).UnixMilli()
			} else if vipPrice >= 150 {
				member.EndTime = time.Now().AddDate(0, 1, 0).UnixMilli()
			} else if vipPrice >= 10 {
				member.EndTime = time.Now().AddDate(0, 0, 1).UnixMilli()
			}

			if err := tx.Create(&member).Error; err != nil {
				return err
			}
		}
		return nil
	})

	return
}
