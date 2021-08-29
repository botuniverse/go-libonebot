package libonebot

import (
	"sync"

	log "github.com/sirupsen/logrus"
)

type OneBot struct {
	Platform string
	Config   *Config

	eventListenChans     []chan marshaledEvent
	eventListenChansLock sync.RWMutex

	handlers         map[string]Handler
	extendedHandlers map[string]Handler

	commClosers     []commCloser
	commClosersLock sync.Mutex
	wg              sync.WaitGroup
}

func NewOneBot(platform string, config *Config) *OneBot {
	return &OneBot{
		Platform: platform,
		Config:   config,

		eventListenChans:     make([]chan marshaledEvent, 0),
		eventListenChansLock: sync.RWMutex{},

		handlers:         make(map[string]Handler),
		extendedHandlers: make(map[string]Handler),

		commClosers: make([]commCloser, 0),
		wg:          sync.WaitGroup{},
	}
}

func (ob *OneBot) Run() {
	if ob.Platform == "" {
		log.Errorf("OneBot 无法启动, 没有提供 OneBot 平台名称")
		return
	}
	if ob.Config == nil {
		log.Errorf("OneBot 无法启动, 没有提供 OneBot 配置")
		return
	}
	ob.startCommMethods()
	log.Infof("OneBot 已启动")
	ob.wg.Add(1)
	ob.wg.Wait()
	log.Infof("OneBot 已关闭")
}

func (ob *OneBot) Shutdown() {
	ob.commClosersLock.Lock()
	for _, closer := range ob.commClosers {
		closer.Close()
	}
	ob.commClosers = make([]commCloser, 0)
	ob.commClosersLock.Unlock()
	ob.wg.Done()
}
