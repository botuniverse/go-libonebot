package libonebot

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type httpWebhookComm struct {
	ob          *OneBot
	url         string
	accessToken string
	httpClient  *http.Client
}

func (comm *httpWebhookComm) post(event marshaledEvent) {
	req, _ := http.NewRequest(http.MethodPost, comm.url, bytes.NewReader(event.bytes))
	req.Header.Set("Content-Type", "application/json")
	if comm.accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+comm.accessToken)
	}
	req.Header.Set("User-Agent", comm.ob.GetUserAgent())
	req.Header.Set("X-OneBot-Version", OneBotVersion)
	req.Header.Set("X-Impl", comm.ob.Impl)
	req.Header.Set("X-Platform", comm.ob.Platform)
	req.Header.Set("X-Self-ID", comm.ob.SelfID)

	resp, err := comm.httpClient.Do(req)
	if err != nil {
		comm.ob.Logger.Errorf("通过 HTTP Webhook (%v) 推送事件 `%v` 失败, 错误: %v", comm.url, event.name, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		comm.ob.Logger.Errorf("通过 HTTP Webhook (%v) 推送事件 `%v` 失败, 状态码: %v", comm.url, event.name, resp.StatusCode)
		return
	}

	if resp.StatusCode == http.StatusOK {
		// handle action requests in the response body
		var isBinary bool
		contentType := resp.Header.Get("Content-Type")
		if strings.HasPrefix(contentType, "application/json") {
			isBinary = false
			contentType = "application/json"
		} else if strings.HasPrefix(contentType, "application/msgpack") {
			isBinary = true
			contentType = "application/msgpack"
		} else {
			// reject unsupported content types
			comm.ob.Logger.Warnf("响应头中的 Content-Type 不支持, 已忽略")
			return
		}
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			comm.ob.Logger.Warnf("动作请求列表读取失败, 已忽略, 错误: %v", err)
			return
		}
		requests, err := decodeRequestList(respBody, isBinary)
		if err != nil {
			comm.ob.Logger.Warnf("动作请求列表解析失败, 已忽略, 错误: %v", err)
			return
		}
		for _, request := range requests {
			comm.ob.handleRequest(&request) // response is ignored
		}
	}
}

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

	comm := &httpWebhookComm{
		ob:          ob,
		url:         c.URL,
		accessToken: c.AccessToken,
		httpClient: &http.Client{
			Timeout: time.Duration(c.Timeout) * time.Millisecond, // 0 for no timeout
		},
	}

	eventChan := ob.openEventListenChan()
	defer ob.closeEventListenChan(eventChan)

	for {
		select {
		case event := <-eventChan:
			comm.ob.Logger.Debugf("通过 HTTP Webhook (%v) 推送事件 `%v`", comm.url, event.name)
			go comm.post(event)
		case <-ctx.Done():
			ob.Logger.Infof("HTTP Webhook (%v) 已关闭", c.URL)
			return
		}
	}
}
