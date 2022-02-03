package libonebot

import (
	"context"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/tevino/abool/v2"
)

type wsReverseComm struct {
	wsCommCommon
	config            ConfigCommWSReverse
	url               string
	accessToken       string
	reconnectInterval time.Duration
	isShutdown        *abool.AtomicBool
}

func (comm *wsReverseComm) connectAndServe(ctx context.Context) {
	comm.ob.Logger.Debugf("WebSocket Reverse (%v) 开始连接", comm.url)

	header := http.Header{}
	if comm.accessToken != "" {
		header.Set("Authorization", "Bearer "+comm.accessToken)
	}
	header.Set("User-Agent", comm.ob.GetUserAgent())
	header.Set("X-OneBot-Version", OneBotVersion)
	header.Set("X-Impl", comm.ob.Impl)
	header.Set("X-Platform", comm.ob.Platform)
	header.Set("X-Self-ID", comm.ob.SelfID)
	conn, _, err := websocket.DefaultDialer.Dial(comm.url, header)
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
			go comm.handleRequest(conn, connWriteLock, messageBytes, messageType, RequestComm{
				Method: CommMethodWSReverse,
				Config: comm.config,
			})
		}
	}()

	eventChan := comm.ob.openEventListenChan()
	defer comm.ob.closeEventListenChan(eventChan)

loop:
	for {
		select {
		case event := <-eventChan:
			comm.ob.Logger.Debugf("通过 WebSocket Reverse (%v) 推送事件 `%v`", comm.url, event.name)
			go comm.pushEvent(conn, connWriteLock, event)
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
		wsCommCommon:      wsCommCommon{ob: ob},
		config:            c,
		url:               c.URL,
		accessToken:       c.AccessToken,
		reconnectInterval: time.Duration(c.ReconnectInterval) * time.Millisecond,
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
		ob.Logger.Infof("WebSocket Reverse (%v) 将在 %v 秒后尝试重连", comm.url, c.ReconnectInterval)
		time.Sleep(comm.reconnectInterval)
	}
	ob.Logger.Infof("WebSocket Reverse (%v) 已关闭", comm.url)
}
