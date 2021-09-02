package libonebot

import (
	"fmt"
)

// ParamGetter 用于在处理动作请求时方便地获取参数, 当参数不存在或参数错误时,
// 向 ResponseWriter 写入错误信息.
type ParamGetter struct {
	params EasierMap
	w      ResponseWriter
}

// NewParamGetter 创建一个 ParamGetter 对象.
func NewParamGetter(w ResponseWriter, r *Request) *ParamGetter {
	return &ParamGetter{
		params: r.Params,
		w:      w,
	}
}

func errorParam(err error) error {
	return fmt.Errorf("参数错误: %v", err)
}

// GetBool 获取一个布尔类型参数.
func (getter *ParamGetter) GetBool(key string) (bool, bool) {
	val, err := getter.params.GetBool(key)
	if err != nil {
		getter.w.WriteFailed(RetCodeParamError, errorParam(err))
		return val, false
	}
	return val, true
}

// GetInt64 获取一个整数类型参数.
func (getter *ParamGetter) GetInt64(key string) (int64, bool) {
	val, err := getter.params.GetInt64(key)
	if err != nil {
		getter.w.WriteFailed(RetCodeParamError, errorParam(err))
		return val, false
	}
	return val, true
}

// GetString 获取一个字符串类型参数.
func (getter *ParamGetter) GetString(key string) (string, bool) {
	val, err := getter.params.GetString(key)
	if err != nil {
		getter.w.WriteFailed(RetCodeParamError, errorParam(err))
		return val, false
	}
	return val, true
}

// GetMessage 获取一个消息类型参数.
func (getter *ParamGetter) GetMessage(key string) (Message, bool) {
	val, err := getter.params.GetMessage(key)
	if err != nil {
		getter.w.WriteFailed(RetCodeParamError, errorParam(err))
		return val, false
	}
	return val, true
}
