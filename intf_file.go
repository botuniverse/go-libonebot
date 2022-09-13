// 接口定义 - 文件接口

package libonebot

// 文件动作
// https://12.onebot.dev/interface/file/actions/

const (
	ActionUploadFile           = "upload_file"            // 上传文件
	ActionUploadFileFragmented = "upload_file_fragmented" // 分片上传文件
	ActionGetFile              = "get_file"               // 获取文件
	ActionGetFileFragmented    = "get_file_fragmented"    // 分片获取文件
)
