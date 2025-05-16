package protos

type KLineModel struct {
	// 时间戳，毫秒级别，必要字段
	Timestamp int64 `json:"timestamp"`

	// 开盘时间
	OpenTimestamp int64 `json:"openTimestamp"`

	// 收盘时间
	CloseTimestamp int64 `json:"closeTimestamp"`

	// 开盘价，必要字段
	Open string `json:"open"`

	// 收盘价，必要字段
	Close string `json:"close"`

	// 最高价，必要字段
	High string `json:"high"`

	// 最低价，必要字段
	Low string `json:"low"`

	// 成交量，非必须字段
	Volume string `json:"volume"`

	// 成交额，非必须字段，如果需要展示技术指标'EMV'和'AVP'，则需要为该字段填充数据。
	Turnover string `json:"turnover"`
}

type SpotAssetsModel struct {
	Currency  string `json:"currency"`
	Name      string `json:"name"`
	Available string `json:"available"`
	Frozen    string `json:"frozen"`
	Total     string `json:"total"`
	UpdateAt  string `json:"updateAt"`
}

type DepthModel struct {
	Timestamp int64      `json:"ts"`     // Create time(Timestamp in milliseconds)
	Asks      [][]string `json:"asks"`   // 卖盘 [价位, 挂单量] Order book on sell side
	Bids      [][]string `json:"bids"`   // 买盘 [价位, 挂单量] Order book on buy side
	Amount    string     `json:"amount"` // Total number of current price depth
	Price     string     `json:"price"`  // The price at current depth
	Symbol    string     `json:"symbol"` // Trading pair
}

type BestBidOfferStruct struct {
	Symbol          string `json:"symbol"`
	QuoteTime       int64  `json:"quoteTime"`
	HighestBidPrice string `json:"highestBidPrice"` // 买盘1价
	HighestBidQty   string `json:"highestBidQty"`   // 买盘1量
	LowestAskPrice  string `json:"lowestAskPrice"`  // 卖盘1价
	LowestAskQty    string `json:"lowestAskQty"`    // 卖盘1量
}

type PriceModel struct {
	Timestamp int64  `json:"ts"`     // time(Timestamp in milliseconds)
	Symbol    string `json:"symbol"` // Trading pair
	Price     string `json:"price"`  // The price at current depth
}

type TickerModel struct {
	Mart string `json:"mart,omitempty"`

	// Currency pair
	Symbol string `json:"symbol,omitempty"`
	// Last trading price
	Price string `json:"last,omitempty"`
	// Recent lowest ask
	LowestAsk string `json:"lowestAsk,omitempty"`
	// Recent highest bid
	HighestBid string `json:"highestBid,omitempty"`
	//
	SyncAtMilli int64 `json:"syncAt,omitempty"`
}

type TickerTableModel struct {
	// Currency pair
	Symbol string `json:"symbol"`

	Binance           string `json:"binance"`
	BinanceLowestAsk  string `json:"binanceLowestAsk"`
	BinanceHighestBid string `json:"binanceHighestBid"`
	BinanceSyncAt     int64  `json:"binanceSyncAt"`

	Okx           string `json:"okx"`
	OkxLowestAsk  string `json:"okxLowestAsk"`
	OkxHighestBid string `json:"okxHighestBid"`
	OkxSyncAt     int64  `json:"okxSyncAt"`

	Gate           string `json:"gate"`
	GateLowestAsk  string `json:"gateLowestAsk"`
	GateHighestBid string `json:"gateHighestBid"`
	GateSyncAt     int64  `json:"gateSyncAt"`

	Mexc           string `json:"mexc"`
	MexcLowestAsk  string `json:"mexcLowestAsk"`
	MexcHighestBid string `json:"mexcHighestBid"`
	MexcSyncAt     int64  `json:"mexcSyncAt"`

	// BitmartLast       string `json:"bitmartLast"`
	// BitmartLowestAsk  string `json:"bitmartLowestAsk"`
	// BitmartHighestBid string `json:"bitmartHighestBid"`
	// BitmartSyncAt     int64  `json:"bitmartSyncAt"`

	// BiconomyLast       string `json:"biconomyLast"`
	// BiconomyLowestAsk  string `json:"biconomyLowestAsk"`
	// BiconomyHighestBid string `json:"biconomyHighestBid"`
	// BiconomySyncAt     int64  `json:"biconomySyncAt"`

	// XtLast       string `json:"xtLast"`
	// XtLowestAsk  string `json:"xtLowestAsk"`
	// XtHighestBid string `json:"xtHighestBid"`
	// XtSyncAt     int64  `json:"xtSyncAt"`

	Spread      float64 `json:"spread"`
	SpreadRatio float64 `json:"spreadRatio"`
}

func NewTickerTableModel() *TickerTableModel {
	return &TickerTableModel{
		Symbol: "",

		Binance:           "0",
		BinanceLowestAsk:  "0",
		BinanceHighestBid: "0",
		BinanceSyncAt:     0,

		Okx:           "0",
		OkxLowestAsk:  "0",
		OkxHighestBid: "0",
		OkxSyncAt:     0,

		Gate:           "0",
		GateLowestAsk:  "0",
		GateHighestBid: "0",
		GateSyncAt:     0,

		Mexc:           "0",
		MexcLowestAsk:  "0",
		MexcHighestBid: "0",
		MexcSyncAt:     0,

		// Bitmart:       "0",
		// BitmartLowestAsk:  "0",
		// BitmartHighestBid: "0",
		// BitmartSyncAt:     0,

		// Biconomy:       "0",
		// BiconomyLowestAsk:  "0",
		// BiconomyHighestBid: "0",
		// BiconomySyncAt:     0,

		// Xt:       "0",
		// XtLowestAsk:  "0",
		// XtHighestBid: "0",
		// XtSyncAt:     0,

		Spread: 0.0,
	}
}

// func (p *TickerTableModel) CalcSpread() {
// 	// max, _ := strconv.ParseFloat(p.Binance, 64)
// 	// min, _ := strconv.ParseFloat(p.Binance, 64)

// 	lasts := []float64{}
// 	for _, last := range []string{p.Binance, p.Okx, p.GateLast, p.MexcLast} {
// 		lastFloat, _ := strconv.ParseFloat(last, 64)
// 		lasts = append(lasts, lastFloat)
// 	}

// 	if lasts[2] > lasts[3] {
// 		p.Spread = lasts[2] - lasts[3]
// 		if lasts[3] > 0 {
// 			p.SpreadRatio = p.Spread / lasts[3]
// 		}
// 	} else {
// 		p.Spread = lasts[3] - lasts[2]
// 		if lasts[2] > 0 {
// 			p.SpreadRatio = p.Spread / lasts[2]
// 		}
// 	}

// 	// for _, last := range lasts {
// 	// 	if last > maxLast {
// 	// 		maxLast = last
// 	// 	}
// 	// 	if last < minLast {
// 	// 		minLast = last
// 	// 	}
// 	// }

// 	// if minLast <= 0 || maxLast <= 0 {
// 	// 	return 0
// 	// }

// 	// return maxLast - minLast
// }
