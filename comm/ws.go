package comm

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type wsComm struct {
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (comm *wsComm) handle(w http.ResponseWriter, r *http.Request) {
	log.Infof("WebSocket 请求: %v", r.URL)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Errorf("WebSocket 连接失败: %v", err)
		return
	}
	log.Infof("WebSocket 连接成功")
	defer conn.Close()
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Errorf("WebSocket 接收错误: %v", err)
			break
		}
		log.Infof("WebSocket 接收到消息: %v", string(message))
	}
}

// Start a WebSocket commmunication task.
func StartWSTask(host string, port uint16) {
	addr := fmt.Sprintf("%s:%d", host, port)
	log.Infof("正在启动 WebSocket 通信方式, 监听地址: %v", addr)

	comm := &wsComm{}
	mux := http.NewServeMux()
	mux.HandleFunc("/", comm.handle)

	go func() {
		if err := http.ListenAndServe(addr, mux); err != nil && err != http.ErrServerClosed {
			log.Errorf("WebSocket 通信方式启动失败, 错误: %v", err)
		}
		log.Info("WebSocket 通信方式已退出")
	}()
}
