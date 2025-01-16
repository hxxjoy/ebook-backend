package controller

import (
	"book-fiber/system/helper"

	"github.com/gofiber/fiber/v2"
)

type BaseController interface {
	Success(ctx *fiber.Ctx, data interface{}) error
	Error(ctx *fiber.Ctx, message ...string) error
	ErrorAuth(ctx *fiber.Ctx) error
	Response(ctx *fiber.Ctx, response helper.Response) error
}

// 2. 基础实现
type baseController struct {
	http helper.HttpHelper
}

func NewBaseController(
	http helper.HttpHelper,
) BaseController {
	return &baseController{
		http: http,
	}
}

// 实现接口方法
func (b *baseController) Success(ctx *fiber.Ctx, data interface{}) error {
	return b.http.Success(ctx, data)
}

func (b *baseController) Error(ctx *fiber.Ctx, message ...string) error {
	return b.http.Error(ctx, message...)
}

func (b *baseController) ErrorAuth(ctx *fiber.Ctx) error {
	return b.http.ErrorAuth(ctx)
}

func (b *baseController) Response(ctx *fiber.Ctx, data helper.Response) error {
	return b.http.Response(ctx, data)
}
