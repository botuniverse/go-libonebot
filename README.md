# (WIP) Go LibOneBot

[![OneBot](https://img.shields.io/badge/OneBot-v12-black?logo=data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAHAAAABwCAMAAADxPgR5AAAAGXRFWHRTb2Z0d2FyZQBBZG9iZSBJbWFnZVJlYWR5ccllPAAAAAxQTFRF////29vbr6+vAAAAk1hCcwAAAAR0Uk5T////AEAqqfQAAAKcSURBVHja7NrbctswDATQXfD//zlpO7FlmwAWIOnOtNaTM5JwDMa8E+PNFz7g3waJ24fviyDPgfhz8fHP39cBcBL9KoJbQUxjA2iYqHL3FAnvzhL4GtVNUcoSZe6eSHizBcK5LL7dBr2AUZlev1ARRHCljzRALIEog6H3U6bCIyqIZdAT0eBuJYaGiJaHSjmkYIZd+qSGWAQnIaz2OArVnX6vrItQvbhZJtVGB5qX9wKqCMkb9W7aexfCO/rwQRBzsDIsYx4AOz0nhAtWu7bqkEQBO0Pr+Ftjt5fFCUEbm0Sbgdu8WSgJ5NgH2iu46R/o1UcBXJsFusWF/QUaz3RwJMEgngfaGGdSxJkE/Yg4lOBryBiMwvAhZrVMUUvwqU7F05b5WLaUIN4M4hRocQQRnEedgsn7TZB3UCpRrIJwQfqvGwsg18EnI2uSVNC8t+0QmMXogvbPg/xk+Mnw/6kW/rraUlvqgmFreAA09xW5t0AFlHrQZ3CsgvZm0FbHNKyBmheBKIF2cCA8A600aHPmFtRB1XvMsJAiza7LpPog0UJwccKdzw8rdf8MyN2ePYF896LC5hTzdZqxb6VNXInaupARLDNBWgI8spq4T0Qb5H4vWfPmHo8OyB1ito+AysNNz0oglj1U955sjUN9d41LnrX2D/u7eRwxyOaOpfyevCWbTgDEoilsOnu7zsKhjRCsnD/QzhdkYLBLXjiK4f3UWmcx2M7PO21CKVTH84638NTplt6JIQH0ZwCNuiWAfvuLhdrcOYPVO9eW3A67l7hZtgaY9GZo9AFc6cryjoeFBIWeU+npnk/nLE0OxCHL1eQsc1IciehjpJv5mqCsjeopaH6r15/MrxNnVhu7tmcslay2gO2Z1QfcfX0JMACG41/u0RrI9QAAAABJRU5ErkJggg==)](https://github.com/botuniverse/onebot/pull/108)
[![Go Reference](https://pkg.go.dev/badge/github.com/botuniverse/go-libonebot.svg)](https://pkg.go.dev/github.com/botuniverse/go-libonebot)

> 目前大体 API 已经成型，但还有一些细节尚未完成……

Go LibOneBot 可以帮助 OneBot 实现者快速在新的聊天机器人平台实现 OneBot v12 接口标准。

具体而言，Go LibOneBot 通过 `OneBot`、`Config`、`ActionMux`、`Event` 等类型的抽象，让 OneBot 实现者只需编写少量代码即可完成一个 OneBot 实现，而无需关心过多 OneBot 标准所定义的通信方式的细节。

基于 LibOneBot 实现 OneBot 时，OneBot 实现者只需专注于编写与聊天机器人平台对接的逻辑，包括通过长轮询或 webhook 方式从机器人平台获得事件，并将其转换为 OneBot 事件，以及处理 OneBot 动作请求，并将其转换为对机器人平台 API 的调用。

## 用法

一个什么都不做的 OneBot 实现：

```go
package main

import (
    libob "github.com/botuniverse/go-libonebot"
)

func main() {
    config := &libob.Config{} // 创建空 Config
    ob := libob.NewOneBot("nothing", config) // 创建 OneBot 实例
    ob.HandleFunc(func(w libob.ResponseWriter, r *libob.Request) {
        // 对所有动作请求都返回 OK
        w.WriteOK()
    })
    ob.Run() // 运行 OneBot 实例
}
```

通过交互命令行输入“私聊消息”的实现：

```go
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync/atomic"

	libob "github.com/botuniverse/go-libonebot"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type OneBotREPL struct {
	*libob.OneBot // 嵌入 OneBot 对象
	config        *REPLConfig
	lastMessageID uint64
}

type Config struct {
	OneBot libob.Config `mapstructure:",squash"` // 嵌入 LibOneBot 配置
	REPL   REPLConfig
}

type REPLConfig struct {
	SelfID string `mapstructure:"self_id"`
	UserID string `mapstructure:"user_id"`
}

const defaultConfigString = `
[heartbeat]
enabled = true
interval = 10

[auth]
access_token = ""

[repl]
self_id = "bot"
user_id = "user"
`

func loadConfig() *Config {
	// 使用 viper 库加载配置
	v := viper.New()
	v.SetConfigType("toml")
	v.ReadConfig(strings.NewReader(defaultConfigString)) // 加载默认配置
	v.SetConfigFile("config.toml")
	err := v.MergeInConfig() // 合并配置文件内容
	if err != nil && os.IsNotExist(err) {
		fmt.Println("配置文件不存在, 正在写入默认配置到 config.toml")
		v.WriteConfigAs("config.toml")
	}
	config := &Config{}
	v.Unmarshal(config)
	fmt.Printf("配置加载成功: %+v\n", config)
	return config
}

func main() {
	// 加载配置
	config := loadConfig()

	// 创建 OneBot 实例
	ob := &OneBotREPL{
		OneBot:        libob.NewOneBot("repl" /* 聊天平台名称，用作扩展动作名等的前缀 */, &config.OneBot),
		config:        &config.REPL,
		lastMessageID: 0,
	}

	// 修改日志配置
	logFile, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	ob.Logger.SetOutput(logFile)
	ob.Logger.SetLevel(logrus.InfoLevel)

	// 通过 ActionMux 注册动作处理函数，该 mux 变量可在多个 OneBot 实例复用
	mux := libob.NewActionMux()
	// 注册 get_version 动作处理函数
	mux.HandleFunc(libob.ActionGetVersion, func(w libob.ResponseWriter, r *libob.Request) {
		// 返回一个映射类型的数据（序列化为 JSON 对象或 MsgPack 映射）
		w.WriteData(map[string]string{
			"version": "1.0.0",
		})
	})
	// 注册 get_self_id 动作处理函数
	mux.HandleFunc(libob.ActionGetSelfInfo, func(w libob.ResponseWriter, r *libob.Request) {
		w.WriteData(map[string]interface{}{
			"user_id": ob.config.SelfID, // 返回配置中指定的 self_id
		})
	})
	mux.HandleFunc(libob.ActionSendMessage, func(w libob.ResponseWriter, r *libob.Request) {
		// 创建 ParamGetter 来获取参数，也可以直接用 r.Params.GetXxx
		p := libob.NewParamGetter(w, r)
		userID, ok := p.GetString("user_id")
		if !ok {
			return
		}
		if userID != ob.config.UserID {
			// user_id 不匹配，返回 RetCodeLogicError
			w.WriteFailed(libob.RetCodeLogicError, fmt.Errorf("无法发送给用户 `%v`", userID))
			return
		}
		msg, ok := p.GetMessage("message")
		if !ok {
			return
		}
		fmt.Println(msg.ExtractText()) // 提取消息中的纯文本并打印在控制台
		// 返回消息 ID
		w.WriteData(map[string]interface{}{
			"message_id": fmt.Sprint(atomic.AddUint64(&ob.lastMessageID, 1)),
		})
	})
	mux.HandleFuncExtended("test", func(w libob.ResponseWriter, r *libob.Request) {
		// 该扩展动作通过 repl_test 动作名来调用
		w.WriteData("It works!") // 返回一个字符串
	})

	ob.Handle(mux)
	go ob.Run() // 启动 OneBot 实例

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("请开始对话 (输入 exit 退出):")
	// 循环读取命令行输入
	for {
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		if text == "exit" {
			ob.Shutdown()
			break
		}
		// 构造 OneBot 私聊消息事件并通过 OneBot 对象推送到机器人业务端
		go ob.Push(&libob.MessageEvent{
			Event: libob.Event{
				SelfID:     ob.config.SelfID,
				Type:       libob.EventTypeMessage,
				DetailType: "private",
			},
			UserID:  ob.config.UserID,
			Message: libob.Message{libob.TextSegment(text)},
		})
	}
}
```

关于上面示例中所涉及的类型、函数的更多细节，请参考 [Godoc 文档](https://pkg.go.dev/github.com/botuniverse/go-libonebot)。
