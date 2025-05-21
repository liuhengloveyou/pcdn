package channels

import "pcdn-server/protos"

var (
	// 买VIP会员成功支付后， 发消息过来
	PayTaskCannel = make(chan *protos.VIPOrderStruct, 100)

	// bot状态变更时发消息过来
	TaskCannel = make(chan *protos.BotTaskEvent, 100)

	BinanceCrawlerChan = make(chan *protos.CrawlerEvent, 100)
	OkxCrawlerChan     = make(chan *protos.CrawlerEvent, 100)
	GateIoCrawlerChan  = make(chan *protos.CrawlerEvent, 100)
)

func PostPayTask(event *protos.VIPOrderStruct) {
	PayTaskCannel <- event
}

func PostBotTask(event *protos.BotTaskEvent) {
	TaskCannel <- event
}

func PostToGateIoCrawler(event *protos.CrawlerEvent) {
	GateIoCrawlerChan <- event
}

func PostToBinanceCrawler(event *protos.CrawlerEvent) {
	BinanceCrawlerChan <- event
}

func PostToOkxCrawler(event *protos.CrawlerEvent) {
	OkxCrawlerChan <- event
}
