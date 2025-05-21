package models

type OrderSideType string

const (
	ORDER_SIDE_BUY  OrderSideType = "BUY"
	ORDER_SIDE_SELL OrderSideType = "SELL"
)

type OrderTypeType string

const (
	ORDER_TYPE_LIMIT  OrderTypeType = "LIMIT"  // 限价单
	ORDER_TYPE_MARKET OrderTypeType = "MARKET" // 市价单
// LIMIT_MAKER 限价只挂单
// IMMEDIATE_OR_CANCEL IOC单 (无法立即成交的部分就撤销,订单在失效前会尽量多的成交。)
// FILL_OR_KILL FOK单 (无法全部立即成交就撤销,如果无法全部成交,订单会失效。)
)

type OrderProp struct {
	MartDomain    string        `json:"mart"`
	Symbol        string        `json:"symbol"`
	Side          OrderSideType `json:"side"`
	Type          OrderTypeType `json:"type"`
	Quantity      string        `json:"quantity"`
	Price         string        `json:"price,omitempty"`
	Notional      string        `json:"notional,omitempty"`
	MartOrderId   string        `json:"mart_order_id,omitempty"`
	ClientOrderId string        `json:"client_order_id,omitempty"`
	Timestamp     string        `json:"timestamp,omitempty"`

	Delayed int64  `json:"delayed"`
	ErrMsg  string `json:"errMsg"`
}

var OrderSides = []OrderSideType{
	ORDER_SIDE_BUY,
	ORDER_SIDE_SELL,
}

type OrderModel struct {
	Model

	MartDomain    string `json:"martDomain" gorm:"column:mart_domain;type:varchar(45);not null;"`
	MartOrderId   string `json:"martOrderId" gorm:"column:mart_order_id;type:varchar(128);"`
	ClientOrderId string `json:"client_order_id,omitempty" gorm:"column:client_order_id;type:varchar(45);"`
	Symbol        string `json:"symbol" gorm:"column:symbol;type:varchar(45);not null;"`
	Side          string `json:"side" gorm:"column:side;type:varchar(45);not null;"`
	Type          string `json:"type" gorm:"column:type;type:varchar(45);not null;"`
	Quantity      string `json:"quantity,omitempty" gorm:"column:quantity;type:varchar(45);"`
	Price         string `json:"price,omitempty" gorm:"column:price;type:varchar(45);"`
	Notional      string `json:"notional,omitempty" gorm:"column:motional;type:varchar(100);"`
	DelayedMilli  int64  `json:"delayed,omitempty" gorm:"column:delayed_milli;"`
	ErrMsg        string `json:"errMsg,omitempty" gorm:"column:err_msg;"`
}

func (OrderModel) TableName() string {
	return "orders"
}
