package onebot

import (
	"fmt"

	"github.com/tidwall/gjson"
)

type Params struct {
	JSON gjson.Result
}

func errorMissingParam(key string) error {
	return fmt.Errorf("参数 `%v` 不存在", key)
}

func errorInvalidParam(key string) error {
	return fmt.Errorf("参数 `%v` 无效", key)
}

func (params *Params) Get(key string) (gjson.Result, error) {
	val := params.JSON.Get(key)
	if !val.Exists() {
		return gjson.Result{}, errorMissingParam(key)
	}
	return val, nil
}

func (params *Params) GetBool(key string) (bool, error) {
	val := params.JSON.Get(key)
	if !val.Exists() {
		return false, errorMissingParam(key)
	}
	if val.Type != gjson.True && val.Type != gjson.False {
		return false, errorInvalidParam(key)
	}
	return val.Bool(), nil
}

func (params *Params) GetInt(key string) (int64, error) {
	val := params.JSON.Get(key)
	if !val.Exists() {
		return 0, errorMissingParam(key)
	}
	if val.Type != gjson.Number {
		return 0, errorInvalidParam(key)
	}
	return val.Int(), nil
}

func (params *Params) GetString(key string) (string, error) {
	val := params.JSON.Get(key)
	if !val.Exists() {
		return "", errorMissingParam(key)
	}
	if val.Type != gjson.String {
		return "", errorInvalidParam(key)
	}
	return val.Str, nil
}

func (params *Params) GetMessage(key string) (Message, error) {
	val, err := params.Get(key)
	if err != nil {
		return nil, err
	}

	if val.Type == gjson.String {
		return Message{TextSegment(val.Str)}, nil
	}

	if val.IsObject() {
		return MessageFromJSON("[" + val.Raw + "]")
	} else if val.IsArray() {
		return MessageFromJSON(val.Raw)
	} else {
		return nil, errorInvalidParam(key)
	}
}

type ParamGetter struct {
	params *Params
	w      ResponseWriter
}

func NewParamGetter(params *Params, w ResponseWriter) *ParamGetter {
	return &ParamGetter{params, w}
}

func (getter *ParamGetter) Get(key string) (gjson.Result, bool) {
	val, err := getter.params.Get(key)
	if err != nil {
		getter.w.WriteFailed(RetCodeParamError, err)
		return val, false
	}
	return val, true
}

func (getter *ParamGetter) GetBool(key string) (bool, bool) {
	val, err := getter.params.GetBool(key)
	if err != nil {
		getter.w.WriteFailed(RetCodeParamError, err)
		return val, false
	}
	return val, true
}

func (getter *ParamGetter) GetInt(key string) (int64, bool) {
	val, err := getter.params.GetInt(key)
	if err != nil {
		getter.w.WriteFailed(RetCodeParamError, err)
		return val, false
	}
	return val, true
}

func (getter *ParamGetter) GetString(key string) (string, bool) {
	val, err := getter.params.GetString(key)
	if err != nil {
		getter.w.WriteFailed(RetCodeParamError, err)
		return val, false
	}
	return val, true
}

func (getter *ParamGetter) GetMessage(key string) (Message, bool) {
	val, err := getter.params.GetMessage(key)
	if err != nil {
		getter.w.WriteFailed(RetCodeParamError, err)
		return val, false
	}
	return val, true
}
