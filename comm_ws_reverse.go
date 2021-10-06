package libonebot

import (
	"context"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/tevino/abool/v2"
)

type wsReverseComm struct {
	ob                *OneBot
	url               string
	reconnectInterval time.Duration
	isShutdown        *abool.AtomicBool
}

func (comm *wsReverseComm) connectAndServe(ctx context.Context) {
	conn, _, err := websocket.DefaultDialer.Dial(comm.url, nil)
	if err != nil {
		comm.ob.Logger.Errorf("WebSocket Reverse (%v) 连接失败, 错误: %v", comm.url, err)
		return
	}
	comm.ob.Logger.Infof("WebSocket Reverse (%v) 连接成功", comm.url)

	// protect concurrent writes to the same connection
	connWriteLock := &sync.Mutex{}

	connCtx, connCancel := context.WithCancel(context.Background())
	isClosed := abool.New()
	checkError := func(err error) bool {
		if err != nil {
			if isClosed.IsNotSet() {
				connCancel() // this will be called for only one time
				if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					comm.ob.Logger.Infof("WebSocket Reverse (%v) 连接断开", comm.url)
				} else {
					comm.ob.Logger.Errorf("WebSocket Reverse (%v) 连接异常断开, 错误: %v", comm.url, err)
				}
			}
			isClosed.Set()
			return true
		}
		return false
	}

	wsClientWG := &sync.WaitGroup{}
	wsClientWG.Add(1)
	go func() {
		defer wsClientWG.Done()
		for {
			// this is the only one place we read from the connection, no need to lock
			messageType, messageBytes, err := conn.ReadMessage()
			if checkError(err) {
				break
			}
			isBinary := messageType == websocket.BinaryMessage
			resp := comm.ob.decodeAndHandleRequest(messageBytes, isBinary)
			respBytes, _ := comm.ob.encodeResponse(resp, isBinary)
			connWriteLock.Lock()
			err = conn.WriteMessage(messageType, respBytes)
			connWriteLock.Unlock()
			if checkError(err) {
				break
			}
		}
	}()

	eventChan := comm.ob.openEventListenChan()
	defer comm.ob.closeEventListenChan(eventChan)

loop:
	for {
		select {
		case event := <-eventChan:
			comm.ob.Logger.Debugf("通过 WebSocket Reverse (%v) 推送事件 `%v`", comm.url, event.name)
			connWriteLock.Lock()
			err := conn.WriteMessage(websocket.TextMessage, event.bytes)
			connWriteLock.Unlock()
			if checkError(err) {
				break loop
			}
		case <-connCtx.Done(): // connection closed
			break loop
		case <-ctx.Done(): // onebot shutdown
			// try close the connection gracefully
			err = conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""), time.Time{})
			if err != nil {
				// be rude if necessary
				conn.Close()
			}
			comm.isShutdown.Set()
			break loop
		}
	}

	wsClientWG.Wait() // wait the ws client goroutine to finish
}

func commRunWSReverse(c ConfigCommWSReverse, ob *OneBot, ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	ob.Logger.Infof("正在启动 WebSocket Reverse (%v)...", c.URL)

	u, err := url.Parse(c.URL)
	if err != nil {
		ob.Logger.Errorf("WebSocket Reverse (%v) 启动失败, URL 不合法, 错误: %v", c.URL, err)
		return
	}
	if u.Scheme != "ws" && u.Scheme != "wss" {
		ob.Logger.Errorf("WebSocket Reverse (%v) 启动失败, URL 不合法, 必须使用 WS 或 WSS 协议", c.URL)
		return
	}

	if c.ReconnectInterval == 0 {
		ob.Logger.Errorf("WebSocket Reverse 重连间隔必须大于 0")
		return
	}

	comm := wsReverseComm{
		ob:                ob,
		url:               c.URL,
		reconnectInterval: time.Duration(c.ReconnectInterval) * time.Second,
		isShutdown:        abool.New(),
	}

	go func() {
		<-ctx.Done()
		comm.isShutdown.Set()
	}()

	for {
		comm.connectAndServe(ctx)
		if comm.isShutdown.IsSet() {
			break
		}
		time.Sleep(comm.reconnectInterval)
		ob.Logger.Infof("WebSocket Reverse (%v) 尝试重连", comm.url)
	}
	ob.Logger.Infof("WebSocket Reverse (%v) 已关闭", comm.url)
}
