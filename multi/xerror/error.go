package xerror

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
)

type ValidateError map[string]string

func (v ValidateError) Error() string {
	b, _ := json.Marshal(v)
	return fmt.Sprintf("%s", b)
}

var _ BizError = (*bizError)(nil)

type BizError interface {
	// i 为了避免被其他包实现
	i()
	Code() int
	// WithError 设置错误信息
	WithError(err error) BizError

	Error() string

	// StackError 获取带堆栈的错误信息
	StackError() error
}

type bizError struct {
	code    int
	message string // 错误描述
	stack   error  // 含有堆栈信息的错误
}

type DBErrorResponse map[string]BizError

func (c DBErrorResponse) Error() string {
	b, _ := json.Marshal(c)
	return string(b)
}

func NewError(code int, language string, errMul ...*ErrorMul) BizError {
	biz := &bizError{
		message: "",
		code:    code,
	}
	var (
		defaultErrMul = DefaultErrorMul
	)
	if len(errMul) > 0 {
		defaultErrMul = errMul[0]
	}
	biz.message = defaultErrMul.Get(code, language)
	return biz
}

func (e *bizError) i() {}

func (e *bizError) Code() int {
	return e.code
}

func (e *bizError) Error() string {
	return fmt.Sprintf(`%s`, e.message)
}

func (e *bizError) WithError(err error) BizError {
	e.stack = errors.WithStack(err)
	return e
}

func (e *bizError) Message() string {
	return e.message
}

func (e *bizError) StackError() error {
	return e.stack
}
