package protos

type CrawlerMethod int

const (
	CrawlerRun  CrawlerMethod = 1 // 启动
	CrawlerStop CrawlerMethod = 2 // 停止
)

type CrawlerEvent struct {
	Method BotTaskMethod
}
