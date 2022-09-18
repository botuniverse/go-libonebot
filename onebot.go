package libonebot

import (
	"context"
	"fmt"
	"regexp"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	Version       = "0.6.0" // LibOneBot 版本号
	OneBotVersion = "12"    // OneBot 标准版本
)

// OneBot 表示一个 OneBot 实例.
type OneBot struct {
	Impl   string
	Self   *Self // 机器人自身标识, 多机器人账号复用 OneBot 对象时为 nil
	Config *Config
	Logger *logrus.Logger

	eventListenChans     []chan marshaledEvent
	eventListenChansLock *sync.RWMutex

	actionHandler Handler

	cancel context.CancelFunc
	wg     *sync.WaitGroup
}

var (
	implPlatformRegex = regexp.MustCompile(`^[a-z][\-a-z0-9]*(\.[\-a-z0-9]+)*$`)
)

// NewOneBot 创建一个新的 OneBot 实例.
//
// 参数:
//   impl: OneBot 实现名称, 不能为空
//   self: OneBot 实例对应的机器人自身标识, 不能为 nil
//   config: OneBot 配置, 不能为 nil
func NewOneBot(impl string, self *Self, config *Config) *OneBot {
	if impl == "" {
		panic("必须提供 OneBot 实现名称")
	}
	if !implPlatformRegex.MatchString(impl) {
		panic("OneBot 实现名称不合法")
	}
	if self == nil {
		panic("必须提供机器人自身标识")
	}
	if self.Platform == "" {
		panic("必须提供机器人平台名称")
	}
	if !implPlatformRegex.MatchString(self.Platform) {
		panic("机器人平台名称不合法")
	}
	if self.UserID == "" {
		panic("必须提供 OneBot 实例对应的机器人用户 ID")
	}
	if config == nil {
		panic("必须提供 OneBot 配置")
	}
	return newOneBotUnchecked(impl, self, config)
}

// NewOneBotMultiSelf 创建一个新的多机器人账号复用的 OneBot 实例.
//
// 参数:
//   impl: OneBot 实现名称, 不能为空
//   config: OneBot 配置, 不能为 nil
func NewOneBotMultiSelf(impl string, config *Config) *OneBot {
	if impl == "" {
		panic("必须提供 OneBot 实现名称")
	}
	if !implPlatformRegex.MatchString(impl) {
		panic("OneBot 实现名称不合法")
	}
	if config == nil {
		panic("必须提供 OneBot 配置")
	}
	return newOneBotUnchecked(impl, nil, config)
}

func newOneBotUnchecked(impl string, self *Self, config *Config) *OneBot {
	return &OneBot{
		Impl:   impl,
		Self:   self,
		Config: config,
		Logger: logrus.New(),

		eventListenChans:     make([]chan marshaledEvent, 0),
		eventListenChansLock: &sync.RWMutex{},

		actionHandler: nil,

		cancel: nil,
		wg:     &sync.WaitGroup{},
	}
}

// Run 运行 OneBot 实例.
//
// 该方法会阻塞当前线程, 直到 Shutdown 被调用.
func (ob *OneBot) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	ob.cancel = cancel

	ob.startCommMethods(ctx)
	ob.startHeartbeat(ctx)

	ob.Logger.Infof("OneBot 已启动")
	<-ctx.Done()
}

// Shutdown 停止 OneBot 实例.
func (ob *OneBot) Shutdown() {
	ob.cancel()  // this will stop everything (comm methods, heartbeat, etc)
	ob.wg.Wait() // wait for everything to completely stop
	ob.Logger.Infof("OneBot 已关闭")
}

// GetUserAgent 获取 OneBot 实例的 User-Agent.
func (ob *OneBot) GetUserAgent() string {
	return fmt.Sprintf("OneBot/%v LibOneBot/%v", OneBotVersion, Version)
}

func (ob *OneBot) startCommMethods(ctx context.Context) {
	if ob.Config.Comm.HTTP != nil {
		for _, c := range ob.Config.Comm.HTTP {
			ob.wg.Add(1)
			go commRunHTTP(c, ob, ctx, ob.wg)
		}
	}

	if ob.Config.Comm.HTTPWebhook != nil {
		for _, c := range ob.Config.Comm.HTTPWebhook {
			ob.wg.Add(1)
			go commRunHTTPWebhook(c, ob, ctx, ob.wg)
		}
	}

	if ob.Config.Comm.WS != nil {
		for _, c := range ob.Config.Comm.WS {
			ob.wg.Add(1)
			go commRunWS(c, ob, ctx, ob.wg)
		}
	}

	if ob.Config.Comm.WSReverse != nil {
		for _, c := range ob.Config.Comm.WSReverse {
			ob.wg.Add(1)
			go commRunWSReverse(c, ob, ctx, ob.wg)
		}
	}
}

func (ob *OneBot) startHeartbeat(ctx context.Context) {
	if !ob.Config.Heartbeat.Enabled {
		return
	}

	if ob.Config.Heartbeat.Interval == 0 {
		ob.Logger.Errorf("心跳间隔必须大于 0")
		return
	}

	ob.wg.Add(1)
	go func() {
		defer ob.wg.Done()

		ticker := time.NewTicker(time.Duration(ob.Config.Heartbeat.Interval) * time.Millisecond)
		defer ticker.Stop()

		ob.Logger.Infof("心跳开始")
		for {
			select {
			case <-ticker.C:
				ob.Logger.Debugf("扑通")
				event := MakeHeartbeatMetaEvent(time.Now(), int64(ob.Config.Heartbeat.Interval))
				ob.Push(&event)
			case <-ctx.Done():
				ob.Logger.Infof("心跳停止")
				return
			}
		}
	}()
}
