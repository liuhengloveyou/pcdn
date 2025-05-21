package repos

import (
	"pcdn-server/common"
	"pcdn-server/protos"
)

type vipRepo struct {
}

func (p *vipRepo) LoadMemberInfo(uid uint64) (m *protos.VIPMemberStruct, err error) {

	m = &protos.VIPMemberStruct{}
	err = common.OrmCli.Where("uid = ?", uid).Take(m).Error

	return
}
