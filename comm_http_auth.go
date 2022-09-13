// OneBot Connect - 通信方式 - HTTP - 鉴权
// https://12.onebot.dev/connect/communication/http/#_1

package libonebot

import "net/http"

type httpAuthorizer struct {
	accessToken string
}

func (auth *httpAuthorizer) authorize(r *http.Request) bool {
	if auth.accessToken != "" {
		// try authorize with `Authorization` header
		if r.Header.Get("Authorization") == "Bearer "+auth.accessToken {
			return true
		}
		// try authorize with `access_token` query parameter
		if r.URL.Query().Get("access_token") == auth.accessToken {
			return true
		}
		return false
	}
	return true
}
