package protos

// 定单状态
type VIPOrderStatType int

const (
	ORDERSTAT_CREATE_OK VIPOrderStatType = 1
	ORDERSTAT_PAYING    VIPOrderStatType = 2
	ORDERSTAT_PAY_OK    VIPOrderStatType = 3
	ORDERSTAT_PAYERROR  VIPOrderStatType = 4
)

type VIPOrderStruct struct {
	Model

	OrderId string `json:"orderId" gorm:"column:order_id;type:varchar(128);"`

	// 商品信息
	VipPrice string `json:"vipPrice" gorm:"column:vip_price;type:varchar(128);"`
	// 钱包地址
	WalletAddr string `json:"walletAddr"  gorm:"column:wallet_addr;type:varchar(128);"`

	// 收款地址
	Receiver string `json:"receiver"  gorm:"column:receiver;type:varchar(128);"`
	// 支付时间
	PayTime uint64 `json:"payTime" gorm:"column:pay_time;"`

	Status   VIPOrderStatType `json:"status" gorm:"column:status;type:int"`
	TranHash string           `json:"tranHash" gorm:"column:tran_hash;type:varchar(128);"`
}

func (VIPOrderStruct) TableName() string {
	return "vip_orders"
}

type VIPMemberStruct struct {
	Model

	EndTime int64 `json:"endTime" db:"end_time" gorm:"column:end_time;"`
}

func (VIPMemberStruct) TableName() string {
	return "vip_members"
}
