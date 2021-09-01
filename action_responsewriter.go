package libonebot

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
