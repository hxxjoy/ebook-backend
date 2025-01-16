package main

import (
	"fmt"
	"os"
	"runtime/debug"
	"time"

	"book-fiber/routes"
	"book-fiber/system/config"
	"book-fiber/system/database"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var log = logrus.New()

func main() {
	// 1. 加载配置文件
	config.MustLoad("configs/app.yaml")
	// 2. 创建 Fiber 实例
	app := fiber.New(fiber.Config{
		AppName:      config.C.App.Name,
		ReadTimeout:  config.C.Server.ReadTimeout,
		WriteTimeout: config.C.Server.WriteTimeout,
		//ErrorHandler: customErrorHandler,
	})

	// 3. 使用中间件
	app.Use(logger.New())
	app.Use(cors.New())
	app.Use(recoveryMiddleware)
	log.SetFormatter(&logrus.JSONFormatter{})
	file, err := os.OpenFile("tmp/error.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(file)
	}
	// 4. 初始化数据库连接
	if err := initDB(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 5. 注册路由
	routes.SetupRoutes(app)

	// 6. 启动服务器
	addr := fmt.Sprintf("%s:%d", config.C.Server.Host, config.C.Server.Port)
	if err := app.Listen(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func initDB() error {
	dsn := config.C.Database.GetDSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxOpenConns(config.C.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(config.C.Database.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(config.C.Database.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(config.C.Database.ConnMaxIdleTime)

	// 将 db 设置到 database 包中
	database.SetDB(db)
	return nil
}

func recoveryMiddleware(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			// 获取 panic 信息和堆栈
			stack := debug.Stack()
			detail := ""
			if isDevelopment() {
				detail = string(stack)
			}
			// 获取请求参数
			var queryParams string
			if c.Context().QueryArgs() != nil {
				queryParams = string(c.Context().QueryArgs().String())
			}

			// 获取请求体
			var body string
			if len(c.Body()) > 0 {
				body = string(c.Body())
			}
			log.WithFields(logrus.Fields{
				"error":  r,
				"stack":  detail,
				"path":   c.Path(),
				"url":    c.OriginalURL(),
				"params": queryParams,
				"body":   body,
			}).Error("Panic recovered")

			// 返回统一的错误响应
			c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":  "Internal Server Error",
				"detail": detail,
			})
		}
	}()
	return c.Next() // 执行下一个中间件
}

func isDevelopment() bool {
	return os.Getenv("GO_ENV") == "development"
}

func LogMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		// 记录请求日志
		log.Printf(
			"[%s] %s %s %v",
			c.Method(),
			c.Path(),
			time.Since(start),
			err,
		)

		return err
	}
}
