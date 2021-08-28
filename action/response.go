package action

import "encoding/json"

type ResponseStatus struct{ string }

func (status ResponseStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(status.string)
}

var (
	StatusOK     = ResponseStatus{"ok"}
	StatusFailed = ResponseStatus{"failed"}
)

const (
	RetCodeOK = 0

	RetCodeInvalidRequest = 11001
	RetCodeActionNotFound = 11002
	RetCodeMissingParam   = 11003
	RetCodeInvalidParam   = 11004

	RetCodeDatabaseError   = 12100
	RetCodeFilesystemError = 12200
	RetCodePlatformError   = 12300
	RetCodeLogicError      = 12400
)

type Response struct {
	Status  ResponseStatus `json:"status"`
	RetCode int            `json:"retcode"`
	Data    interface{}    `json:"data"`
	Message string         `json:"message"`
	Echo    interface{}    `json:"echo,omitempty"`
}

func OKResponse(data interface{}) Response {
	return Response{
		Status:  StatusOK,
		RetCode: RetCodeOK,
		Data:    data,
	}
}

func FailedResponse(retCode int, message string) Response {
	return Response{
		Status:  StatusFailed,
		RetCode: retCode,
		Message: message,
	}
}

type ResponseWriter struct {
	resp *Response
}

func (w ResponseWriter) WriteOK() {
	w.resp.Status = StatusOK
	w.resp.RetCode = RetCodeOK
	w.resp.Message = ""
}

func (w ResponseWriter) WriteData(data interface{}) {
	w.WriteOK()
	w.resp.Data = data
}

func (w ResponseWriter) WriteFailed(retCode int, message string) {
	w.resp.Status = StatusFailed
	w.resp.RetCode = retCode
	w.resp.Message = message
}
