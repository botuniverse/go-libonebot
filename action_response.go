package libonebot

import (
	"encoding/json"
)

type actionStatus struct{ string }

func (s actionStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.string)
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

type ResponseWriter struct {
	resp *Response
}

func (w ResponseWriter) WriteOK() {
	w.resp.Status = statusOK
	w.resp.RetCode = RetCodeOK
	w.resp.Message = ""
}

func (w ResponseWriter) WriteData(data interface{}) {
	w.WriteOK()
	w.resp.Data = data
}

func (w ResponseWriter) WriteFailed(retCode int, err error) {
	w.resp.Status = statusFailed
	w.resp.RetCode = retCode
	w.resp.Message = err.Error()
}
