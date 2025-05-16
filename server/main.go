package main

import (
	"arbitrage/api"
	"arbitrage/common"
	"arbitrage/protos"
	"arbitrage/repos"
	"arbitrage/tasks"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	BuildTime string
	CommitID  string

	Sig string
)

var (
	showVer = flag.Bool("version", false, "Print version information.")
	initSys = flag.Bool("init", false, "初始化系统环境")
)

// func init() {
// 	os.Setenv("HTTP_PROXY", "http://127.0.0.1:10808")
// 	os.Setenv("HTTPS_PROXY", "http://127.0.0.1:10808")
// 	os.Setenv("NO_PROXY", "localhost,127.0.0.1")
// }

func sigHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM)

	go func() {
		s := <-c
		Sig = "service is suspend ..."
		fmt.Println("Got signal:", s)
	}()
}

func main() {
	flag.Parse()
	sigHandler()

	if *showVer {
		fmt.Printf("%s\t%s\n", BuildTime, CommitID)
		return
	}

	if *initSys {
		InitSystemEnv()
		return
	}

	go tg.InitTG()
	martsws.InitMartsWS()
	crawler.Run()
	tasks.RunTasks()

	if err := api.InitAndRunHttpApi(common.ServConfig.Addr); err != nil {
		panic("HTTPAPI: " + err.Error())
	}
}

func InitSystemEnv() {
	if err := repos.InitModels(common.OrmCli); err != nil {
		panic(err)
	}
}

func testBybit() {
	fmt.Println(">>>testBybit")

	mart := marts.GetMartByName(protos.BYBIT)
	if err := mart.Init(&protos.MartParamModel{
		AccessKey: "uqBDYJtzy01cEbjL2i",
		SecretKey: "8S1OKUZbDs7m3tRjDeou2rSnvVA6Mh4kFne3",
	}); err != nil {
		panic(err)
	}

	data := protos.OrderProp{
		ClientOrderId: fmt.Sprintf("%v", time.Now().UnixMilli()),
		Symbol:        "BTC_USDT",
		Side:          protos.ORDER_SIDE_BUY,
		Type:          protos.ORDER_TYPE_LIMIT,
		Price:         "80000",
		Quantity:      "1",
	}
	oid, e := mart.NewOrder(&data)
	fmt.Println("buy>>>>>>>>>", oid, e)
	oid, e = mart.CancelOrder("BTC_USDT", oid)
	fmt.Println("cal>>>>>>>>>", oid, e)
}

func testHuobi() {
	fmt.Println(">>>testHuobi")

	mart := marts.GetMartByName(protos.HTX)
	if err := mart.Init(&protos.MartParamModel{
		AccessKey: "dab4c45e6f-6fdffe7a-42e3f109-6c66d",
		SecretKey: "1f1a500b-a7d7cef2-e4b46269-faebe",
	}); err != nil {
		panic(err)
	}

	data := protos.OrderProp{
		ClientOrderId: fmt.Sprintf("%v", time.Now().UnixMilli()),
		Symbol:        "BIO_USDT",
		Side:          protos.ORDER_SIDE_BUY,
		Type:          protos.ORDER_TYPE_LIMIT,
		Price:         "0.5",
		Quantity:      "1",
	}
	oid, e := mart.NewOrder(&data)
	fmt.Println("buy>>>>>>>>>", oid, e)
	// oid, e = mart.CancelOrder("BIO_USDT", oid)
	// fmt.Println("cal>>>>>>>>>", oid, e)
}

func testGateio() {
	mart := marts.GetMartByName(protos.GATE)
	if err := mart.Init(&protos.MartParamModel{
		AccessKey: "61ad814e0624cc5d74887c0eb15cc14b",
		SecretKey: "93df35255a5f2ec0368c08815502cab5e21ac6d7ba9c9420ccae3466cd716d72",
	}); err != nil {
		panic(err)
	}

	tickers, err := mart.GetTickers()
	fmt.Println("price>>>", tickers, err)
}

// apikey = "0ec2af95-2cd9-4ee1-bb68-b74dde2ed8c1" secretkey = "686BD71634E40AEF769B3F686A8CAC4B" IP = "" 备注名 = "22" 权限 = "读取/提现/交易"
func testOkx() {
	mart := marts.GetMartByName(protos.OKX)
	if err := mart.Init(&protos.MartParamModel{
		AccessKey:  "0ec2af95-2cd9-4ee1-bb68-b74dde2ed8c1",
		SecretKey:  "686BD71634E40AEF769B3F686A8CAC4B",
		Passphrase: "WWW-222-wwww",
	}); err != nil {
		panic(err)
	}

	resp, err := mart.GetSpotAccountAssets()
	fmt.Printf(">>>%#v %v\n\n", resp, err)

	data := protos.OrderProp{
		ClientOrderId: fmt.Sprintf("%v", time.Now().UnixMilli()),
		Symbol:        "BIO_USDT",
		Side:          protos.ORDER_SIDE_BUY,
		Type:          protos.ORDER_TYPE_LIMIT,
		Price:         "0.5",
		Quantity:      "1",
	}
	oid, e := mart.NewOrder(&data)
	fmt.Println("buy>>>>>>>>>", oid, e)
	oid, e = mart.CancelOrder("BIO_USDT", oid)
	fmt.Println("cal>>>>>>>>>", oid, e)

	// param := &market.GetTickersParam{
	// 	InstType: okxapi.InstTypeSPOT,
	// }
	// req, resp := market.NewGetTickers(param)
	// if err := client.Do(req, resp); err != nil {
	// 	panic(err)
	// }
	// log.Println(req, resp.(*market.GetTickersResponse))
}

func testBinance() {
	mart := marts.GetMartByName(protos.BINANCE)
	if err := mart.Init(&protos.MartParamModel{
		AccessKey: "arWICCgv6uiCMxJC0o5wUqMrUER4nDaNAFAEGn0IqfV731N8KAgV9b8SrjhF5ifV",
		SecretKey: "EOMmXwEq2QFKEatavnX0EYsyXnki3uKPAwl8j3FEVpnHDJXURB7py3louQ1IEWab",
	}); err != nil {
		panic(err)
	}

	resp, err := mart.GetPrice("MEME_USDT")
	fmt.Println(">>>>>>>>>", resp, err)
	// data := protos.OrderProp{
	// 	ClientOrderId: fmt.Sprintf("%v", time.Now().UnixMilli()),
	// 	Symbol:        "MEME_USDT",
	// 	Side:          protos.ORDER_SIDE_BUY,
	// 	Type:          protos.ORDER_TYPE_LIMIT,
	// 	Price:         "0.007",
	// 	Quantity:      "1000",
	// }

	// resp, err := mart.NewOrder(&data)
	// fmt.Println("order>>>", resp, err)
	// resp, err = mart.CancelOrder("MEME_USDT", resp)
	// fmt.Println("cancel order>>>", resp, err)
}

func testMartBiconomy() {
	mart := marts.GetMartByName(protos.BICONOMY)
	if err := mart.Init(&protos.MartParamModel{
		AccessKey: "e85b8238-d1a8-4488-9f6c-188f34898ea3",
		SecretKey: "8d2068a3-a230-4aa2-832d-f78855130986",
		// AccessKey: "29cf7cfd-0ab0-4090-858d-e23cf39765d5",
		// SecretKey: "db953657-daed-42d2-9618-701715a5c411",
	}); err != nil {
		panic(err)
	}

	price, err := mart.GetPrice("MOLI_USDT")
	fmt.Println("price>>>", price, err)

	// assets, err := mart.GetSpotAccountAssets()
	// fmt.Println("assets>>>", assets, err)

	// depth, err := mart.GetDepth("MOLI_USDT", 10)
	// fmt.Println("GetDepth>>>", mart.GetDomain(), depth, err)
	// if depth == nil {
	// 	return
	// }

	// kline, err := mart.LatestKLine("MOLI_USDT")
	// fmt.Println("kline>>>", kline, err)

	// data := protos.OrderProp{
	// 	Symbol:   "MOLI_USDT",
	// 	Side:     protos.ORDER_SIDE_BUY,
	// 	Type:     protos.ORDER_TYPE_LIMIT,
	// 	Price:    "0.100",
	// 	Quantity: "100",
	// }
	// id, err := mart.NewOrder(&data)
	// fmt.Println("NewOrder", id, err)
	// id, err = mart.CancelOrder("MOLI_USDT", id)
	// fmt.Println("CancelOrder", id, err)

	return
}

func testMexc() {
	// mx0vglH4CTbudP0ed9
	// fd0fcfcd82e54bf1acd29c4f1ec70d09
	mexc := marts.GetMartByName(protos.MEXC)
	if err := mexc.Init(&protos.MartParamModel{
		AccessKey: "mx0vglM1ZOW6z0XIHj",
		SecretKey: "9d6e0665353242fca9f69f1ecb4a9185",
	}); err != nil {
		panic(err)
	}

	price, err := mexc.GetBookTicker("MOLI_USDT")
	fmt.Println("GetDomain>>>", mexc.GetDomain(), price, err)

	// depth, err := mexc.GetDepth("MOLI_USDT", 1)
	// fmt.Println("GetDomain>>>", mexc.GetDomain(), depth)
	// if depth == nil {
	// 	return
	// }

	// // get_balance
	// res, err := mexc.GetSpotAccountAssets()
	// fmt.Println("assets:: ", res, err)

	// kline
	// res, err := mexc.LatestKLine("MOLIUSDT")
	// fmt.Println(res, err)

	// data := protos.OrderProp{
	// 	Symbol:   "MOLI_USDT",
	// 	Side:     protos.ORDER_SIDE_BUY,
	// 	Type:     protos.ORDER_TYPE_LIMIT,
	// 	Price:    "0.0100",
	// 	Quantity: "100",
	// }
	// orderId, err := mexc.NewOrder(&data)
	// fmt.Println("orderId: ", orderId, err)

	// time.Sleep(time.Millisecond * 100)
	// orderId, err = mexc.CancelOrder("MOLI_USDT", orderId)
	// fmt.Println("cancle orderId: ", orderId, err)
	return
}

func testXt() {
	martXt := marts.GetMartByName(protos.XT)
	if err := martXt.Init(&protos.MartParamModel{
		AccessKey: "362a8910-dee3-450f-a6a9-3042f5c87477",
		SecretKey: "07eb576761c171ad4bd0646caf533c60b14fe31d",
	}); err != nil {
		panic(err)
	}
	// res, err := martXt.GetPrice("MOLI_USDT")
	// fmt.Println(res, err)

	// depth
	// res, err := martXt.GetDepth("MOLI_USDT")
	// fmt.Println(res, err)

	// get_balance
	res, err := martXt.GetSpotAccountAssets()
	fmt.Println(res, err)

	// data := protos.OrderProp{
	// 	ClientOrderId: fmt.Sprintf("%v", time.Now().UnixMicro()),
	// 	Symbol:        "MOLI_USDT",
	// 	Side:          protos.ORDER_SIDE_BUY,
	// 	Type:          protos.ORDER_TYPE_LIMIT,
	// 	Price:         "0.100",
	// 	Quantity:      "1",
	// }
	// oid, err := martXt.NewOrder(&data)
	// fmt.Println("NewOrder", oid, err)

	// oid, err = martXt.CancelOrder("MOLI_USDT", oid)
	// fmt.Println("CancelOrder", oid, err)

}

func testBitMart() {
	martBitmart := marts.GetMartByName(protos.BITMART)
	if err := martBitmart.Init(&protos.MartParamModel{
		AccessKey: "0c683f1898d84a81fc70bb6e621dd6d88b0af05e",
		SecretKey: "bd815ae93a90e82c1f048683eedc749b5071430c63a3557ee520b349b29b0c62",
		Memo:      "bot-1",
	}); err != nil {
		panic(err)
	}
	fmt.Println(">>>", martBitmart.GetDomain())

	res, err := martBitmart.GetPrice("MOLI_USDT")
	fmt.Printf("%#v %v\n", res, err)

	// depth
	// res, err := martBitmart.GetDepth("MOLI_USDT", 10)
	// fmt.Printf("%#v %v\n", res, err)

	// // get_balance
	// res, err := martBitmart.GetSpotAccountAssets()
	// fmt.Println(res, err)

	data := protos.OrderProp{
		ClientOrderId: fmt.Sprintf("%v", time.Now().UnixMilli()),
		Symbol:        "MOLI_USDT",
		Side:          protos.ORDER_SIDE_SELL,
		Type:          protos.ORDER_TYPE_LIMIT,
		Price:         "0.100",
		Quantity:      "100",
	}
	id, err := martBitmart.NewOrder(&data)
	fmt.Println(id, err)

	martBitmart.CancelOrder("MOLI_USDT", id)
	// rst, err := martXt.LatestKLine("MOLI_USDT")
	// fmt.Println("kline>>>", rst, err)

}
