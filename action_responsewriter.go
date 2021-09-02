package libonebot

// ResponseWriter 封装了对 Response 的修改操作.
type ResponseWriter struct {
	resp *Response
}

// WriteOK 向 Response 写入成功状态.
func (w ResponseWriter) WriteOK() {
	w.resp.Status = statusOK
	w.resp.RetCode = RetCodeOK
	w.resp.Message = ""
}

// WriteData 向 Response 写入成功状态, 并写入返回数据.
func (w ResponseWriter) WriteData(data interface{}) {
	w.WriteOK()
	w.resp.Data = data
}

// WriteFailed 向 Response 写入失败状态, 并写入返回码和错误信息.
func (w ResponseWriter) WriteFailed(retCode int, err error) {
	w.resp.Status = statusFailed
	w.resp.RetCode = retCode
	w.resp.Message = err.Error()
}
