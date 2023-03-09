package xerror

import (
	"sync"
)

var (
	// DefaultErrorMul default
	DefaultErrorMul = &ErrorMul{
		data: sync.Map{},
	}
)

func init() {
	DefaultErrorMul.Set(IllegalAccessToken, "zh-cn", "非法token")
	DefaultErrorMul.Set(IllegalAccessToken, "en", "Illegal token")
	DefaultErrorMul.Set(IllegalCertificate, "zh-cn", "非法凭证")
	DefaultErrorMul.Set(IllegalCertificate, "en", "Illegal certificate")
	DefaultErrorMul.Set(IllegalParameter, "zh-cn", "非法参数")
	DefaultErrorMul.Set(IllegalParameter, "en", "Illegal parameter")
	DefaultErrorMul.Set(RecordNotFound, "zh-cn", "记录未找到")
	DefaultErrorMul.Set(RecordNotFound, "en", "record not found")
	DefaultErrorMul.Set(AlreadyExists, "zh-cn", "已经存在")
	DefaultErrorMul.Set(AlreadyExists, "en", "already exists")
	DefaultErrorMul.Set(SortParameterMismatch, "zh-cn", "排序参数不匹配")
	DefaultErrorMul.Set(SortParameterMismatch, "en", "sort parameter mismatch")

}

// ErrorMul error multi-language
type ErrorMul struct {
	data sync.Map
}

// Set send key lan val to mul
func (e *ErrorMul) Set(key int, lan string, val string) *ErrorMul {
	if v, ok := e.data.Load(key); !ok {
		e.data.Store(key, map[string]string{lan: val})
	} else {
		v.(map[string]string)[lan] = val
	}
	return e
}

// Get send key lan to get value
func (e *ErrorMul) Get(key int, lan string) string {
	if val, ok := e.data.Load(key); ok {
		return val.(map[string]string)[lan]
	}
	return ""
}

// DeleteKey delete key
func (e *ErrorMul) DeleteKey(key int) *ErrorMul {
	e.data.Delete(key)
	return e
}
