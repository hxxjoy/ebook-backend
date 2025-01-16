package v1

import (
	"book-fiber/home/controller"
	"book-fiber/home/middleware"
	"book-fiber/system/container"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(router fiber.Router) {

	c := container.GetContainer()

	bookController, err := container.GetController[*controller.BookController](c, "book")
	if err != nil {
		panic(err)
	}
	userController, err := container.GetController[*controller.UserController](c, "user")
	if err != nil {
		panic(err)
	}
	authController, err := container.GetController[*controller.AuthController](c, "auth")
	if err != nil {
		panic(err)
	}
	authService := c.GetAuthService()

	// 认证相关路由
	auth := router.Group("/auth")
	auth.Post("/login", authController.Login)
	auth.Post("/register", authController.Register)
	auth.Post("/refresh", authController.RefreshToken)
	auth.Post("/logout", authController.Logout)
	router.Post("/send-email-code", authController.SendEmailCode)

	// 用户相关路由（需要认证）
	user := router.Group("/user", middleware.Auth(authService))
	user.Get("/profile", userController.GetProfile)
	user.Put("/profile", userController.UpdateProfile)
	user.Put("/password", userController.ChangePassword)

	books := router.Group("/book")

	books.Get("/list", bookController.GetBookList)
	//books.Get("/search/:search", bookController.GetSearchList)
	books.Get("/:slug", bookController.GetOne)
	books.Get("/chapters/:slug", bookController.GetChapters)
	books.Get("/chapter/:slug", bookController.GetChapter)
	bookCategoryController, err := container.GetController[*controller.BookCategoryController](c, "book_category")
	if err != nil {
		panic(err)
	}

	bookCategory := router.Group("/book-category")

	bookCategory.Get("/list", bookCategoryController.GetList)
	bookCategory.Get("/books/:slug/:page", bookCategoryController.GetCategoryBooks)

}
