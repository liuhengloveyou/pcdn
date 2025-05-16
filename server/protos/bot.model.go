package protos

import (
	"database/sql/driver"
	"fmt"

	"github.com/bytedance/sonic"
)

type BotTaskMethod int

const (
	BotTaskRun  BotTaskMethod = 1 // 启动
	BotTaskStop BotTaskMethod = 2 // 停止
)

type BotTaskEvent struct {
	Method BotTaskMethod

	Uid   uint64
	BotID uint64

	Payload any
}

type BotRunType int

const (
	BotIsRunNone BotRunType = 0
	BotIsRunning BotRunType = 1 // 正在运行
	BotIsNotRun  BotRunType = 2 // 已经停
	BotIsRunEnd  BotRunType = 3
)

// 机器人下单模式
type BotOrderModeType string

const (
	BotOrderRandom    BotOrderModeType = "random"
	BotOrderBuyFirst  BotOrderModeType = "buyFirst"
	BotOrderSellFirst BotOrderModeType = "sellFirst"
)

// 机器人价格模式
type BotPriceModeType string

const (
	BotPriceNormalMode BotPriceModeType = "normalMode"
	BotPriceSafeMode   BotPriceModeType = "safeMode"
)

// 机器人启动模式
type BotLaunchModeType string

const (
	BotImmediateLaunch BotLaunchModeType = "immediateStart" // 立即启动
	BotScheduledLaunch BotLaunchModeType = "scheduledStart" // 定时启动
)

// 机器人数量分布模式
type BotQuantModeType string

const (
	BotContinuousDistribution BotQuantModeType = "continuousDist" // 连续分布
	BotRandomDistribution     BotQuantModeType = "randomDist"     // 随机分布
)

type BotType int

const (
	// Risk Manager
	BotTypeBalanceMonitor         BotType = 300
	BotTypePriceMonitor           BotType = 301
	BotTypeBalanceScheduledReport BotType = 302
	BotTypeEmergencyBrake         BotType = 303
	BotTypeTGScheduledReport      BotType = 304
	BotTypeTGAlert                BotType = 305

	// 套利
	BotTypeObserveAndReport   BotType = 500
	BotTypeArbirageDepth1     BotType = 501
	BotTypeAutoArbirageDepth1 BotType = 502
	BotTypeTimerReport        BotType = 503
)

// 深度1对敲观察报告
type ObserveAndReportBotModel struct {
	TGChatId string `json:"tgChatId"`

	MartA   string `json:"martA"`
	SymbolA string `json:"symbolA"`

	MartB   string `json:"martB"`
	SymbolB string `json:"symbolB"`

	// 最小价差比例
	// (高价 - 低价) / 低价
	MinSpreadRatio string `json:"minSpreadRatio"`
	// 最小挂单量
	MinVolume string `json:"minVolume"`

	// 每单最小挂单量
	MinOrderSize int64 `json:"minOrderSize"`
	// 每单最大挂单量
	MaxOrderSize int64 `json:"maxOrderSize"`

	// 提醒间隔秒数
	TimeInterval int `json:"timeInterval"`

	// 启动模式：
	LaunchMode BotLaunchModeType `json:"launchMode"`
}

func (m *ObserveAndReportBotModel) Scan(src interface{}) error {
	if src == nil {
		return nil
	}

	b, _ := src.([]byte)
	return sonic.Unmarshal(b, m)
}
func (m ObserveAndReportBotModel) Value() (driver.Value, error) {
	return sonic.Marshal(m)
}

// 深度1对敲观察报告
type TimerReportBotModel struct {
	TGChatId string `json:"tgChatId"`

	MartA   string `json:"martA"`
	SymbolA string `json:"symbolA"`

	MartB   string `json:"martB"`
	SymbolB string `json:"symbolB"`

	// 提醒间隔秒数
	TimeInterval int `json:"timeInterval"`

	// 启动模式：
	LaunchMode BotLaunchModeType `json:"launchMode"`
}

func (m *TimerReportBotModel) Scan(src interface{}) error {
	if src == nil {
		return nil
	}

	b, _ := src.([]byte)
	return sonic.Unmarshal(b, m)
}
func (m TimerReportBotModel) Value() (driver.Value, error) {
	return sonic.Marshal(m)
}

// 深度1对敲
type ArbirageDepth1BotModel struct {
	MartA   string `json:"martA"`
	SymbolA string `json:"symbolA"`

	MartB   string `json:"martB"`
	SymbolB string `json:"symbolB"`

	// 最小价差比例
	// (高价 - 低价) / 低价
	MinSpreadRatio string `json:"minSpreadRatio"`
	// 最小挂单量
	MinVolume string `json:"minVolume"`
	// 每单最小挂单量
	MinOrderSize int64 `json:"minOrderSize"`
	// 每单最大挂单量
	MaxOrderSize int64 `json:"maxOrderSize"`

	TGChatId string `json:"tgChatId"`

	// 启动模式：
	LaunchMode BotLaunchModeType `json:"launchMode"`
}

func (m *ArbirageDepth1BotModel) Scan(src interface{}) error {
	if src == nil {
		return nil
	}

	b, _ := src.([]byte)
	return sonic.Unmarshal(b, m)
}
func (m ArbirageDepth1BotModel) Value() (driver.Value, error) {
	return sonic.Marshal(m)
}

// 深度1对敲,自动转账
type AutoArbirageDepth1BotModel struct {
	MartA   string `json:"martA"`
	SymbolA string `json:"symbolA"`

	MartB   string `json:"martB"`
	SymbolB string `json:"symbolB"`

	// 最小价差比例
	// (高价 - 低价) / 低价
	MinSpreadRatio string `json:"minSpreadRatio"`

	// 最小挂单量
	MinVolume string `json:"minVolume"`

	TGChatId string `json:"tgChatId"`

	// 启动模式：
	LaunchMode BotLaunchModeType `json:"launchMode"`
}

func (m *AutoArbirageDepth1BotModel) Scan(src interface{}) error {
	if src == nil {
		return nil
	}

	b, _ := src.([]byte)
	return sonic.Unmarshal(b, m)
}
func (m AutoArbirageDepth1BotModel) Value() (driver.Value, error) {
	return sonic.Marshal(m)
}

type BotModel struct {
	Model

	BotType   BotType    `json:"botType" gorm:"column:bot_type;type:INT;not null;"`
	IsRunning BotRunType `json:"isRunning" gorm:"column:is_running;type:INT;not null;"`
	StartTime int64      `json:"startTime" gorm:"column:start_time;type:INT;"`

	TimerReportBot        *TimerReportBotModel        `json:"timerReportBot,omitempty" gorm:"column:timer_report_bot;type:JSON"`
	ObserveAndReportBot   *ObserveAndReportBotModel   `json:"observeAndReportBot,omitempty" gorm:"column:observe_and_report;type:JSON"`
	ArbirageDepth1Bot     *ArbirageDepth1BotModel     `json:"arbirageDepth1Bot,omitempty" gorm:"column:arbirage_depth1;type:JSON"`
	AutoArbirageDepth1Bot *AutoArbirageDepth1BotModel `json:"autoArbirageDepth1Bot,omitempty" gorm:"column:auto_arbirage_depth1;type:JSON"`
}

func (BotModel) TableName() string {
	return "bots"
}

func (p *BotModel) Key() string {
	return fmt.Sprintf("%d-%d", p.UserId, p.Id)
}
