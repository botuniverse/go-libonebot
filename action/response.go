package action

type ResponseStatus struct{ string }

var (
	StatusOK     = ResponseStatus{"ok"}
	StatusFailed = ResponseStatus{"failed"}
)

const (
	RetCodeOK = 0

	RetCodeInvalidRequest = 11001
	RetCodeMissingAction  = 11002
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
}

func FailedResponse(retCode int, message string) Response {
	return Response{
		Status:  StatusFailed,
		RetCode: retCode,
		Data:    nil,
		Message: message,
	}
}
