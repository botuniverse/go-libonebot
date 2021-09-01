package libonebot

import (
	"errors"
	"fmt"
	"strings"

	"github.com/botuniverse/go-libonebot/utils"
	"github.com/tidwall/gjson"
	"github.com/vmihailenco/msgpack/v5"
)

type Request struct {
	Action Action
	Params *easyMap
	Echo   interface{}
}

func validateActionJSON(actionJSON gjson.Result) error {
	if !actionJSON.Get("action").Exists() {
		return errors.New("动作请求缺少 `action` 字段")
	}
	if actionJSON.Get("action").String() == "" {
		return errors.New("动作请求的 `action` 字段为空")
	}
	if !actionJSON.Get("params").Exists() {
		return errors.New("动作请求缺少 `params` 字段")
	}
	if !actionJSON.Get("params").IsObject() {
		return errors.New("动作请求的 `params` 字段不是一个 JSON 对象")
	}
	return nil
}

func parseTextActionRequest(prefix string, actionBytes []byte) (Request, error) {
	if !gjson.ValidBytes(actionBytes) {
		return Request{}, errors.New("动作请求体不是合法的 JSON")
	}

	actionJSON := gjson.Parse(utils.BytesToString(actionBytes))
	err := validateActionJSON(actionJSON)
	if err != nil {
		return Request{}, err
	}

	var action Action
	fullname := actionJSON.Get("action").String()
	prefix_ul := prefix + "_"
	if strings.HasPrefix(fullname, prefix_ul) {
		// extended action
		action = Action{
			Prefix:     prefix,
			Name:       strings.TrimPrefix(fullname, prefix_ul),
			IsExtended: true,
		}
	} else {
		// core action
		action = Action{
			Prefix:     "",
			Name:       fullname,
			IsExtended: false,
		}
	}

	r := Request{
		Action: action,
		Params: newEasyMapFromJSON(actionJSON.Get("params")),
		Echo:   actionJSON.Get("echo").Value(),
	}
	return r, nil
}

func parseBinaryActionRequest(prefix string, actionBytes []byte) (Request, error) {
	// TODO
	var actionMap map[string]interface{}
	err := msgpack.Unmarshal(actionBytes, &actionMap)
	if err != nil {
		return Request{}, err
	}
	fmt.Printf("actionMap: %#v\n", actionMap)
	panic("")
}
