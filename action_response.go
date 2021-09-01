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

const (
	RetCodeOK = 0

	// Action request error
	RetCodeInvalidRequest = 11001
	RetCodeActionNotFound = 11002
	RetCodeParamError     = 11003

	// Action execution error
	RetCodeDatabaseError   = 12100
	RetCodeFilesystemError = 12200
	RetCodePlatformError   = 12300
	RetCodeLogicError      = 12400

	// Action handler error
	RetCodeBadActionHandler = 13001
)

type Response struct {
	Status  actionStatus `json:"status"`
	RetCode int          `json:"retcode"`
	Data    interface{}  `json:"data"`
	Message string       `json:"message"`
	Echo    interface{}  `json:"echo,omitempty"`
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
