package libonebot

import (
	"fmt"
)

type ParamGetter struct {
	params *easyMap
	w      ResponseWriter
}

func NewParamGetter(params *easyMap, w ResponseWriter) *ParamGetter {
	return &ParamGetter{params, w}
}

func errorParam(err error) error {
	return fmt.Errorf("参数错误: %v", err)
}

func (getter *ParamGetter) GetBool(key string) (bool, bool) {
	val, err := getter.params.GetBool(key)
	if err != nil {
		getter.w.WriteFailed(RetCodeParamError, errorParam(err))
		return val, false
	}
	return val, true
}

func (getter *ParamGetter) GetInt64(key string) (int64, bool) {
	val, err := getter.params.GetInt64(key)
	if err != nil {
		getter.w.WriteFailed(RetCodeParamError, errorParam(err))
		return val, false
	}
	return val, true
}

func (getter *ParamGetter) GetString(key string) (string, bool) {
	val, err := getter.params.GetString(key)
	if err != nil {
		getter.w.WriteFailed(RetCodeParamError, errorParam(err))
		return val, false
	}
	return val, true
}

func (getter *ParamGetter) GetMessage(key string) (Message, bool) {
	val, err := getter.params.GetMessage(key)
	if err != nil {
		getter.w.WriteFailed(RetCodeParamError, errorParam(err))
		return val, false
	}
	return val, true
}
