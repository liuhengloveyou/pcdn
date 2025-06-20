package api

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"pcdn-server/common"

	gocommon "github.com/liuhengloveyou/go-common"
	passport "github.com/liuhengloveyou/passport/face"
	passportprotos "github.com/liuhengloveyou/passport/protos"
	"go.uber.org/zap"
)

type ApiStruct struct {
	Handler    func(http.ResponseWriter, *http.Request)
	Method     string
	NeedLogin  bool
	NeedAccess bool
}

var (
	Logger *zap.SugaredLogger

	Apis map[string]ApiStruct = make(map[string]ApiStruct, 1000)
)

func init() {
	Logger = common.Logger.Sugar()

	Apis = map[string]ApiStruct{
		"/uploadImgByForm": {
			Handler: uploadImgByForm,
		},
	}

	initTcApi()
	initDeviceManagerApi()
	initBusinessLogApi()
}

func InitAndRunHttpApi(addr string) error {
	// http.HandleFunc("/ws", ws.WebsocketHandler)
	passport.InitAndRunHttpApi(nil)

	http.Handle("/", &HttpApiServer{})
	// http.Handle("/", http.FileServer(http.Dir("www")))

	s := &http.Server{
		Addr:           addr,
		ReadTimeout:    10 * time.Minute,
		WriteTimeout:   10 * time.Minute,
		MaxHeaderBytes: 1 << 20,
	}

	fmt.Println("go: ", addr)
	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

type HttpApiServer struct{}

func (p *HttpApiServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin") //请求头部
	if origin != "" {
		w.Header().Add("Access-Control-Allow-Origin", origin)
		w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		w.Header().Add("Access-Control-Allow-Headers", "Origin, X-Requested-With, X-Extra-Header, Content-Type, Accept, Authorization")
		w.Header().Add("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
		w.Header().Add("Access-Control-Allow-Credentials", "true")
	}
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	apiName := r.Header.Get("X-API")
	if apiName == "" {
		URL, _ := url.Parse(r.RequestURI)
		apiName = URL.Path
	}
	fmt.Println("api: ", apiName)

	if apiName == "" {
		Logger.Warnf("api not found: %v\n", apiName)
		gocommon.HttpErr(w, http.StatusMethodNotAllowed, -1, "?API")
		return
	}

	apiHandler, ok := Apis[apiName]
	if !ok {
		gocommon.HttpErr(w, http.StatusMethodNotAllowed, -1, apiName)
		return
	}

	if apiHandler.Method != "" && apiHandler.Method != r.Method {
		gocommon.HttpErr(w, http.StatusMethodNotAllowed, -1, "method not allowed")
		return
	}

	if apiHandler.NeedLogin {
		if ReadSessionFromRequest(r) == nil {
			gocommon.HttpJsonErr(w, http.StatusUnauthorized, common.ErrNoAuth)
			return
		}
	}

	apiHandler.Handler(w, r)
}

func ReadSessionFromRequest(r *http.Request) (sessionUser *passportprotos.User) {
	sessionStr := r.Header.Get("go-sess")
	if sessionStr != "" {
		sessionUserJson := struct {
			Data passportprotos.User `json:"data"`
		}{}

		if err := json.Unmarshal([]byte(sessionStr), &sessionUserJson); err != nil {
			Logger.Errorf("ReadSessionFromRequest json ERR: %v\n", err)
			return
		}

		return &sessionUserJson.Data
	}

	sessUser := passport.GetSessionUser(r)
	if sessUser.UID <= 0 {
		Logger.Errorf("passport.GetSessionUser ERR: %#v\n", sessUser)
		return
	}

	return &sessUser
}

func uploadImgByForm(w http.ResponseWriter, r *http.Request) {
	var (
		dir string
		fp  string
	)

	sessionUser := ReadSessionFromRequest(r)
	if sessionUser.TenantID <= 0 || sessionUser.UID <= 0 {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrNoAuth)
		return
	}

	flen, _ := strconv.ParseInt(r.Header.Get("Content-Length"), 10, 64)
	if flen == 0 || flen > common.MAX_UPLOAD_LEN {
		Logger.Error("uploadImgByForm Content-Length ERR: ", flen)
		gocommon.HttpErr(w, http.StatusBadRequest, -1, "文件太大")
		return

	}

	r.ParseMultipartForm(32 << 20)
	act := r.FormValue("act")
	fileType := r.FormValue("type")
	if fileType == "" {
		Logger.Error("uploadImgByForm fileType nil")
		gocommon.HttpErr(w, http.StatusBadRequest, -1, "文件类型错误")
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		Logger.Error("uploadImgByForm FormFile err: ", err)
		gocommon.HttpErr(w, http.StatusBadRequest, -1, "读上传文件错误")
		return
	}
	defer file.Close()

	fileBuff, err := io.ReadAll(file)
	if err != nil {
		Logger.Error("uploadImgByForm ReadAll err: ", err)
		gocommon.HttpErr(w, http.StatusBadRequest, -1, "读上传文件错误")
		return
	}

	md5Val := fmt.Sprintf("%x", md5.Sum(fileBuff))
	imgDir := common.ServConfig.ImgDir
	if act == "StoreIn" {
		imgDir = common.ServConfig.SotreImgDir
	}

	dir = fmt.Sprintf("%s/%s", imgDir, md5Val[:2])
	if err := os.MkdirAll(dir, 0755); err != nil {
		gocommon.HttpErr(w, http.StatusOK, -1, "文件系统错误")
		Logger.Error("uploadImgByForm mkdir ERR: ", dir, err)
		return
	}

	fp = fmt.Sprintf("%s/%s%s", dir, md5Val, fileType)
	Logger.Info("uploadImgByForm fn: ", fp)

	if !gocommon.IsExists(fp) {
		if err := os.WriteFile(fp, fileBuff, 0755); err != nil {
			Logger.Error("FileUpload err: ", err)
			gocommon.HttpErr(w, http.StatusInternalServerError, -1, "写文件失败")
			return
		}
	}

	Logger.Info("uploadImgByForm ok: ", fp)

	gocommon.HttpErr(w, http.StatusOK, 0, fmt.Sprintf("%s/%s%s", md5Val[:2], md5Val, fileType))
}
