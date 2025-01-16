package controller

import (
	"book-fiber/home/service"
	"book-fiber/model"
	system "book-fiber/system/controller"
	"github.com/gofiber/fiber/v2"
)

type BookController struct {
	system.BaseController
	service *service.BookService
}

func NewBookController(base system.BaseController, service *service.BookService) *BookController {
	return &BookController{BaseController: base, service: service}
}

func (c *BookController) GetBookList(ctx *fiber.Ctx) error {
	query := new(model.BookQuery)

	if err := ctx.QueryParser(query); err != nil {
		return c.Error(ctx, "Params Error")
	}
	result, err := c.service.GetList(query)
	if err != nil {
		return c.Error(ctx)
	}

	return c.Success(ctx, result)
}

func (c *BookController) GetOne(ctx *fiber.Ctx) error {
	result, err := c.service.GetOne(ctx.Params("slug"))
	if err != nil {
		return c.Error(ctx)
	}
	return c.Success(ctx, result)
}

func (c *BookController) GetChapters(ctx *fiber.Ctx) error {
	result, err := c.service.GetChapters(ctx.Params("slug"))
	if err != nil {
		return c.Error(ctx)
	}
	return c.Success(ctx, result)
}

func (c *BookController) GetChapter(ctx *fiber.Ctx) error {
	result, err := c.service.GetChapter(ctx.Params("slug"))
	if err != nil {
		return c.Error(ctx)
	}
	return c.Success(ctx, result)
}
