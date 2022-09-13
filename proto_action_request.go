// OneBot Connect - 数据协议 - 动作请求
// https://12.onebot.dev/connect/data-protocol/action-request/

package libonebot

import (
	"errors"

	"github.com/botuniverse/go-libonebot/utils"
	"github.com/tidwall/gjson"
	"github.com/vmihailenco/msgpack/v5"
)

const (
	CommMethodNone        = 0 // 无通信方式 (OneBot 实现内部构造的请求)
	CommMethodHTTP        = 1 // HTTP 通信方式
	CommMethodHTTPWebhook = 2 // HTTP Webhook 通信方式
	CommMethodWS          = 3 // WebSocket 通信方式
	CommMethodWSReverse   = 4 // 反向 WebSocket 通信方式
)

// RequestCommMethod 表示接收动作请求的通信方式.
type RequestComm struct {
	Method int         // 通信方式
	Config interface{} // 通信方式配置
}

// Request 表示一个动作请求.
type Request struct {
	Comm   RequestComm // 接收动作请求的通信方式
	Action string      // 动作名称
	Params EasierMap   // 动作参数
	Echo   string      // 动作请求的 echo 字段
}

func validateRequestMap(m EasierMap) error {
	if action, err := m.GetString("action"); err != nil {
		return errors.New("`action` 字段不存在或类型错误")
	} else if action == "" {
		return errors.New("`action` 字段为空")
	}
	if _, err := m.GetMap("params"); err != nil {
		return errors.New("`params` 字段不存在或类型错误")
	}
	_, err1 := m.Get("echo")
	_, err2 := m.GetString("echo")
	if err1 == nil && err2 != nil {
		return errors.New("`echo` 字段类型错误")
	}
	return nil
}

func parseRequestFromMap(m map[string]interface{}, reqComm RequestComm) (Request, error) {
	em := EasierMapFromMap(m)
	err := validateRequestMap(em)
	if err != nil {
		return Request{}, err
	}

	action, _ := em.GetString("action")
	params, _ := em.GetMap("params")
	echo, _ := em.GetString("echo")
	r := Request{
		Comm:   reqComm,
		Action: action,
		Params: params,
		Echo:   echo,
	}
	return r, nil
}

func decodeRequest(actionBytes []byte, isBinary bool, reqComm RequestComm) (Request, error) {
	var actionRequestMap map[string]interface{}
	if isBinary {
		err := msgpack.Unmarshal(actionBytes, &actionRequestMap)
		if err != nil || actionRequestMap == nil {
			return Request{}, errors.New("不是一个 MsgPack 映射")
		}
	} else {
		if !gjson.ValidBytes(actionBytes) {
			return Request{}, errors.New("不是合法的 JSON")
		}
		m, ok := gjson.Parse(utils.BytesToString(actionBytes)).Value().(map[string]interface{})
		if !ok || m == nil {
			return Request{}, errors.New("不是一个 JSON 对象")
		}
		actionRequestMap = m
	}
	return parseRequestFromMap(actionRequestMap, reqComm)
}

func decodeRequestList(actionBytes []byte, isBinary bool, reqComm RequestComm) ([]Request, error) {
	var actionRequestMapList []map[string]interface{}
	if isBinary {
		err := msgpack.Unmarshal(actionBytes, &actionRequestMapList)
		if err != nil || actionRequestMapList == nil {
			return nil, errors.New("不是一个 MsgPack 映射数组")
		}
	} else {
		if !gjson.ValidBytes(actionBytes) {
			return nil, errors.New("不是合法的 JSON")
		}
		m, ok := gjson.Parse(utils.BytesToString(actionBytes)).Value().([]interface{})
		if !ok || m == nil {
			return nil, errors.New("不是一个 JSON 数组")
		}
		actionRequestMapList = make([]map[string]interface{}, 0, len(m))
		for _, v := range m {
			mm, ok := v.(map[string]interface{})
			if !ok || mm == nil {
				return nil, errors.New("不是一个 JSON 对象数组")
			}
			actionRequestMapList = append(actionRequestMapList, mm)
		}
	}
	requests := make([]Request, 0, len(actionRequestMapList))
	for _, actionRequestMap := range actionRequestMapList {
		r, err := parseRequestFromMap(actionRequestMap, reqComm)
		if err != nil {
			return nil, err
		}
		requests = append(requests, r)
	}
	return requests, nil
}
