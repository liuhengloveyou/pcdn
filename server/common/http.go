package common

import (
	"io"
	"net/http"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/go-playground/validator/v10"
	"github.com/liuhengloveyou/passport/common"
)

func ReadJsonBodyFromRequest(r *http.Request, dst interface{}, validate string) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	Logger.Sugar().Debugf("readJsonBodyFromRequest body: '%v'\n", string(body))

	if err = sonic.Unmarshal(body, dst); err != nil {
		Logger.Sugar().Errorf("readJsonBodyFromRequest JSON: '%v' %v\n", string(body), err)
		return err
	}

	if len(validate) == 0 {
		return nil
	}

	if strings.Compare("struct", validate) == 0 {
		if err := common.Validate.Struct(dst); err != nil {
			Logger.Sugar().Errorf("Validate ERR: %v\n", err)
			if _, ok := err.(*validator.InvalidValidationError); !ok {
				return common.ErrParam
			}
		}
	}

	return nil
}
