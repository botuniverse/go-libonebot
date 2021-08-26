package comm

import (
	"net/url"

	"github.com/botuniverse/go-libonebot/event"
	log "github.com/sirupsen/logrus"
)

// Start an HTTP Webhook communication task.
func StartHTTPWebhookTask(urlString string, eventEmitter *event.EventEmitter) {
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

	eventChan := eventEmitter.OpenOutChan()
	go func() {
		for eventBytes := range eventChan {
			log.Debugf("EventBytes: %v", eventBytes)
		}
		log.Warnf("HTTP Webhook (%v) 已关闭", urlString)
	}()
}
