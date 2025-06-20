package common

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	gocommon "github.com/liuhengloveyou/go-common"
)

func SaveBase64Img(name, data string) (string, error) {
	var (
		extName string = ""
		dir     string = ""
		fp      string = ""
	)

	dotIdx := strings.LastIndex(name, ".")
	if dotIdx > 0 {
		extName = name[dotIdx:]
	}
	if extName == "" {
		Logger.Sugar().Error("SaveBase64Img extName ERR: ", name)
		return "", fmt.Errorf("extName ERR")
	}

	base64Data := data[strings.Index(data, ",")+1:]
	fileBuff, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		fmt.Println("SaveBase64Img base64 ERR: ", dir, err)
		Logger.Sugar().Error("SaveBase64Img base64 ERR: ", dir, err)
		return "", fmt.Errorf("SaveBase64Img base64 ERR")
	}
	md5Val := fmt.Sprintf("%x", md5.Sum(fileBuff))

	dir = fmt.Sprintf("%s/%s", ServConfig.ImgDir, md5Val[:2])
	if err := os.MkdirAll(dir, 0755); err != nil {
		Logger.Sugar().Error("SaveBase64Img mkdir ERR: ", dir, err)
		return "", fmt.Errorf("SaveBase64Img mkdir ERR")
	}

	fp = fmt.Sprintf("%s/%s%s", dir, md5Val, extName)
	Logger.Sugar().Info("SaveBase64Img fn: ", fp)

	if !gocommon.IsExists(fp) {
		if err := os.WriteFile(fp, fileBuff, 0755); err != nil {
			Logger.Sugar().Error("FileUpload err: ", err)
			return "", fmt.Errorf("WriteFile ERR")
		}
	}

	Logger.Sugar().Info("SaveBase64Img ok: ", fp)

	return fmt.Sprintf("%s/%s%s", md5Val[:2], md5Val, extName), nil
}
