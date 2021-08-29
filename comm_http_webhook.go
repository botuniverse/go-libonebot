package libonebot

import (
	"bytes"
	"net/http"
	"net/url"
	"sync"
)

func commStartHTTPWebhook(c ConfigCommHTTPWebhook, ob *OneBot) commCloser {
	ob.Logger.Infof("正在启动 HTTP Webhook (%v)...", c.URL)

	u, err := url.Parse(c.URL)
	if err != nil {
		ob.Logger.Warnf("HTTP Webhook (%v) 启动失败, URL 不合法, 错误: %v", c.URL, err)
		return nil
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		ob.Logger.Warnf("HTTP Webhook (%v) 启动失败, URL 不合法, 必须使用 HTTP 或 HTTPS 协议", c.URL)
		return nil
	}

	eventChan := ob.openEventListenChan()
	httpClient := &http.Client{}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for event := range eventChan {
			// TODO: use special User-Agent
			// TODO: check status code
			// TODO: timeout
			ob.Logger.Debugf("通过 HTTP Webhook (%v) 推送事件 %v", c.URL, event.name)
			httpClient.Post(c.URL, "application/json", bytes.NewReader(event.bytes))
		}
	}()

	return func() {
		ob.closeEventListenChan(eventChan)
		wg.Wait()
		ob.Logger.Infof("HTTP Webhook (%v) 已关闭", c.URL)
	}
}
