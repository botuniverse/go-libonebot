package libonebot

import (
	"github.com/botuniverse/go-libonebot/utils"
)

type easierMap struct {
	utils.EasyMap
}

func easierMapFromMap(m map[string]interface{}) easierMap {
	return easierMapFromEasyMap(utils.EasyMapFromMap(m))
}

func easierMapFromEasyMap(m utils.EasyMap) easierMap {
	return easierMap{
		EasyMap: m,
	}
}

func (m easierMap) GetMap(key string) (easierMap, error) {
	val, err := m.EasyMap.GetMap(key)
	if err != nil {
		return easierMap{}, err
	}
	return easierMapFromEasyMap(val), nil
}

func (m easierMap) GetMapArray(key string) ([]easierMap, error) {
	val, err := m.EasyMap.GetMapArray(key)
	if err != nil {
		return nil, err
	}
	arr := make([]easierMap, len(val))
	for i, v := range val {
		arr[i] = easierMapFromEasyMap(v)
	}
	return arr, nil
}

func (m easierMap) GetMessage(key string) (Message, error) {
	val, err := m.Get(key)
	if err != nil {
		return nil, err
	}
	return messageFromInterface(val)
}
