package libonebot

import (
	"context"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// OneBot 表示一个 OneBot 实例.
type OneBot struct {
	Platform string
	SelfID   string
	Config   *Config
	Logger   *logrus.Logger

	eventListenChans     []chan marshaledEvent
	eventListenChansLock sync.RWMutex

	actionHandler Handler

	commClosers     []commCloser
	commClosersLock sync.Mutex

	cancel context.CancelFunc
}

// NewOneBot 创建一个新的 OneBot 实例.
//
// 参数:
//   platform: OneBot 实现平台名称, 应和扩展动作名称、扩展参数等前缀相同, 不能为空
//   selfID: OneBot 实例对应的机器人自身 ID, 不能为空
//   config: OneBot 配置, 不能为 nil
func NewOneBot(platform string, selfID string, config *Config) *OneBot {
	if platform == "" {
		panic("必须提供 OneBot 平台名称")
	}
	if selfID == "" {
		panic("必须提供 OneBot 实例对应的机器人自身 ID")
	}
	if config == nil {
		panic("必须提供 OneBot 配置")
	}
	return &OneBot{
		Platform: platform,
		SelfID:   selfID,
		Config:   config,
		Logger:   logrus.New(),

		eventListenChans:     make([]chan marshaledEvent, 0),
		eventListenChansLock: sync.RWMutex{},

		actionHandler: nil,

		commClosers:     make([]commCloser, 0),
		commClosersLock: sync.Mutex{},

		cancel: nil,
	}
}

// Run 运行 OneBot 实例.
//
// 该方法会阻塞当前线程, 直到 Shutdown 被调用.
func (ob *OneBot) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	ob.cancel = cancel

	ob.startCommMethods()
	if ob.Config.Heartbeat.Enabled {
		go ob.heartbeat(ctx)
	}

	ob.Logger.Infof("OneBot 已启动")
	<-ctx.Done()
	ob.Logger.Infof("OneBot 已关闭")
}

// Shutdown 停止 OneBot 实例.
func (ob *OneBot) Shutdown() {
	ob.commClosersLock.Lock()
	for _, closer := range ob.commClosers {
		closer.Close()
	}
	ob.commClosers = make([]commCloser, 0)
	ob.commClosersLock.Unlock()
	ob.cancel()
}

func (ob *OneBot) heartbeat(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(ob.Config.Heartbeat.Interval) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			event := MakeHeartbeatMetaEvent()
			ob.Push(&event)
		}
	}
}
