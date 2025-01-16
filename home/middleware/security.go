package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func SecurityMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        // 添加安全headers
        c.Set("X-Frame-Options", "DENY")
        c.Set("X-Content-Type-Options", "nosniff")
        c.Set("X-XSS-Protection", "1; mode=block")
        c.Set("Referrer-Policy", "strict-origin-when-cross-origin")
        
        return c.Next()
    }
}

// 限制请求频率
func RateLimiter() fiber.Handler {
    return limiter.New(limiter.Config{
        Max:        20,
        Expiration: 1 * time.Minute,
    })
}