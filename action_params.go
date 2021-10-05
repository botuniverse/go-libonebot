package libonebot

import (
	"encoding/base64"
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
func (p *ParamGetter) GetBool(key string) (bool, bool) {
	val, err := p.params.GetBool(key)
	if err != nil {
		p.w.WriteFailed(RetCodeBadParam, errorParam(err))
		return val, false
	}
	return val, true
}

// GetInt64 获取一个整数类型参数.
func (p *ParamGetter) GetInt64(key string) (int64, bool) {
	val, err := p.params.GetInt64(key)
	if err != nil {
		p.w.WriteFailed(RetCodeBadParam, errorParam(err))
		return val, false
	}
	return val, true
}

// GetFloat64 获取一个浮点数类型参数.
func (p *ParamGetter) GetFloat64(key string) (float64, bool) {
	val, err := p.params.GetFloat64(key)
	if err != nil {
		p.w.WriteFailed(RetCodeBadParam, errorParam(err))
		return val, false
	}
	return val, true
}

// GetString 获取一个字符串类型参数.
func (p *ParamGetter) GetString(key string) (string, bool) {
	val, err := p.params.GetString(key)
	if err != nil {
		p.w.WriteFailed(RetCodeBadParam, errorParam(err))
		return val, false
	}
	return val, true
}

// GetBytesOrBase64 获取一个字节数组类型参数.
func (p *ParamGetter) GetBytesOrBase64(key string) ([]byte, bool) {
	b, err := p.params.GetBytes(key)
	if err == nil {
		return b, true
	}
	s, err := p.params.GetString(key)
	if err != nil {
		p.w.WriteFailed(RetCodeBadParam, errorParam(err))
		return nil, false
	}
	b, err = base64.StdEncoding.DecodeString(s)
	if err != nil {
		p.w.WriteFailed(RetCodeBadParam, errorParam(err))
		return nil, false
	}
	return b, true
}

// GetMessage 获取一个消息类型参数.
func (p *ParamGetter) GetMessage(key string) (Message, bool) {
	val, err := p.params.GetMessage(key)
	if err != nil {
		p.w.WriteFailed(RetCodeBadParam, errorParam(err))
		return val, false
	}
	return val, true
}
