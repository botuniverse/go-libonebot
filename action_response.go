package libonebot

import (
	"bytes"
	"encoding/json"
	"errors"

	"github.com/vmihailenco/msgpack/v5"
)

const (
	statusOK     = "ok"
	statusFailed = "failed"
)

// RetCodeXxx 表示动作响应返回码.
const (
	RetCodeOK = 0 // 成功

	// 动作请求错误 (类似 HTTP 的 4xx 客户端错误)
	RetCodeRequestErrorBase       = 10000
	RetCodeBadRequest             = 10001 // 无效的动作请求 (格式错误, 必要字段缺失或字段类型错误)
	RetCodeUnsupportedAction      = 10002 // 不支持的动作请求 (OneBot 实现没有实现该动作)
	RetCodeBadParam               = 10003 // 无效的动作请求参数 (参数缺失或参数类型错误)
	RetCodeUnsupportedParam       = 10004 // 不支持的动作请求参数 (OneBot 实现没有实现该参数的语义)
	RetCodeUnsupportedSegment     = 10005 // 不支持的消息段类型 (OneBot 实现没有实现该消息段类型)
	RetCodeBadSegmentData         = 10006 // 无效的消息段参数 (参数缺失或参数类型错误)
	RetCodeUnsupportedSegmentData = 10007 // 不支持的消息段参数 (OneBot 实现没有实现该参数的语义)

	// 动作处理器错误 (类似 HTTP 的 5xx 服务端错误)
	RetCodeHandlerErrorBase     = 20000
	RetCodeBadHandler           = 20001 // 动作处理器实现错误 (如没有正确设置响应状态等)
	RetCodeInternalHandlerError = 20002 // 动作处理器运行时抛出异常

	// 动作执行错误 (OneBot 实现可根据需要在低三位细分)
	RetCodeExecutionErrorBase = 30000
	RetCodeDatabaseError      = 31000 // 数据库错误
	RetCodeFilesystemError    = 32000 // 文件系统错误
	RetCodeNetworkError       = 33000 // 网络错误
	RetCodePlatformError      = 34000 // 聊天平台错误 (如由于聊天平台限制导致消息发送失败等)
	RetCodeLogicError         = 35000 // 动作逻辑错误 (如尝试向不存在的用户发送消息等)
	RetCodeIAmTired           = 36000 // 我不想干了 (一位 OneBot 实现决定罢工)

	// 保留错误段
	RetCodeReservedErrorBase1 = 40000
	RetCodeReservedErrorBase2 = 50000
)

// Response 表示一个动作响应.
type Response struct {
	Status  string      `json:"status"`         // 执行状态 (成功与否)
	RetCode int         `json:"retcode"`        // 返回码
	Data    interface{} `json:"data"`           // 响应数据
	Message string      `json:"message"`        // 错误信息
	Echo    interface{} `json:"echo,omitempty"` // 动作请求的 echo 字段 (原样返回)
}

func failedResponse(retCode int, err error) Response {
	return Response{
		Status:  statusFailed,
		RetCode: retCode,
		Message: err.Error(),
	}
}

func (r Response) encode(isBinary bool) ([]byte, error) {
	if r.Status != statusOK && r.Status != statusFailed {
		return nil, errors.New("`status` 字段值无效")
	}
	if isBinary {
		var buf bytes.Buffer
		enc := msgpack.NewEncoder(&buf)
		enc.SetCustomStructTag("json")
		err := enc.Encode(r)
		return buf.Bytes(), err
	}
	return json.Marshal(r)
}
