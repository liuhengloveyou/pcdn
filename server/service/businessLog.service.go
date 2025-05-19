package service

import (
	"pcdn-server/common"
	"pcdn-server/protos"
	"pcdn-server/repos"
	"time"

	"go.uber.org/zap"
)

type businessLogService struct {
}

func (s *businessLogService) Add(m *protos.BusinessLog) (id uint64, err error) {
	m.CreateTime = time.Now().UnixMilli()
	m.UpdateTime = m.CreateTime

	return repos.BusinessLogRepo.Add(m)
}

func (s *businessLogService) Find(uid uint64, martDomain string, page int) (rr []protos.BusinessLog, err error) {
	if uid <= 0 {
		return nil, nil
	}

	rr, err = repos.BusinessLogRepo.Find(uid, martDomain, page)
	if err != nil {
		logger.Error("Find ERR: ", zap.Error(err))
		return nil, common.ErrService
	}

	return
}
