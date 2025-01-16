package middleware

import (
	"book-fiber/home/service"
	"book-fiber/system/helper"

	"github.com/gofiber/fiber/v2"
)

func Auth(authService *service.AuthService) fiber.Handler {
    return func(c *fiber.Ctx) error {
        // 1. 获取token
        token, err := helper.GetTokenFromHeader(c)
        if err != nil {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "code": -2,
                "msg":  "Unauthorized",
                "data": nil,
            })
        }

        // 2. 解析token获取slug
        slug, err := helper.ParseTokenAndGetSlug(token)
        if err != nil {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "code": -2,
                "msg":  "Invalid token",
                "data": nil,
            })
        }

        // 3. 检查token是否在黑名单中
        if authService.IsTokenBlacklisted(slug) {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "code": -2,
                "msg":  "Token has been revoked",
                "data": nil,
            })
        }

        // 4. 将slug存入上下文
        c.Locals("slug", slug)
        return c.Next()
    }
}