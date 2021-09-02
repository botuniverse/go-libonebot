package libonebot

import (
	"github.com/botuniverse/go-libonebot/utils"
)

type EasierMap struct {
	utils.EasyMap
}

func EasierMapFromMap(m map[string]interface{}) EasierMap {
	return EasierMapFromEasyMap(utils.EasyMapFromMap(m))
}

func EasierMapFromEasyMap(m utils.EasyMap) EasierMap {
	return EasierMap{
		EasyMap: m,
	}
}

func (m EasierMap) GetMap(key string) (EasierMap, error) {
	val, err := m.EasyMap.GetMap(key)
	if err != nil {
		return EasierMap{}, err
	}
	return EasierMapFromEasyMap(val), nil
}

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

func (m EasierMap) GetMessage(key string) (Message, error) {
	val, err := m.Get(key)
	if err != nil {
		return nil, err
	}
	return messageFromInterface(val)
}
