package libonebot

import (
	"bytes"
	"net/http"
	"net/url"
	"sync"

	log "github.com/sirupsen/logrus"
)

func commStartHTTPWebhook(c ConfigCommHTTPWebhook, ob *OneBot) commCloser {
	log.Infof("正在启动 HTTP Webhook (%v)...", c.URL)

	uri, err := url.Parse(c.URL)
	if err != nil {
		log.Warnf("HTTP Webhook (%v) 启动失败, URL 不合法, 错误: %v", c.URL, err)
		return nil
	}
	if uri.Scheme != "http" && uri.Scheme != "https" {
		log.Warnf("HTTP Webhook (%v) 启动失败, URL 不合法, 必须使用 HTTP 或 HTTPS 协议", c.URL)
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
			log.Debugf("通过 HTTP Webhook (%v) 推送事件 %v", c.URL, event.name)
			httpClient.Post(c.URL, "application/json", bytes.NewReader(event.bytes))
		}
		log.Infof("HTTP Webhook (%v) 已关闭", c.URL)
	}()

	return func() {
		ob.closeEventListenChan(eventChan)
		wg.Wait()
	}
}
