package libonebot

import (
	"bytes"
	"context"
	"net/http"
	"net/url"
	"sync"
	"time"
)

func commRunHTTPWebhook(c ConfigCommHTTPWebhook, ob *OneBot, ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	ob.Logger.Infof("正在启动 HTTP Webhook (%v)...", c.URL)

	u, err := url.Parse(c.URL)
	if err != nil {
		ob.Logger.Errorf("HTTP Webhook (%v) 启动失败, URL 不合法, 错误: %v", c.URL, err)
		return
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		ob.Logger.Errorf("HTTP Webhook (%v) 启动失败, URL 不合法, 必须使用 HTTP 或 HTTPS 协议", c.URL)
		return
	}

	eventChan := ob.openEventListenChan()
	httpClient := &http.Client{
		Timeout: time.Duration(c.Timeout) * time.Second, // 0 for no timeout
	}

	for {
		select {
		case event := <-eventChan:
			ob.Logger.Debugf("通过 HTTP Webhook (%v) 推送事件 `%v`", c.URL, event.name)
			req, _ := http.NewRequest(http.MethodPost, c.URL, bytes.NewReader(event.bytes))
			req.Header.Set("Content-Type", "application/json")
			if c.AccessToken != "" {
				req.Header.Set("Authorization", "Bearer "+c.AccessToken)
			}
			req.Header.Set("User-Agent", ob.GetUserAgent())
			req.Header.Set("X-OneBot-Version", OneBotVersion)
			req.Header.Set("X-Platform", ob.Platform)
			req.Header.Set("X-Self-ID", ob.SelfID)
			resp, err := httpClient.Do(req)
			if err != nil {
				ob.Logger.Errorf("通过 HTTP Webhook (%v) 推送事件 `%v` 失败, 错误: %v", c.URL, event.name, err)
				continue
			}
			if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
				ob.Logger.Errorf("通过 HTTP Webhook (%v) 推送事件 `%v` 失败, 状态码: %v", c.URL, event.name, resp.StatusCode)
				continue
			}
			// TODO: call actions
		case <-ctx.Done():
			ob.closeEventListenChan(eventChan)
			ob.Logger.Infof("HTTP Webhook (%v) 已关闭", c.URL)
			return
		}
	}
}
