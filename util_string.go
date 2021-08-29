package libonebot

import (
	"reflect"
	"unsafe"
)

func stringToBytes(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}

func bytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
