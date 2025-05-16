package protos

const BINANCE = "binance"

const OKX = "okx"

// https://www.gate.io/docs/developers/apiv4/zh_CN
const GATE = "gate"

// https://www.bitmart.com/open-api-guide/en-US
const BITMART = "bitmart"

// https://www.mexc.com/zh-CN/mexc-api
const MEXC = "mexc"

// https://doc.xt.com/
const XT = "xt"

// https://github.com/BiconomyOfficial/apidocs
const BICONOMY = "biconomy"

// https://www.htx.com/zh-cn/opend/newApiPages/?id=7ec44bc5-7773-11ed-9966-0242ac110003
const HTX = "huobi"

// https://bybit-exchange.github.io/docs/zh-TW/v5/websocket/public/orderbook
const BYBIT = "bybit"

var MartsInSystem = []string{
	BINANCE,
	OKX,
	GATE,
	BITMART,
	MEXC,
	XT,
	BICONOMY,
	HTX,
	BYBIT,
}

type MartParamModel struct {
	Model

	MartDomain string `json:"domain" gorm:"column:domain;index:idx_domain;type:VARCHAR(45);"`
	AccessKey  string `json:"accessKey" gorm:"column:access_key;type:VARCHAR(128);"`
	SecretKey  string `json:"secretKey,omitempty" gorm:"column:secret_key;type:VARCHAR(128);"`
	Memo       string `json:"memo,omitempty" gorm:"column:memo;type:VARCHAR(128);"`
	Passphrase string `json:"passphrase,omitempty" gorm:"column:passphrase;type:VARCHAR(128);"`

	Active int64 `json:"active" gorm:"column:active;"`
}

func (MartParamModel) TableName() string {
	return "martparams"
}
