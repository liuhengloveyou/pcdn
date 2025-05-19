package service

import (
	"pcdn-server/common"
	"pcdn-server/protos"
	"pcdn-server/repos"

	"go.uber.org/zap"
)

type martParamService struct{}

func (s *martParamService) Set(req *protos.MartParamModel) (uint64, error) {
	if req.UserId <= 0 {
		return 0, common.ErrParam
	}

	if len(req.MartDomain) == 0 {
		logger.Error("SetMartParam Domain ERR:", zap.Any("req", req))
		return 0, common.ErrMartParamDomain
	}
	if len(req.AccessKey) == 0 || len(req.SecretKey) == 0 {
		logger.Error("SetMartParam key ERR:", zap.Any("req", req))
		return 0, common.ErrMartParamKey
	}

	// 系统是否存在Domain
	if !IsMartDomainLegal(req.MartDomain) {
		logger.Error("SetMartParam Domain ERR:", zap.Any("req", req))
		return 0, common.ErrMartParamDomain
	}

	id, err := repos.MartParamRepo.Set(req)
	logger.Debug("SetMartParam:", zap.Any("req", req), zap.Uint64("id", id), zap.Error(err))
	if err != nil {
		logger.Error("SetMartParam.Set ERR: ", zap.Error(err))
		return id, common.ErrService
	}

	return id, nil
}

func (s *martParamService) Find(uid uint64) (rr []protos.MartParamModel, err error) {
	if uid <= 0 {
		return nil, nil
	}

	rr, err = repos.MartParamRepo.Find(uid, "")
	if err != nil {
		logger.Error("Find ERR: ", zap.Error(err))
		return nil, common.ErrService
	}

	return
}

func (s *martParamService) Select(uid, id uint64, martDomain string) (m *protos.MartParamModel, err error) {
	if uid <= 0 {
		return nil, common.ErrParam
	}

	m, err = repos.MartParamRepo.Select(uid, id, martDomain)
	if err != nil {
		logger.Error("select ERR: ", zap.Any("uid", uid), zap.Any("id", id), zap.String("mart", martDomain), zap.Error(err))
		return nil, common.ErrService
	}

	return m, nil
}

func (s *martParamService) Active(id, uid uint64) (int64, error) {
	if uid <= 0 || id <= 0 {
		return 0, common.ErrParam
	}

	r, err := repos.MartParamRepo.UpdateActive(id, uid)
	if err != nil {
		logger.Error("Active ERR: ", zap.Error(err))
		return 0, common.ErrService
	}

	return r, nil
}

func (s *martParamService) Delete(id, uid uint64) (int64, error) {
	if uid <= 0 || id <= 0 {
		return 0, common.ErrParam
	}

	r, err := repos.MartParamRepo.Delete(id, uid)
	common.Logger.Warn("martParamService.Delete", zap.Uint64("id", id), zap.Uint64("uid", uid))
	if err != nil {
		logger.Error("Delete ERR: ", zap.Error(err))
		return 0, common.ErrService
	}

	return r, nil
}
