// helper/response.go

package helper

import (
	"github.com/gofiber/fiber/v2"
)

// 1. 先定义接口
type HttpHelper interface {
	Error(ctx *fiber.Ctx, message ...string) error
	ErrorAuth(ctx *fiber.Ctx) error
	Success(ctx *fiber.Ctx, data interface{}) error
	Response(ctx *fiber.Ctx, response Response) error
}

// 2. 实现接口
type httpHelper struct {
	// 可以添加需要的依赖
	//logger *log.Logger
	//config *config.Config
}

func NewHttpHelper() HttpHelper {
	return &httpHelper{}
}

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// 成功响应构造器
func Success(data any) Response {
	return Response{
		Code:    1,
		Message: "ok",
		Data:    data,
	}
}

// 错误响应构造器
func Error(message string) Response {
	return Response{
		Code:    0,
		Message: message,
		Data:    []interface{}{},
	}
}

// 未授权响应
func ErrorAuth() Response {
	return Response{
		Code:    -1,
		Message: "Please sign in",
		Data:    []interface{}{},
	}
}

// 自定义响应构造器
func Custom(code int, message string, data any) Response {
	return Response{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

// HTTP响应包装器
func (h *httpHelper) Response(ctx *fiber.Ctx, response Response) error {
	return ctx.Status(fiber.StatusOK).JSON(response)
}

// 便捷方法
func (h *httpHelper) Success(ctx *fiber.Ctx, data any) error {
	return h.Response(ctx, Success(data))
}

func (h *httpHelper) Error(ctx *fiber.Ctx, message ...string) error {
	defaultMessage := "Internal error"
	if len(message) > 0 && message[0] != "" {
		defaultMessage = message[0]
	}
	return h.Response(ctx, Error(defaultMessage))
}

func (h *httpHelper) ErrorAuth(ctx *fiber.Ctx) error {
	return h.Response(ctx, ErrorAuth())
}
