package controller

import (
	"book-fiber/home/service"
	"book-fiber/model"
	system "book-fiber/system/controller"
	"book-fiber/system/helper"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	system.BaseController
    userService  *service.UserService
    authService  *service.AuthService
	emailService service.EmailService
}

type LoginRequest struct {
    Email       string `json:"email" validate:"required,email"`
    Password    string `json:"password" validate:"required"`
    RememberMe  bool   `json:"remember_me"`
}

type RegisterRequest struct {
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required,min=8,max=26"`
	//Nickname string `json:"nickname" form:"nickname" validate:"required,min=2,max=26"`
	Code string `json:"code" form:"code" validate:"required,len=6"`
}

type SendEmailRequest struct {
	Email string `json:"email"`
}

func NewAuthController(base system.BaseController, userService *service.UserService, authService *service.AuthService,emailService service.EmailService) *AuthController {
	return &AuthController{BaseController: base, userService: userService, authService: authService, emailService:  emailService}
}

func (c *AuthController) Login(ctx *fiber.Ctx) error {
    var req LoginRequest
    if err := ctx.BodyParser(&req); err != nil {
        return c.Error(ctx, "Invalid request")
    }

    // 验证请求参数
    if err := validate.Struct(req); err != nil {
        return c.Error(ctx, err.Error())
    }

    // 登录处理
    user, tokens, err := c.authService.Login(req.Email, req.Password, req.RememberMe)
    if err != nil {
        return c.Error(ctx, err.Error())
    }

    // 设置刷新令牌到Cookie
    ctx.Cookie(&fiber.Cookie{
        Name:     "refresh_token",
        Value:    tokens.RefreshToken,
        Expires:  time.Now().Add(4320 * time.Hour),
        HTTPOnly: true,
        Secure:   true,
        SameSite: "Strict",
    })

    return c.Success(ctx, fiber.Map{
        "user":         user,
        "access_token": tokens.AccessToken,
    })
}

func (c *AuthController) Register(ctx *fiber.Ctx) error {
	var req RegisterRequest
	// 解析请求体
	if err := ctx.BodyParser(&req); err != nil {
		return c.Error(ctx, "Params Error")
	}
	if err := validate.Struct(req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return c.Error(ctx, validationErrors.Error())
	}
	ip := ctx.IP()
	checkCode := c.emailService.VerifyCode(req.Email, req.Code, ip)
	if !checkCode {
		return c.Error(ctx, "Verification code is incorrect")
	}
	// 检查邮箱是否已存在
	user, _ := c.userService.GetByEmail(req.Email)
	if user != nil {
		return c.Error(ctx, "Email Exists")
	}
	now := time.Now()
	nickname := helper.GenerateNickname()
	user = &model.User{
		Email:        &req.Email,
		Username:     &req.Email,
		PasswordHash: &req.Password,
		Nickname:     &nickname,
		CreatedAt:    &now,
	}

	result, err := c.userService.Create(user)
	if err != nil {
		return c.Error(ctx, "User registration failed, please try again later")
	}
	return c.Success(ctx, result)
}

func (c *AuthController) SendEmailCode(ctx *fiber.Ctx) error {
	// 获取邮箱地址
	var req SendEmailRequest
	if err := ctx.BodyParser(&req); err != nil {
		return c.Error(ctx, "Params Error.")
	}
	ip := ctx.IP()
	// 发送邮件
	err := c.emailService.SendVerificationCode(req.Email, ip)
	if err != nil {
		return c.Error(ctx, err.Error())
	}

	return c.Success(ctx, "Verification code has been sent")
}

func (c *AuthController) RefreshToken(ctx *fiber.Ctx) error {
    refreshToken := ctx.Cookies("refresh_token")
    if refreshToken == "" {
        return c.Error(ctx, "No refresh token")
    }

    tokens, err := c.authService.RefreshTokens(refreshToken)
    if err != nil {
        return c.Error(ctx, "Invalid refresh token")
    }

    // 更新Cookie中的刷新令牌
    ctx.Cookie(&fiber.Cookie{
        Name:     "refresh_token",
        Value:    tokens.RefreshToken,
        Expires:  time.Now().Add(4320 * time.Hour),
        HTTPOnly: true,
        Secure:   true,
        SameSite: "Strict",
    })

    return c.Success(ctx, fiber.Map{
        "access_token": tokens.AccessToken,
    })
}

func (c *AuthController) Logout(ctx *fiber.Ctx) error {
    token, err := helper.GetTokenFromHeader(ctx)
    if err != nil {
        return c.Error(ctx, "Invalid token")
    }

    if err := c.authService.Logout(token); err != nil {
        return c.Error(ctx, "Logout failed")
    }

    // 清除刷新令牌Cookie
    ctx.Cookie(&fiber.Cookie{
        Name:     "refresh_token",
        Value:    "",
        Expires:  time.Now().Add(-24 * time.Hour),
        HTTPOnly: true,
        Secure:   true,
        SameSite: "Strict",
    })

    return c.Success(ctx, "Logged out successfully")
}