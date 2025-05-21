package service

import (
	"context"
	"fmt"

	"pcdn-server/common"
	"pcdn-server/protos"

	"github.com/bytedance/sonic"
	passportprotos "github.com/liuhengloveyou/passport/protos"
	"go.uber.org/zap"
)

type martService struct {
}

func (s *martService) LoadLatestKLine(martDomain, symbol string) ([]protos.KLineModel, error) {

	if len(symbol) == 0 {
		return nil, common.ErrParam
	}

	// 系统是否存在Domain
	if !IsMartDomainLegal(martDomain) {
		logger.Error("LoadLatestKLine Domain ERR:", zap.Any("martDomain", martDomain), zap.Any("symbol", symbol))
		return nil, common.ErrMartParamDomain
	}

	b, err := common.RDB.Get(context.Background(), fmt.Sprintf("latestKLine/%s/%s", martDomain, symbol)).Bytes()
	if err != nil {
		return nil, err
	}

	rst := make([]protos.KLineModel, 0)
	err = sonic.Unmarshal(b, &rst)

	return rst, err
}

func (s *martService) LoadDepth(martDomain, symbol string) (*protos.DepthModel, error) {
	if len(symbol) == 0 {
		return nil, common.ErrParam
	}

	// 系统是否存在Domain
	if !IsMartDomainLegal(martDomain) {
		logger.Error("LoadDepth Domain ERR:", zap.Any("martDomain", martDomain), zap.Any("symbol", symbol))
		return nil, common.ErrMartParamDomain
	}

	b, err := common.RDB.Get(context.Background(), fmt.Sprintf("depth/%s/%s", martDomain, symbol)).Bytes()
	if err != nil {
		return nil, err
	}

	var rst protos.DepthModel
	err = sonic.Unmarshal(b, &rst)

	return &rst, err
}

func (s *martService) LoadSpotWallet(uid uint64, martDomain string) ([]protos.SpotAssetsModel, error) {

	if len(martDomain) == 0 {
		return nil, common.ErrParam
	}

	b, err := common.RDB.Get(context.Background(), fmt.Sprintf("spotAccountWallet/%d/%s", uid, martDomain)).Bytes()
	if err != nil {
		return nil, err
	}

	rst := make([]protos.SpotAssetsModel, 0)
	err = sonic.Unmarshal(b, &rst)

	return rst, err
}

// 下单
func (s *martService) NewOrder(sessionUser *passportprotos.User, prop *protos.OrderProp) error {
	if sessionUser == nil {
		return common.ErrSession
	}
	if len(prop.MartDomain) == 0 || len(prop.Symbol) == 0 {
		return common.ErrParam
	}

	// 查APIKEY
	martParam, err := MartParamService.Select(sessionUser.UID, 0, prop.MartDomain)
	common.Logger.Info("martService.NewOrder: ",
		zap.Uint64("uid", sessionUser.UID),
		zap.Any("prop", prop),
		zap.Any("param", martParam),
		zap.Error(err))
	if err != nil || martParam == nil {
		return common.ErrParam
	}

	// mart := marts.GetMartByName(prop.MartDomain)
	// common.Logger.Info("martService.NewOrder: ",
	// 	zap.Uint64("uid", sessionUser.UID),
	// 	zap.Uint64("uid", sessionUser.UID),
	// 	zap.Any("prop", prop),
	// 	zap.Any("param", martParam),
	// 	zap.Any("mart", mart))
	// if mart == nil {
	// 	return common.ErrService
	// }
	// err = mart.Init(martParam)
	// if err != nil {
	// 	common.Logger.Error("martService.NewOrder: ",
	// 		zap.Uint64("uid", sessionUser.UID),
	// 		zap.Uint64("uid", sessionUser.UID),
	// 		zap.Any("prop", prop),
	// 		zap.Any("param", martParam),
	// 		zap.Error(err))
	// 	return err
	// }

	// _, err = mart.NewOrder(&protos.OrderProp{
	// 	Symbol:        prop.Symbol,
	// 	Side:          protos.ORDER_SIDE_BUY,
	// 	Type:          protos.ORDER_TYPE_LIMIT,
	// 	ClientOrderId: fmt.Sprintf("%v", time.Now().UnixMilli()),
	// 	Quantity:      prop.Quantity,
	// 	Price:         prop.Price,
	// 	Notional:      prop.Notional,
	// })
	// common.Logger.Info("martService.NewOrder: ",
	// 	zap.Uint64("uid", sessionUser.UID),
	// 	zap.Uint64("uid", sessionUser.UID),
	// 	zap.Any("prop", prop),
	// 	zap.Any("param", martParam),
	// 	zap.Error(err))
	// if err != nil {
	// 	return common.ErrService
	// }

	return nil
}
