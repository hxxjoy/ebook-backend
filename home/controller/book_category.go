package controller

import (
	"book-fiber/home/service"
	system "book-fiber/system/controller"
	"github.com/gofiber/fiber/v2"
)

type BookCategoryController struct {
	system.BaseController
	service *service.BookCategoryService
}

func NewBookCategoryController(base system.BaseController, service *service.BookCategoryService) *BookCategoryController {
	return &BookCategoryController{BaseController: base, service: service}
}

func (c *BookCategoryController) GetList(ctx *fiber.Ctx) error {
    result, err := c.service.BuildCategoryTree()
    if err != nil {
		return c.Error(ctx)
	}
	return c.Success(ctx, result)
}

func (c *BookCategoryController) GetCategoryBooks(ctx *fiber.Ctx) error {
    slug := ctx.Params("slug")  // 从URL参数获取分类ID
    page, err := ctx.ParamsInt("page") // 获取 "3" 并转换为整数
    if err != nil || page < 1 {
        page = 1 // 如果转换失败或小于1，设置默认值
    }
    pageSize := ctx.QueryInt("page_size", 24) // 获取每页数量，默认12条

    // 调用service层获取数据
    data, err := c.service.GetBooksByCategory(slug, page, pageSize)
    if err != nil {
        return c.Error(ctx)
    }
    return c.Success(ctx,data)
}

func (c *BookController) GetBooks(ctx *fiber.Ctx) error {
	result, err := c.service.GetChapters(ctx.Params("slug"))
	if err != nil {
		return c.Error(ctx)
	}
	return c.Success(ctx, result)
}
