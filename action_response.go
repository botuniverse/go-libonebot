package libonebot

import (
	"bytes"
	"encoding/json"

	"github.com/vmihailenco/msgpack/v5"
)

type actionStatus struct{ string }

func (s actionStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.string)
}

func (s actionStatus) MarshalMsgpack() ([]byte, error) {
	return msgpack.Marshal(s.string)
}

var (
	statusOK     = actionStatus{"ok"}
	statusFailed = actionStatus{"failed"}
)

// RetCodeXxx 表示动作响应返回码.
const (
	RetCodeOK = 0 // 成功

	// 动作请求错误
	RetCodeInvalidRequest = 11001 // 动作请求无效 (格式错误, 必要字段缺失或字段类型错误)
	RetCodeActionNotFound = 11002 // 动作请求不存在 (OneBot 实现没有实现该动作)
	RetCodeParamError     = 11003 // 动作请求参数错误 (参数缺失或参数类型错误)

	// 动作执行错误
	RetCodeDatabaseError   = 12100 // 数据库错误
	RetCodeFilesystemError = 12200 // 文件系统错误
	RetCodePlatformError   = 12300 // 聊天平台错误
	RetCodeLogicError      = 12400 // 动作逻辑错误 (如尝试向不存在的用户发送消息等)

	// 动作处理器错误
	RetCodeBadActionHandler = 13001 // 动作处理器实现错误
)

// Response 表示一个动作响应.
type Response struct {
	Status  actionStatus `json:"status"`         // 执行状态 (成功与否)
	RetCode int          `json:"retcode"`        // 返回码
	Data    interface{}  `json:"data"`           // 返回数据
	Message string       `json:"message"`        // 错误信息
	Echo    interface{}  `json:"echo,omitempty"` // 动作请求的 echo 字段 (原样返回)
}

func failedResponse(retCode int, err error) Response {
	return Response{
		Status:  statusFailed,
		RetCode: retCode,
		Message: err.Error(),
	}
}

func (r Response) encode(isBinary bool) ([]byte, error) {
	if isBinary {
		var buf bytes.Buffer
		enc := msgpack.NewEncoder(&buf)
		enc.SetCustomStructTag("json")
		err := enc.Encode(r)
		return buf.Bytes(), err
	}
	return json.Marshal(r)
}
