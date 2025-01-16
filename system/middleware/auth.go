// middleware/auth.go
package middleware

import "github.com/gofiber/fiber/v2"

func AdminAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Cookies("admin_token")
		if token == "" {
			return c.Redirect("/login")
		}
		// 验证token
		return c.Next()
	}
}

func APIAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Cookies("api_token")
		if token == "" {
			return c.Redirect("/login")
		}
		// 验证token
		return c.Next()
	}
}
