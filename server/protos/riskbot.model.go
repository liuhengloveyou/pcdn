package protos

import (
	"database/sql/driver"

	"github.com/bytedance/sonic"
)

// 余额监控机器人
type BalanceMonitorBotProp struct {
	Mail1 string `json:"mail1"`
	Mail2 string `json:"mail2"`
	Mail3 string `json:"mail3"`

	RemindInterval        int64  `json:"remindInterval"`
	BalanceAlertThreshold string `json:"balanceAlertThreshold"`
	TokenAlertThreshold   string `json:"tokenAlertThreshold"`

	// 启动模式：
	LaunchMode BotLaunchModeType `json:"launchMode"` // 立即启动 / 定时启动

	// 机器人描述：
	Description string `json:"description"`
}

func (m *BalanceMonitorBotProp) Scan(src interface{}) error {
	if src == nil {
		return nil
	}

	b, _ := src.([]byte)
	return sonic.Unmarshal(b, m)
}
func (m BalanceMonitorBotProp) Value() (driver.Value, error) {
	return sonic.Marshal(m)
}

// 价格监控机器人
type PriceMonitorBotProp struct {
	Mail1 string `json:"mail1"`

	RemindInterval int64  `json:"remindInterval"`
	MinPrice       string `json:"minPrice"`
	MaxPrice       string `json:"maxPrice"`

	// 启动模式：
	LaunchMode BotLaunchModeType `json:"launchMode"` // 立即启动 / 定时启动

	// 机器人描述：
	Description string `json:"description"`
}

func (m *PriceMonitorBotProp) Scan(src interface{}) error {
	if src == nil {
		return nil
	}

	b, _ := src.([]byte)
	return sonic.Unmarshal(b, m)
}
func (m PriceMonitorBotProp) Value() (driver.Value, error) {
	return sonic.Marshal(m)
}

// 余额定时报告
type BalanceScheduleReportBotProp struct {
	Mail1          string `json:"mail1"`
	Mail2          string `json:"mail2"`
	Mail3          string `json:"mail3"`
	RemindInterval int64  `json:"remindInterval"`

	// 启动模式：
	LaunchMode BotLaunchModeType `json:"launchMode"` // 立即启动 / 定时启动

	// 机器人描述：
	Description string `json:"description"`
}

func (m *BalanceScheduleReportBotProp) Scan(src interface{}) error {
	if src == nil {
		return nil
	}

	b, _ := src.([]byte)
	return sonic.Unmarshal(b, m)
}
func (m BalanceScheduleReportBotProp) Value() (driver.Value, error) {
	return sonic.Marshal(m)
}

// 紧急制动
type EmergencyBrakeBotProp struct {
	Mail1 string `json:"mail1"`
	Mail2 string `json:"mail2"`
	Mail3 string `json:"mail3"`

	BalanceAlertThreshold string `json:"balanceAlertThreshold"`
	TokenAlertThreshold   string `json:"tokenAlertThreshold"`

	// 启动模式：
	LaunchMode BotLaunchModeType `json:"launchMode"` // 立即启动 / 定时启动

	// 机器人描述：
	Description string `json:"description"`
}

func (m *EmergencyBrakeBotProp) Scan(src interface{}) error {
	if src == nil {
		return nil
	}

	b, _ := src.([]byte)
	return sonic.Unmarshal(b, m)
}
func (m EmergencyBrakeBotProp) Value() (driver.Value, error) {
	return sonic.Marshal(m)
}

// TG Scheduled Report
type TGScheduledReportRobotProp struct {
	ChatId   string `json:"chatId"`
	UserName string `json:"userName"`

	RemindInterval int64 `json:"remindInterval"`

	// 启动模式：
	LaunchMode BotLaunchModeType `json:"launchMode"` // 立即启动 / 定时启动
}

func (m *TGScheduledReportRobotProp) Scan(src interface{}) error {
	if src == nil {
		return nil
	}

	b, _ := src.([]byte)
	return sonic.Unmarshal(b, m)
}
func (m TGScheduledReportRobotProp) Value() (driver.Value, error) {
	return sonic.Marshal(m)
}

// TG Alert
type TGAlertRobotProp struct {
	ChatId         string `json:"chatId"`
	UserName       string `json:"userName"`
	AvailableMoney string `json:"availableMoney"`
	AvailableToken string `json:"availableToken"`
	TotalMoney     string `json:"totalMoney"`
	TotalToken     string `json:"totalToken"`
	PriceRangeMin  string `json:"priceRangeMin"`
	PriceRangeMax  string `json:"priceRangeMax"`
	RemindInterval int64  `json:"remindInterval"`

	// 启动模式：
	LaunchMode BotLaunchModeType `json:"launchMode"` // 立即启动 / 定时启动
}

func (m *TGAlertRobotProp) Scan(src interface{}) error {
	if src == nil {
		return nil
	}

	b, _ := src.([]byte)
	return sonic.Unmarshal(b, m)
}
func (m TGAlertRobotProp) Value() (driver.Value, error) {
	return sonic.Marshal(m)
}

type RiskBotModel struct {
	Model

	BotType    BotType `json:"botType" gorm:"column:bot_type;type:INT;not null;"`
	MartDomain string  `json:"martDomain" gorm:"column:mart_domain;type:varchar(45);not null;"`
	Symbol     string  `json:"symbol" gorm:"column:symbol;type:VARCHAR(45);not null;"`

	IsRunning BotRunType `json:"isRunning" gorm:"column:is_running;type:INT;not null;"`
	StartTime int64      `json:"startTime" gorm:"column:start_time;type:INT;"`

	BalanceMonitorBot        *BalanceMonitorBotProp        `json:"balanceMonitorBot,omitempty" gorm:"column:balance_monitor_bot;type:JSON"`
	PriceMonitorBot          *PriceMonitorBotProp          `json:"priceMonitorBot,omitempty" gorm:"column:price_monitor_bot;type:JSON"`
	BalanceScheduleReportBot *BalanceScheduleReportBotProp `json:"balanceScheduleReportBot,omitempty" gorm:"column:balance_schedule_report;type:JSON"`
	EmergencyBrakeBot        *EmergencyBrakeBotProp        `json:"emergencyBrakeBot,omitempty" gorm:"column:emergency_brake_bot;type:JSON"`
	TGScheduledReportBot     *TGScheduledReportRobotProp   `json:"tgScheduledReportBot,omitempty" gorm:"column:tg_scheduled_report_bot;type:JSON"`
	TGAlertBot               *TGAlertRobotProp             `json:"tgAlertBot,omitempty" gorm:"column:tg_alert_bot;type:JSON"`
}

func (RiskBotModel) TableName() string {
	return "risk_bots"
}
