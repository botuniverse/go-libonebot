package libonebot

import (
	"bytes"
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"
)

func commStartHTTPWebhook(urlString string, onebot *OneBot) {
	log.Infof("正在启动 HTTP Webhook (%v)...", urlString)

	uri, err := url.Parse(urlString)
	if err != nil {
		log.Warnf("HTTP Webhook (%v) 启动失败, URL 不合法, 错误: %v", urlString, err)
		return
	}
	if uri.Scheme != "http" && uri.Scheme != "https" {
		log.Warnf("HTTP Webhook (%v) 启动失败, URL 不合法, 必须使用 HTTP 或 HTTPS 协议", urlString)
		return
	}

	eventChan := onebot.openEventListenChan()
	httpClient := &http.Client{}
	go func() {
		for eventBytes := range eventChan {
			// TODO: use special User-Agent
			// TODO: check status code
			httpClient.Post(urlString, "application/json", bytes.NewReader(eventBytes))
		}
		log.Warnf("HTTP Webhook (%v) 已关闭", urlString)
	}()
}
