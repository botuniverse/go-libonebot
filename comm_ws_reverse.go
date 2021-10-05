package libonebot

import (
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

func commStartWSReverse(c ConfigCommWSReverse, ob *OneBot) commCloser {
	ob.Logger.Infof("正在启动 WebSocket Reverse (%v)...", c.URL)

	u, err := url.Parse(c.URL)
	if err != nil {
		ob.Logger.Warnf("WebSocket Reverse (%v) 启动失败, URL 不合法, 错误: %v", c.URL, err)
		return nil
	}
	if u.Scheme != "ws" && u.Scheme != "wss" {
		ob.Logger.Warnf("WebSocket Reverse (%v) 启动失败, URL 不合法, 必须使用 WS 或 WSS 协议", c.URL)
		return nil
	}

	conn, _, err := websocket.DefaultDialer.Dial(c.URL, nil)
	if err != nil {
		// TODO: reconnect
		ob.Logger.Warnf("WebSocket Reverse (%v) 启动失败, 错误: %v", c.URL, err)
		return nil
	}

	// protect concurrent writes to the same connection
	connWriteLock := &sync.Mutex{}
	wg := &sync.WaitGroup{}

	eventChan := ob.openEventListenChan()
	wg.Add(1)
	go func() {
		defer wg.Done()
		// keep pushing events throught the connection
		for event := range eventChan {
			ob.Logger.Debugf("通过 WebSocket Reverse (%v) 推送事件 `%v`", c.URL, event.name)
			connWriteLock.Lock()
			conn.WriteMessage(websocket.TextMessage, event.bytes) // TODO: handle err
			connWriteLock.Unlock()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			// this is the only one place we read from the connection, no need to lock
			messageType, messageBytes, err := conn.ReadMessage()
			if err != nil {
				// TODO: reconnect
				if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					ob.Logger.Infof("WebSocket Reverse (%v) 连接断开", c.URL)
				} else {
					ob.Logger.Errorf("WebSocket Reverse (%v) 连接异常断开, 错误: %v", c.URL, err)
				}
				break
			}

			isBinary := messageType == websocket.BinaryMessage
			resp := ob.parseAndHandleActionRequest(messageBytes, isBinary)
			respBytes, err := resp.encode(isBinary)
			if err != nil {
				err := fmt.Errorf("动作响应编码失败, 错误: %v", err)
				ob.Logger.Warn(err)
				respBytes, _ = failedResponse(RetCodeBadHandler, err).encode(isBinary)
			}
			connWriteLock.Lock()
			conn.WriteMessage(messageType, respBytes) // TODO: handle err
			connWriteLock.Unlock()
		}
	}()

	return func() {
		ob.closeEventListenChan(eventChan)
		// try close the connection gracefully
		err := conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""), time.Time{})
		if err != nil {
			// be rude if necessary
			conn.Close()
		}
		wg.Wait()
		ob.Logger.Infof("WebSocket Reverse (%v) 已关闭", c.URL)
	}
}
