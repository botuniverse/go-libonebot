package onebot

import (
	"encoding/json"
	"fmt"

	"github.com/tidwall/gjson"
)

type easyMap struct {
	JSON gjson.Result
}

func newEasyMapFromMap(m map[string]interface{}) *easyMap {
	j, _ := json.Marshal(m)
	return &easyMap{gjson.Parse(bytesToString(j))}
}

func newEasyMapFromJSON(j gjson.Result) *easyMap {
	return &easyMap{j}
}

func (m easyMap) MarshalJSON() ([]byte, error) {
	return stringToBytes(m.JSON.Raw), nil
}

func (m *easyMap) UnmarshalJSON(data []byte) error {
	s := bytesToString(data)
	if !gjson.Valid(s) {
		return fmt.Errorf("JSON 语法错误")
	}
	j := gjson.Parse(s)
	if !j.IsObject() {
		return fmt.Errorf("必须是 JSON 对象")
	}
	m.JSON = j
	return nil
}

func errorMissingField(key string) error {
	return fmt.Errorf("`%v` 字段不存在", key)
}

func errorInvalidField(key string) error {
	return fmt.Errorf("`%v` 字段无效", key)
}

func (m *easyMap) Get(key string) (gjson.Result, error) {
	val := m.JSON.Get(key)
	if !val.Exists() {
		return gjson.Result{}, errorMissingField(key)
	}
	return val, nil
}

func (m *easyMap) GetBool(key string) (bool, error) {
	val := m.JSON.Get(key)
	if !val.Exists() {
		return false, errorMissingField(key)
	}
	if val.Type != gjson.True && val.Type != gjson.False {
		return false, errorInvalidField(key)
	}
	return val.Bool(), nil
}

func (m *easyMap) GetInt64(key string) (int64, error) {
	val := m.JSON.Get(key)
	if !val.Exists() {
		return 0, errorMissingField(key)
	}
	if val.Type != gjson.Number {
		return 0, errorInvalidField(key)
	}
	return val.Int(), nil
}

func (m *easyMap) GetFloat64(key string) (float64, error) {
	val := m.JSON.Get(key)
	if !val.Exists() {
		return 0, errorMissingField(key)
	}
	if val.Type != gjson.Number {
		return 0, errorInvalidField(key)
	}
	return val.Float(), nil
}

func (m *easyMap) GetString(key string) (string, error) {
	val := m.JSON.Get(key)
	if !val.Exists() {
		return "", errorMissingField(key)
	}
	if val.Type != gjson.String {
		return "", errorInvalidField(key)
	}
	return val.Str, nil
}

func (m *easyMap) GetMessage(key string) (Message, error) {
	val := m.JSON.Get(key)
	if !val.Exists() {
		return Message{}, errorMissingField(key)
	}
	return MessageFromJSON(val)
}

func (m *easyMap) GetMap(key string) (easyMap, error) {
	val := m.JSON.Get(key)
	if !val.Exists() {
		return easyMap{}, errorMissingField(key)
	}

	if val.IsObject() {
		return easyMap{val}, nil
	} else {
		return easyMap{}, errorInvalidField(key)
	}
}

func (m *easyMap) GetArray(key string) ([]interface{}, error) {
	val := m.JSON.Get(key)
	if !val.Exists() {
		return nil, errorMissingField(key)
	}

	if val.IsArray() {
		return val.Value().([]interface{}), nil
	} else {
		return nil, errorInvalidField(key)
	}
}
