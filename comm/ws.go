package comm

import (
	"fmt"
	"net/http"

	"github.com/botuniverse/go-libonebot/utils"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

type wsComm struct {
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (comm *wsComm) handle(w http.ResponseWriter, r *http.Request) {
	log.Infof("收到来自 %v 的 WebSocket 连接请求", r.RemoteAddr)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Errorf("WebSocket 连接失败, 错误: %v", err)
		return
	}
	log.Infof("WebSocket 连接成功")
	defer conn.Close()

	ws_chan := make(chan []byte) // TODO: channel size
	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					log.Infof("WebSocket 连接断开")
				} else {
					log.Errorf("WebSocket 连接异常断开, 错误: %v", err)
				}
				break
			}
			ws_chan <- message
		}
	}()

	for {
		select {
		case messageBytes := <-ws_chan:
			message := utils.BytesToString(messageBytes)
			log.Debugf("WebSocket message: %v", message)
			if !gjson.Valid(message) {
				log.Warnf("Action 请求体不是合法的 JSON, 已忽略")
				continue
			}
			actionRequest := gjson.Parse(message)
			actionResponse := handleAction(actionRequest)
			conn.WriteMessage(websocket.TextMessage, utils.StringToBytes(actionResponse.String()))
		}
	}
}

// Start a WebSocket commmunication task.
func StartWSTask(host string, port uint16) {
	addr := fmt.Sprintf("%s:%d", host, port)
	log.Infof("正在启动 WebSocket (%v)...", addr)

	comm := &wsComm{}
	mux := http.NewServeMux()
	mux.HandleFunc("/", comm.handle)

	go func() {
		if err := http.ListenAndServe(addr, mux); err != nil && err != http.ErrServerClosed {
			log.Errorf("WebSocket (%v) 启动失败, 错误: %v", addr, err)
		}
		log.Infof("WebSocket (%v) 已关闭", addr)
	}()
}
