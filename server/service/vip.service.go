package service

import (
	"arbitrage/channels"
	"arbitrage/common"
	"arbitrage/protos"
	"fmt"
	"math/rand"
	"strings"
	"time"

	passportprotos "github.com/liuhengloveyou/passport/protos"
	"go.uber.org/zap"
)

type vipOrderService struct {
}

func (s *vipOrderService) MemberInfo(sess *passportprotos.User) (m *protos.VIPMemberStruct, err error) {

	err = common.OrmCli.Where("uid = ?", sess.UID).Take(&m).Error
	if err != nil {
		common.Logger.Error("MemberInfo db ERR: ", zap.Error(err))
		return nil, common.ErrService
	}

	return
}

func (s *vipOrderService) Save(sess *passportprotos.User, m *protos.VIPOrderStruct) (err error) {
	// 如果有正在进行中的订单
	var rr []protos.VIPOrderStruct
	tx := common.OrmCli.Table("vip_orders")
	tx = tx.Where("uid = ? AND status = ?", sess.UID, protos.ORDERSTAT_PAYING)
	err = tx.Order("id desc").Find(&rr).Error
	if len(rr) > 0 {
		return common.ErrVIPPaying
	}

	m.CreateTime = time.Now().UnixMilli()
	m.UpdateTime = m.CreateTime
	m.UserId = sess.UID
	m.Status = protos.ORDERSTAT_CREATE_OK
	m.OrderId = fmt.Sprintf("vip-%v%v", m.VipPrice, time.Now().UnixMicro())
	m.Receiver = common.ServConfig.Receiver
	m.VipPrice = fmt.Sprintf("%s.%d", m.VipPrice, rand.New(rand.NewSource(time.Now().UnixMicro())).Intn(1000))

	err = common.OrmCli.Save(m).Error

	return
}

func (s *vipOrderService) Payed(sess *passportprotos.User, oid int64, orderId string) (err error) {
	var m protos.VIPOrderStruct

	err = common.OrmCli.Where("id = ?", oid).Take(&m).Error
	if err != nil {
		common.Logger.Error("Payed db ERR: ", zap.Any("oid", oid), zap.String("orderId", orderId), zap.Error(err))
		return common.ErrService
	}
	if m.UserId != sess.UID || strings.Compare(m.OrderId, orderId) != 0 {
		common.Logger.Error("Payed db ERR: ", zap.Any("oid", oid), zap.String("orderId", orderId), zap.Any("m", m))
		return common.ErrService
	}

	now := time.Now()
	// 更新数据库到正在PAYING状态
	err = common.OrmCli.Table("vip_orders").Where("id = ?", oid).UpdateColumns(map[string]interface{}{"update_time": now.UnixMilli(), "status": protos.ORDERSTAT_PAYING}).Error
	if err != nil {
		common.Logger.Error("Payed db ERR: ", zap.Any("oid", oid), zap.String("orderId", orderId), zap.Error(err))
		return common.ErrService
	}

	// 查支付记录
	channels.PostPayTask(&m)

	return
}
