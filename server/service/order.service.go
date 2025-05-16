package service

import (
	"arbitrage/common"
	"arbitrage/protos"
	"arbitrage/repos"

	"go.uber.org/zap"
)

type orderService struct {
}

func (s *orderService) Select(uid int64) (m *protos.OrderModel, err error) {
	if uid <= 0 {
		return nil, common.ErrParam
	}

	// m, err = repos.BotRepo.Select(uid)
	// if err != nil {
	// 	logger.Error("select ERR: ", zap.Error(err))
	// 	return nil, common.ErrService
	// }

	return nil, nil
}

func (s *orderService) Find(uid uint64, martDomain, symbol string, page int) (rr []protos.OrderModel, err error) {
	if page < 1 {
		page = 1
	}

	// if len(martDomain) == 0 || len(symbol) == 0 {
	// 	return nil, common.ErrParam
	// }

	rr, err = repos.OrderRepo.Find(uid, martDomain, symbol, page)
	if err != nil {
		logger.Error("Find ERR: ", zap.Error(err))
		return nil, common.ErrService
	}

	return
}
