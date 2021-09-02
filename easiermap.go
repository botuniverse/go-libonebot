package libonebot

import (
	"github.com/botuniverse/go-libonebot/utils"
)

// EasierMap 在 EasyMap 的基础上添加了 GetMessage 方法.
type EasierMap struct {
	utils.EasyMap
}

// EasierMapFromMap 从 map[string]interface{} 创建一个 EasierMap.
//
// 参数:
//   m: 要封装的 map[string]interface{}, 不能为 nil
func EasierMapFromMap(m map[string]interface{}) EasierMap {
	return EasierMapFromEasyMap(utils.EasyMapFromMap(m))
}

// EasierMapFromEasyMap 从 EasyMap 创建一个 EasierMap.
func EasierMapFromEasyMap(m utils.EasyMap) EasierMap {
	return EasierMap{
		EasyMap: m,
	}
}

// GetMap 获取 map[string]interface{} 类型的字段值, 并封装为 EasierMap.
func (m EasierMap) GetMap(key string) (EasierMap, error) {
	val, err := m.EasyMap.GetMap(key)
	if err != nil {
		return EasierMap{}, err
	}
	return EasierMapFromEasyMap(val), nil
}

// GetMapArray 获取 []map[string]interface{} 类型的字段值, 并封装为 []EasierMap.
func (m EasierMap) GetMapArray(key string) ([]EasierMap, error) {
	val, err := m.EasyMap.GetMapArray(key)
	if err != nil {
		return nil, err
	}
	arr := make([]EasierMap, len(val))
	for i, v := range val {
		arr[i] = EasierMapFromEasyMap(v)
	}
	return arr, nil
}

// GetMessage 获取 Message 类型的字段值.
func (m EasierMap) GetMessage(key string) (Message, error) {
	val, err := m.Get(key)
	if err != nil {
		return nil, err
	}
	return messageFromInterface(val)
}
