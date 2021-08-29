package libonebot

import (
	"errors"
	"strings"

	"github.com/tidwall/gjson"
)

type Request struct {
	Action Action
	Params *easyMap
	Echo   interface{}
}

func validateActionJSON(actionJSON gjson.Result) error {
	if !actionJSON.Get("action").Exists() {
		return errors.New("Action 请求体缺少 `action` 字段")
	}
	if actionJSON.Get("action").String() == "" {
		return errors.New("Action 请求体的 `action` 字段为空")
	}
	if !actionJSON.Get("params").Exists() {
		return errors.New("Action 请求体缺少 `params` 字段")
	}
	if !actionJSON.Get("params").IsObject() {
		return errors.New("Action 请求体的 `params` 字段不是一个 JSON 对象")
	}
	return nil
}

func parseActionRequest(prefix string, actionBody string) (Request, error) {
	if !gjson.Valid(actionBody) {
		return Request{}, errors.New("Action 请求体不是合法的 JSON")
	}

	actionJSON := gjson.Parse(actionBody)
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
