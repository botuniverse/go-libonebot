package libonebot

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/tevino/abool/v2"
)

type wsComm struct {
	ob          *OneBot
	addr        string
	accessToken string
}

var wsUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (comm *wsComm) handle(w http.ResponseWriter, r *http.Request) {
	comm.ob.Logger.Debugf("收到来自 %v 的 WebSocket (%v) 连接请求", r.RemoteAddr, comm.addr)

	// authorization
	if comm.accessToken != "" {
		if r.Header.Get("Authorization") != "Bearer "+comm.accessToken {
			comm.ob.Logger.Errorf("请求头中的 Authorization 不匹配")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}

	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		comm.ob.Logger.Errorf("WebSocket (%v) 连接失败, 错误: %v", comm.addr, err)
		return
	}
	comm.ob.Logger.Infof("WebSocket (%v) 连接成功", comm.addr)
	defer conn.Close()

	// protect concurrent writes to the same connection
	connWriteLock := &sync.Mutex{}

	isClosed := abool.New()
	checkError := func(err error) bool {
		if err != nil {
			if isClosed.IsNotSet() {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					comm.ob.Logger.Infof("WebSocket (%v) 连接断开", comm.addr)
				} else {
					comm.ob.Logger.Errorf("WebSocket (%v) 连接异常断开, 错误: %v", comm.addr, err)
				}
			}
			isClosed.Set()
			return true
		}
		return false
	}

	eventChan := comm.ob.openEventListenChan()
	defer comm.ob.closeEventListenChan(eventChan)
	go func() {
		// keep pushing events throught the connection
		for event := range eventChan {
			comm.ob.Logger.Debugf("通过 WebSocket (%v) 推送事件 `%v`", comm.addr, event.name)
			connWriteLock.Lock()
			err := conn.WriteMessage(websocket.TextMessage, event.bytes)
			connWriteLock.Unlock()
			if checkError(err) {
				break
			}
		}
	}()

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
}

func commRunWS(c ConfigCommWS, ob *OneBot, ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	addr := fmt.Sprintf("%s:%d", c.Host, c.Port)
	ob.Logger.Infof("正在启动 WebSocket (%v)...", addr)

	comm := &wsComm{
		ob:          ob,
		addr:        addr,
		accessToken: c.AccessToken,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", comm.handle)
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			ob.Logger.Errorf("WebSocket (%v) 启动失败, 错误: %v", addr, err)
		}
	}()

	<-ctx.Done()
	if err := server.Shutdown(context.TODO()); err != nil {
		ob.Logger.Errorf("WebSocket (%v) 关闭失败, 错误: %v", addr, err)
	} else {
		ob.Logger.Infof("WebSocket (%v) 已关闭", addr)
	}
}
