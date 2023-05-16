package response

import (
	fiber "github.com/gofiber/fiber/v2"
	"github.com/olongfen/toolkit/consts"
	"github.com/olongfen/toolkit/multi/xerror"
	"github.com/olongfen/toolkit/scontext"
	"net/http"
)

// Response http response
type Response struct {
	status int
	//
	Code     int         `json:"code"`
	Data     interface{} `json:"data"`
	Message  string      `json:"message"`
	Language string      `json:"language"`
	Errors   interface{} `json:"errors"`
}

// NewResponse new
func NewResponse() *Response {
	return &Response{status: http.StatusOK}
}

// SetErrors set error
func (r *Response) SetErrors(errs interface{}) *Response {
	r.Errors = errs
	return r
}

// SetMessage set message
func (r *Response) SetMessage(msg string) *Response {
	r.Message = msg
	return r
}

// Success response success
func (r *Response) Success(ctx *fiber.Ctx, data interface{}) error {
	r.Data = data
	userCtx := ctx.UserContext()
	lan := scontext.GetLanguage(userCtx)
	if r.Message == "" {
		switch lan {
		case consts.English:
			r.Message = "success"
		case consts.TraditionalChinese:
			r.Message = "成功"
		default:
			r.Message = "成功"
		}

	}
	r.Language = scontext.GetLanguage(ctx.UserContext())
	return ctx.Status(r.status).JSON(r)
}

// ErrorHandler fiber error handler
var ErrorHandler = func(ctx *fiber.Ctx, err error) error {
	status := fiber.StatusOK
	userCtx := ctx.UserContext()
	lan := scontext.GetLanguage(userCtx)
	resp := NewResponse()
	switch lan {
	case consts.English:
		resp.Message = "failed"
	case consts.TraditionalChinese:
		resp.Message = "失敗"
	default:
		resp.Message = "失败"
	}
	resp.Code = -1
	switch err.(type) {
	case *fiber.Error:
		// 处理内部错误返回
		e := err.(*fiber.Error)
		resp.Code = e.Code
	case xerror.BizError:
		// 处理自定义业务错误返回
		e := err.(xerror.BizError)
		resp.Code = e.Code()
		resp.Message = e.Error()
	case xerror.ValidateError:
		resp.SetErrors(err.(xerror.ValidateError))
		resp.Code = xerror.IllegalParameter
		resp.Message = xerror.NewError(xerror.IllegalParameter, lan).Error()
	case xerror.DBErrorResponse:
		// 处理数据库错误返回
		var (
			m = map[string]string{}
		)
		e := err.(xerror.DBErrorResponse)
		for k, v := range e {
			resp.Code = v.Code()
			m[k] = v.Error()
		}
		resp.SetErrors(m)
	default:
		status = fiber.StatusInternalServerError
	}
	/*	// 处理内部错误返回
		if e, ok := err.(*fiber.Error); ok {
			xlog.Log.Error("HTTP Error", zap.Error(err))
			resp.Code = e.Code
			resp.Message = "failed"
		}
		// 处理自定义业务错误返回
		if e, ok := err.(err_mul.BizError); ok {
			xlog.Log.Error("Business Error", zap.Error(e.StackError()))
			resp.Code = e.Code()
			resp.Message = e.Error()
		}
		// 处理数据库错误返回
		if e, ok := err.(err_mul.DBErrorResponse); ok {
			var (
				m = map[string]string{}
			)
			for k, v := range e {
				resp.Code = v.Code()
				m[k] = v.Error()
			}
			resp.Message = "failed"
			resp.SetErrors(m)
		}*/

	if resp.Errors == nil {
		resp.Errors = map[string]any{"error": err.Error()}
	}

	return ctx.Status(status).JSON(resp)
}
