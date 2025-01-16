package controller

import (
	"book-fiber/home/service"
	system "book-fiber/system/controller"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

type UserController struct {
	system.BaseController
	service      *service.UserService
	emailService service.EmailService
}

func NewUserController(base system.BaseController, service *service.UserService, emailService service.EmailService) *UserController {
	return &UserController{BaseController: base, service: service, emailService: emailService}
}

func (c *UserController) GetProfile(ctx *fiber.Ctx) error {
	slug := ctx.Locals("slug").(string)
	user, err := c.service.GetBySlug(slug)
	if err != nil {
		return c.Error(ctx, "User not found")
	}
	return c.Success(ctx, user)
}

func (c *UserController) UpdateProfile(ctx *fiber.Ctx) error {
	// 更新用户资料
	return nil
}

func (c *UserController) ChangePassword(ctx *fiber.Ctx) error {
	// 修改密码
	return nil
}
