package libonebot

import (
	"sync"

	"github.com/sirupsen/logrus"
)

// OneBot 表示一个 OneBot 实例.
type OneBot struct {
	Platform string
	Config   *Config
	Logger   *logrus.Logger

	eventListenChans     []chan marshaledEvent
	eventListenChansLock sync.RWMutex

	actionHandler Handler

	commClosers     []commCloser
	commClosersLock sync.Mutex
	wg              sync.WaitGroup
}

// NewOneBot 创建一个新的 OneBot 实例.
//
// 参数:
//   platform: OneBot 实现平台名称, 用作动作名称等的前缀, 不能为空
//   config: OneBot 配置, 不能为 nil
func NewOneBot(platform string, config *Config) *OneBot {
	if platform == "" {
		panic("必须提供 OneBot 平台名称")
	}
	if config == nil {
		panic("必须提供 OneBot 配置")
	}
	return &OneBot{
		Platform: platform,
		Config:   config,
		Logger:   logrus.New(),

		eventListenChans:     make([]chan marshaledEvent, 0),
		eventListenChansLock: sync.RWMutex{},

		actionHandler: nil,

		commClosers: make([]commCloser, 0),
		wg:          sync.WaitGroup{},
	}
}

// Run 运行 OneBot 实例.
//
// 该方法会阻塞当前线程, 直到 Shutdown 被调用.
func (ob *OneBot) Run() {
	ob.startCommMethods()
	ob.Logger.Infof("OneBot 已启动")
	ob.wg.Add(1)
	ob.wg.Wait()
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
	ob.wg.Done()
}
