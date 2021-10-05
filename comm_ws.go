package libonebot

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type wsComm struct {
	ob   *OneBot
	addr string
}

var wsUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (comm *wsComm) handle(w http.ResponseWriter, r *http.Request) {
	comm.ob.Logger.Infof("收到来自 %v 的 WebSocket (%v) 连接请求", r.RemoteAddr, comm.addr)
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		comm.ob.Logger.Errorf("WebSocket (%v) 连接失败, 错误: %v", comm.addr, err)
		return
	}
	comm.ob.Logger.Infof("WebSocket (%v) 连接成功", comm.addr)
	defer conn.Close()

	// protect concurrent writes to the same connection
	connWriteLock := &sync.Mutex{}

	eventChan := comm.ob.openEventListenChan()
	defer comm.ob.closeEventListenChan(eventChan)
	go func() {
		// keep pushing events throught the connection
		for event := range eventChan {
			comm.ob.Logger.Debugf("通过 WebSocket (%v) 推送事件 `%v`", comm.addr, event.name)
			connWriteLock.Lock()
			conn.WriteMessage(websocket.TextMessage, event.bytes) // TODO: handle err
			connWriteLock.Unlock()
		}
	}()

	for {
		// this is the only one place we read from the connection, no need to lock
		messageType, messageBytes, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				comm.ob.Logger.Infof("WebSocket (%v) 连接断开", comm.addr)
			} else {
				comm.ob.Logger.Errorf("WebSocket (%v) 连接异常断开, 错误: %v", comm.addr, err)
			}
			break
		}

		isBinary := messageType == websocket.BinaryMessage
		resp := comm.ob.parseAndHandleActionRequest(messageBytes, isBinary)
		respBytes, err := resp.encode(isBinary)
		if err != nil {
			err := fmt.Errorf("动作响应编码失败, 错误: %v", err)
			comm.ob.Logger.Warn(err)
			respBytes, _ = failedResponse(RetCodeBadHandler, err).encode(isBinary)
		}
		connWriteLock.Lock()
		conn.WriteMessage(messageType, respBytes) // TODO: handle err
		connWriteLock.Unlock()
	}
}

func commStartWS(c ConfigCommWS, ob *OneBot) commCloser {
	addr := fmt.Sprintf("%s:%d", c.Host, c.Port)
	ob.Logger.Infof("正在启动 WebSocket (%v)...", addr)

	comm := &wsComm{ob: ob, addr: addr}
	mux := http.NewServeMux()
	mux.HandleFunc("/", comm.handle)
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			ob.Logger.Errorf("WebSocket (%v) 启动失败, 错误: %v", addr, err)
		} else {
			ob.Logger.Infof("WebSocket (%v) 已关闭", addr)
		}
	}()

	return func() {
		if err := server.Shutdown(context.TODO() /* TODO */); err != nil {
			ob.Logger.Errorf("WebSocket (%v) 关闭失败, 错误: %v", addr, err)
		}
		// TODO: wg.Wait() 后再输出已关闭
	}
}
