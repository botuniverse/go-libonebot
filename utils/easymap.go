package utils

import (
	"fmt"
	"strconv"
)

type EasyMap struct {
	inner map[string]interface{}
}

func EasyMapFromMap(m map[string]interface{}) EasyMap {
	if m == nil {
		panic("must not be nil")
	}
	return EasyMap{m}
}

func (m EasyMap) Value() map[string]interface{} {
	return m.inner
}

func (m EasyMap) errorMissingField(key string) error {
	return fmt.Errorf("`%v` 字段不存在", key)
}

func (m EasyMap) errorInvalidField(key string) error {
	return fmt.Errorf("`%v` 字段是无效值", key)
}

func (m EasyMap) Get(key string) (interface{}, error) {
	val, ok := m.inner[key]
	if !ok {
		return nil, m.errorMissingField(key)
	}
	return val, nil
}

func (m EasyMap) GetBool(key string) (bool, error) {
	val, err := m.Get(key)
	if err != nil {
		return false, err
	}
	switch val := val.(type) {
	case bool:
		return val, nil
	case string:
		return val == "true", nil
	default:
		return false, m.errorInvalidField(key)
	}
}

func (m EasyMap) GetInt64(key string) (int64, error) {
	val, err := m.Get(key)
	if err != nil {
		return 0, err
	}
	switch val := val.(type) {
	case int64:
		return val, nil
	case int32:
		return int64(val), nil
	case int:
		return int64(val), nil
	case float64:
		return int64(val), nil
	case float32:
		return int64(val), nil
	case string:
		return strconv.ParseInt(val, 10, 64)
	default:
		return 0, m.errorInvalidField(key)
	}
}

func (m EasyMap) GetFloat64(key string) (float64, error) {
	val, err := m.Get(key)
	if err != nil {
		return 0, err
	}
	switch val := val.(type) {
	case float64:
		return val, nil
	case float32:
		return float64(val), nil
	case int64:
		return float64(val), nil
	case int32:
		return float64(val), nil
	case int:
		return float64(val), nil
	case string:
		return strconv.ParseFloat(val, 64)
	default:
		return 0, m.errorInvalidField(key)
	}
}

func (m EasyMap) GetString(key string) (string, error) {
	val, err := m.Get(key)
	if err != nil {
		return "", err
	}
	if val == nil {
		return "", m.errorInvalidField(key)
	}
	switch val := val.(type) {
	case string:
		return val, nil
	default:
		return "", m.errorInvalidField(key)
	}
}

func (m EasyMap) GetMap(key string) (EasyMap, error) {
	val, err := m.Get(key)
	if err != nil {
		return EasyMap{}, err
	}
	if val == nil {
		return EasyMap{}, m.errorInvalidField(key)
	}
	switch val := val.(type) {
	case map[string]interface{}:
		return EasyMapFromMap(val), nil
	default:
		return EasyMap{}, m.errorInvalidField(key)
	}
}

func (m EasyMap) GetArray(key string) ([]interface{}, error) {
	val, err := m.Get(key)
	if err != nil {
		return nil, err
	}
	if val == nil {
		return nil, m.errorInvalidField(key)
	}
	switch val := val.(type) {
	case []interface{}:
		return val, nil
	default:
		return nil, m.errorInvalidField(key)
	}
}

func (m EasyMap) GetMapArray(key string) ([]EasyMap, error) {
	val, err := m.Get(key)
	if err != nil {
		return nil, err
	}
	if val == nil {
		return nil, m.errorInvalidField(key)
	}
	switch val := val.(type) {
	case []map[string]interface{}:
		maps := make([]EasyMap, len(val))
		for i, m := range val {
			maps[i] = EasyMapFromMap(m)
		}
		return maps, nil
	default:
		return nil, m.errorInvalidField(key)
	}
}

func (m EasyMap) Set(key string, value interface{}) {
	m.inner[key] = value
}
