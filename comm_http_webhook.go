package libonebot

import (
	"bytes"
	"context"
	"net/http"
	"net/url"
	"sync"
)

func commRunHTTPWebhook(c ConfigCommHTTPWebhook, ob *OneBot, ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	ob.Logger.Infof("正在启动 HTTP Webhook (%v)...", c.URL)

	u, err := url.Parse(c.URL)
	if err != nil {
		ob.Logger.Warnf("HTTP Webhook (%v) 启动失败, URL 不合法, 错误: %v", c.URL, err)
		return
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		ob.Logger.Warnf("HTTP Webhook (%v) 启动失败, URL 不合法, 必须使用 HTTP 或 HTTPS 协议", c.URL)
		return
	}

	eventChan := ob.openEventListenChan()
	httpClient := &http.Client{}

	for {
		select {
		case event := <-eventChan:
			// TODO: use special User-Agent
			// TODO: check status code
			// TODO: timeout
			ob.Logger.Debugf("通过 HTTP Webhook (%v) 推送事件 `%v`", c.URL, event.name)
			httpClient.Post(c.URL, "application/json", bytes.NewReader(event.bytes))
		case <-ctx.Done():
			ob.closeEventListenChan(eventChan)
			ob.Logger.Infof("HTTP Webhook (%v) 已关闭", c.URL)
			return
		}
	}
}
