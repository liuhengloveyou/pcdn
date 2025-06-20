package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"pcdn-server/api"
	"pcdn-server/common"
	"pcdn-server/repos"
	"pcdn-server/tasks"
	"pcdn-server/tcpservice"

	"github.com/qiniu/go-sdk/v7/auth"
)

var (
	BuildTime string
	CommitID  string

	Sig string
)

var (
	showVer = flag.Bool("version", false, "打印版本号")
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
		common.Logger.Sugar().Warn("Got signal:", s)
	}()
}

const ak = "20Ko8rtFdLh4YZcf_QkC_Hj4aMYsLs3dpJWNUOqc"
const sk = "gs63OaLilp0qR_74BICW-XFeNRPddcVa6JrF0C_r"

type NodeStatReq struct {
	NodeIDs []string `json:"nodeIDs"`
	Start   string   `json:"start"`
	End     string   `json:"end"`
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

	repos.InitRepos()
	go tcpservice.InitTcpService(common.ServConfig.TcpServerAddr)
	go tasks.RunTasks()
	go tcpservice.RunTcpTasks()

	if err := api.InitAndRunHttpApi(common.ServConfig.HttpServerAddr); err != nil {
		panic("HTTPAPI: " + err.Error())
	}
}

func InitSystemEnv() {
	if err := repos.InitModels(common.OrmCli); err != nil {
		panic(err)
	}
}

func qiniuGet() {
	params := url.Values{}
	params.Add("VendorDeviceIds", "[ld12345678902]")
	// params.Add("isp", "移动")
	// params.Add("usbw", fmt.Sprint(1000))
	// params.Add("bwNum", fmt.Sprint(5))
	// params.Add("natType", "public")
	// params.Add("dialType", "staticNetSingle")

	reqUrl := "https://api.niulinkcloud.com/v1/ant/mapping"
	reqUrl = fmt.Sprintf("%s?%s", reqUrl, params.Encode())
	fmt.Println("rqurl: ", reqUrl)
	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		return
	}
	// 设置Content-Type必须在调用SignRequestV2方法之前，因为Content-Type也需要参与签名
	req.Header.Set("Content-Type", "application/json")
	token, err := auth.New(ak, sk).SignRequestV2(req)
	if err != nil {
		return
	}
	req.Header.Set("Authorization", "Qiniu "+token)

	httpCli := &http.Client{}
	httpResp, err := httpCli.Do(req)
	if err != nil {
		return
	}
	respBody, err := io.ReadAll(httpResp.Body)
	defer httpResp.Body.Close()
	if err != nil {
		return
	}
	if httpResp.StatusCode == 200 {
		println(fmt.Sprintf("resp body:%s", string(respBody)))
	} else {
		println(fmt.Sprintf("http resp code: %d, resp body:%s", httpResp.StatusCode, string(respBody)))
	}
}

func qiniuPost() {
	reqData := map[string]interface{}{
		"vendorDeviceIds": []string{"ld12345678902"},
	}
	reqDataBytes, err := json.Marshal(reqData)
	if err != nil {
		return
	}
	reqUrl := "https://api.niulinkcloud.com/v1/ant/mapping"
	req, err := http.NewRequest("POST", reqUrl, bytes.NewReader(reqDataBytes))
	if err != nil {
		return
	}
	// 设置Content-Type必须在调用SignRequestV2方法之前，因为Content-Type也需要参与签名
	req.Header.Set("Content-Type", "application/json")
	token, err := auth.New(ak, sk).SignRequestV2(req)
	if err != nil {
		return
	}
	req.Header.Set("Authorization", "Qiniu "+token)

	httpCli := &http.Client{}
	httpResp, err := httpCli.Do(req)
	if err != nil {
		return
	}
	respBody, err := ioutil.ReadAll(httpResp.Body)
	defer httpResp.Body.Close()
	if err != nil {
		return
	}
	if httpResp.StatusCode == 200 {
		println(fmt.Sprintf("resp body:%s", string(respBody)))
	} else {
		println(fmt.Sprintf("http resp code: %d, resp body:%s", httpResp.StatusCode, string(respBody)))
	}
}
