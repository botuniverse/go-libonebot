package libonebot_test

import (
	"fmt"
	"sync/atomic"
	"time"

	libob "github.com/botuniverse/go-libonebot"
	"github.com/sirupsen/logrus"
)

func Example_1() {
	// 示例: 什么都不做的 OneBot 实现

	// 创建空 Config
	config := &libob.Config{}
	// 创建机器人自身标识
	self := &libob.Self{
		Platform: "nothing",
		UserID:   "id_of_bot",
	}
	// 创建 OneBot 实例
	ob := libob.NewOneBot("go-onebot-nothing", self, config)
	// 运行 OneBot 实例
	ob.Run()
}

var ob *libob.OneBot

func Example_2() {
	// 示例: 修改和使用 Logger
	ob.Logger.SetLevel(logrus.InfoLevel)
	ob.Logger.Infof("这是一个 INFO 日志")
}

func Example_3() {
	// 示例: 扩展 Config 和 OneBot 类型

	type MyConfig struct {
		libob.Config
		SelfID string
		UserID string
	}

	type MyOneBot struct {
		*libob.OneBot
		config *MyConfig
	}

	const Impl = "go-onebot-tg"
	const Platform = "tg"

	config := &MyConfig{ /* ... */ }
	self := &libob.Self{
		Platform: Platform,
		UserID:   config.SelfID,
	}
	ob := &MyOneBot{
		OneBot: libob.NewOneBot(Impl, self, &config.Config),
		config: config,
	}
}

const PlatformPrefix = "myplat"

func Example_4() {
	// 示例: 使用 ActionMux 注册动作处理器

	mux := libob.NewActionMux()

	// 注册 get_status 动作处理函数
	mux.HandleFunc(libob.ActionGetStatus, func(w libob.ResponseWriter, r *libob.Request) {
		w.WriteData(map[string]interface{}{
			"good":                            true,
			"online":                          true,
			PlatformPrefix + "special_status": "元气满满", // 扩展动作响应
		})
	})

	// 注册 myplat.some_action 扩展动作处理函数
	mux.HandleFunc(PlatformPrefix+".some_action", func(w libob.ResponseWriter, r *libob.Request) {
		w.WriteData("It works!") // 返回一个字符串 (返回什么都行)
	})

	// 注册 mux 为动作请求处理器
	ob.Handle(mux)
}

var mux *libob.ActionMux

func Example_5() {
	// 示例: 使用 ParamGetter 获取动作参数
	mux.HandleFunc(libob.ActionGetUserInfo, func(w libob.ResponseWriter, r *libob.Request) {
		p := libob.NewParamGetter(w, r)
		userID, ok := p.GetString("user_id") // 获取标准动作参数
		if !ok {
			return
		}
		nocache, ok := p.GetBool(PlatformPrefix + ".nocache") // 获取扩展参数
		w.WriteData(map[string]interface{}{
			"user_id":  userID,
			"nickname": userID,
		})
	})
}

var lastMessageID = uint64(0)

func Example_6() {
	// 示例: 构造并推送事件

	// 生成或获取消息 ID
	messageID := fmt.Sprint(atomic.AddUint64(&lastMessageID, 1))
	// 构造消息对象
	message := libob.Message{
		libob.MentionSegment("some_user"),
		libob.TextSegment(" 你好啊～"),
	}
	// 构造消息的替代表示
	alt_message := "@some_user 你好啊～"
	// 构造事件对象
	event := libob.MakePrivateMessageEvent(time.Now(), messageID, message, alt_message, "sender_id")
	// 推送事件
	ob.Push(&event)
}

func Example_7() {
	// 示例: 扩展标准事件

	type MyGroupMessageEvent struct {
		libob.GroupMessageEvent // 嵌入标准事件

		// 扩展字段
		Anonymous string `json:"myplat.anonymous"`
	}

	event := MyGroupMessageEvent{
		GroupMessageEvent: libob.MakeGroupMessageEvent(time.Now(), "message_id", libob.Message{}, "alt_message", "group_id", "user_id"),
		Anonymous:         "齐天大圣",
	}
	ob.Push(&event)
}

func Example_8() {
	// 示例: 多机器人账号复用 OneBot 对象

	config := &libob.Config{ /* ... */ }
	ob := libob.NewOneBotMultiSelf("go-onebot-multi", config)

	mux := libob.NewActionMux()
	mux.HandleFunc(libob.ActionSendMessage, func(w libob.ResponseWriter, r *libob.Request) {
		if r.Self == nil {
			w.WriteFailed(libob.RetCodeWhoAmI, fmt.Errorf("未指定机器人账号"))
			return
		}
		// 通过 r.Self 获得用户指定的机器人自身标识
		_ = r.Self.Platform
		_ = r.Self.UserID
	})

	ob.Handle(mux)
	go ob.Run()

	event1 := libob.MakeFriendIncreaseNoticeEvent(time.Now(), "friend_id")
	ob.PushWithSelf(&event1, &libob.Self{Platform: "myplat1", UserID: "bot_id_1"})
	event2 := libob.MakeFriendIncreaseNoticeEvent(time.Now(), "friend_id")
	ob.PushWithSelf(&event2, &libob.Self{Platform: "myplat1", UserID: "bot_id_2"})
}
