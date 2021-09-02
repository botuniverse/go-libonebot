package utils

import (
	"fmt"
	"strconv"
)

// EasyMap 封装了一个类 map[string]interface{} 的数据结构,
// 提供了更加方便的方法来获取其中的字段值.
type EasyMap struct {
	inner map[string]interface{}
}

// EasyMapFromMap 从 map[string]interface{} 创建一个 EasyMap.
//
// 参数:
//   m: 要封装的 map[string]interface{}, 不能为 nil
func EasyMapFromMap(m map[string]interface{}) EasyMap {
	if m == nil {
		panic("must not be nil")
	}
	return EasyMap{m}
}

// Value 获取 EasyMap 内部数据结构的 map[string]interface{} 形式.
func (m EasyMap) Value() map[string]interface{} {
	return m.inner
}

func (m EasyMap) errorMissingField(key string) error {
	return fmt.Errorf("`%v` 字段不存在", key)
}

func (m EasyMap) errorInvalidField(key string) error {
	return fmt.Errorf("`%v` 字段是无效值", key)
}

// Get 获取任意类型的字段值.
func (m EasyMap) Get(key string) (interface{}, error) {
	val, ok := m.inner[key]
	if !ok {
		return nil, m.errorMissingField(key)
	}
	return val, nil
}

// GetBool 获取布尔类型的字段值, 如果字段是字符串, 则会尝试转换.
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

// GetInt64 获取整数类型的字段值, 如果字段是字符串, 则会尝试转换.
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

// GetFloat64 获取浮点数类型的字段值, 如果字段是字符串, 则会尝试转换.
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

// GetString 获取字符串类型的字段值.
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

// GetBytes 获取字节数组类型的字段值.
func (m EasyMap) GetBytes(key string) ([]byte, error) {
	val, err := m.Get(key)
	if err != nil {
		return nil, err
	}
	if val == nil {
		return nil, m.errorInvalidField(key)
	}
	switch val := val.(type) {
	case []byte:
		return val, nil
	default:
		return nil, m.errorInvalidField(key)
	}
}

// GetMap 获取 map[string]interface{} 类型的字段值, 并封装为 EasyMap.
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

// GetArray 获取 []interface{} 类型的字段值.
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

// GetMapArray 获取 []map[string]interface{} 类型的字段值, 并封装为 []EasyMap.
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

// Set 设置任意类型的字段值.
func (m EasyMap) Set(key string, value interface{}) {
	m.inner[key] = value
}
