package onebot

import (
	"bytes"
	"net/http"
	"net/url"
	"sync"

	log "github.com/sirupsen/logrus"
)

func commStartHTTPWebhook(urlString string, onebot *OneBot) commCloser {
	log.Infof("正在启动 HTTP Webhook (%v)...", urlString)

	uri, err := url.Parse(urlString)
	if err != nil {
		log.Warnf("HTTP Webhook (%v) 启动失败, URL 不合法, 错误: %v", urlString, err)
		return nil
	}
	if uri.Scheme != "http" && uri.Scheme != "https" {
		log.Warnf("HTTP Webhook (%v) 启动失败, URL 不合法, 必须使用 HTTP 或 HTTPS 协议", urlString)
		return nil
	}

	eventChan := onebot.openEventListenChan()
	httpClient := &http.Client{}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for event := range eventChan {
			// TODO: use special User-Agent
			// TODO: check status code
			// TODO: timeout
			log.Debugf("通过 HTTP Webhook (%v) 推送事件 %v", urlString, event.name)
			httpClient.Post(urlString, "application/json", bytes.NewReader(event.bytes))
		}
		log.Infof("HTTP Webhook (%v) 已关闭", urlString)
	}()

	return func() {
		onebot.closeEventListenChan(eventChan)
		wg.Wait()
	}
}
