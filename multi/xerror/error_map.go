package xerror

import (
	"github.com/olongfen/toolkit/consts"
	"sync"
)

var (
	// DefaultErrorMul default
	DefaultErrorMul = &ErrorMul{
		data: sync.Map{},
	}
)

func init() {
	DefaultErrorMul.Set(IllegalAccessToken, consts.SimplifiedChinese, "非法token")
	DefaultErrorMul.Set(IllegalAccessToken, consts.TraditionalChinese, "非法token")
	DefaultErrorMul.Set(IllegalAccessToken, consts.English, "Illegal token")

	DefaultErrorMul.Set(IllegalCertificate, consts.SimplifiedChinese, "非法凭证")
	DefaultErrorMul.Set(IllegalCertificate, consts.TraditionalChinese, "非法憑證")
	DefaultErrorMul.Set(IllegalCertificate, consts.English, "Illegal certificate")

	DefaultErrorMul.Set(IllegalParameter, consts.SimplifiedChinese, "非法参数")
	DefaultErrorMul.Set(IllegalParameter, consts.TraditionalChinese, "非法參數")
	DefaultErrorMul.Set(IllegalParameter, consts.English, "Illegal parameter")

	DefaultErrorMul.Set(RecordNotFound, consts.SimplifiedChinese, "记录未找到")
	DefaultErrorMul.Set(RecordNotFound, consts.TraditionalChinese, "記錄未找到")
	DefaultErrorMul.Set(RecordNotFound, consts.English, "record not found")

	DefaultErrorMul.Set(AlreadyExists, consts.SimplifiedChinese, "已经存在,不允许重复创建")
	DefaultErrorMul.Set(AlreadyExists, consts.TraditionalChinese, "已經存在,不允許重複創建")
	DefaultErrorMul.Set(AlreadyExists, consts.English, "already exists,duplicate creation is not allowed")

	DefaultErrorMul.Set(SortParameterMismatch, consts.SimplifiedChinese, "排序参数不匹配")
	DefaultErrorMul.Set(SortParameterMismatch, consts.TraditionalChinese, "排序參數不匹配")
	DefaultErrorMul.Set(SortParameterMismatch, consts.English, "sort parameter mismatch")

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
