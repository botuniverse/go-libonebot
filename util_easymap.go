package onebot

import (
	"encoding/json"
	"fmt"

	"github.com/tidwall/gjson"
)

type easyMap struct {
	JSON gjson.Result
}

func easyMapFromMap(m map[string]interface{}) easyMap {
	j, _ := json.Marshal(m)
	return easyMap{gjson.Parse(bytesToString(j))}
}

func easyMapFromJSON(j gjson.Result) easyMap {
	return easyMap{j}
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
	val, err := m.Get(key)
	if err != nil {
		return nil, err
	}

	if val.Type == gjson.String {
		return Message{TextSegment(val.Str)}, nil
	}

	if val.IsObject() {
		return MessageFromJSON("[" + val.Raw + "]")
	} else if val.IsArray() {
		return MessageFromJSON(val.Raw)
	} else {
		return nil, errorInvalidField(key)
	}
}

func (m *easyMap) GetMap(key string) (easyMap, error) {
	val, err := m.Get(key)
	if err != nil {
		return easyMap{}, err
	}

	if val.IsObject() {
		return easyMap{val}, nil
	} else {
		return easyMap{}, errorInvalidField(key)
	}
}

func (m *easyMap) GetArray(key string) ([]interface{}, error) {
	val, err := m.Get(key)
	if err != nil {
		return nil, err
	}

	if val.IsArray() {
		return val.Value().([]interface{}), nil
	} else {
		return nil, errorInvalidField(key)
	}
}
