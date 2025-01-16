package container

import (
	"book-fiber/home/controller"
	"book-fiber/home/repository"
	"book-fiber/home/service"
	"book-fiber/system/cache"
	"book-fiber/system/config"
	baseController "book-fiber/system/controller"
	"book-fiber/system/database"
	"book-fiber/system/helper"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	once     sync.Once
	instance *Container
)

type Container struct {
	db           *gorm.DB
	cache        *cache.Cache
	repositories map[string]interface{}
	services     map[string]interface{}
	controllers  map[string]interface{}
	config       *config.Config
	baseCtl      baseController.BaseController
	http         helper.HttpHelper
	redis        *redis.Client
}

// Singleton pattern to get the container instance
func GetContainer() *Container {
	once.Do(func() {
		redisAddr := fmt.Sprintf("%s:%d",config.C.Redis.Host, config.C.Redis.Port )
		instance = &Container{
			db:           database.GetDB(),
			cache:        cache.NewCache(config.C.Cache.Enabled, time.Duration(config.C.Cache.TTL)*time.Minute, time.Duration(config.C.Cache.Clearup)*time.Minute, config.C.Cache.MaxSize),
			repositories: make(map[string]interface{}),
			services:     make(map[string]interface{}),
			controllers:  make(map[string]interface{}),
			config:       config.C,
			http:         helper.NewHttpHelper(),
			redis:       redis.NewClient(&redis.Options{
				Addr:     redisAddr,
				Password: config.C.Redis.Password,
				DB:       config.C.Redis.DB,
			}),
		}
		instance.baseCtl = baseController.NewBaseController(instance.http)
	})
	return instance
}

func (c *Container) GetHttp() helper.HttpHelper {
	return c.http
}

func (c *Container) GetCache() *cache.Cache {
	return c.cache
}

func (c *Container) GetConfig() *config.Config {
	return c.config
}

func GetController[T any](c *Container, name string) (T, error) {
	var zero T

	// 检查是否已存在
	if ctrl, exists := c.controllers[name]; exists {
		if result, ok := ctrl.(T); ok {
			return result, nil
		}
		return zero, fmt.Errorf("invalid controller type for %s", name)
	}

	// 初始化新的 controller
	switch name {
	case "book":
		ctrl := c.GetBookController()
		// 直接将初始化的 controller 转换为请求的类型 T
		if result, ok := any(ctrl).(T); ok {
			return result, nil
		}
	case "auth":
		ctrl := c.GetAuthController()
		// 直接将初始化的 controller 转换为请求的类型 T
		if result, ok := any(ctrl).(T); ok {
			return result, nil
		}
	case "user":
		ctrl := c.GetUserController()
		// 直接将初始化的 controller 转换为请求的类型 T
		if result, ok := any(ctrl).(T); ok {
			return result, nil
		}
	case "book_category":
		ctrl := c.GetBookCategoryController()
		// 直接将初始化的 controller 转换为请求的类型 T
		if result, ok := any(ctrl).(T); ok {
			return result, nil
		}
	}

	return zero, fmt.Errorf("controller %s not found", name)
}

func (c *Container) GetBookRepository() *repository.BookRepository {
	if ctrl, exists := c.repositories["book"]; exists {
		return ctrl.(*repository.BookRepository)
	}

	repo := repository.NewBookRepository(c.db, c.cache, c.config)
	c.repositories["book"] = repo
	return repo
}
func (c *Container) GetBookService() *service.BookService {
	if ctrl, exists := c.services["book"]; exists {
		return ctrl.(*service.BookService)
	}

	repo := c.GetBookRepository()
	svc := service.NewBookService(repo)
	c.services["book"] = svc
	return svc
}

func (c *Container) GetBookCategoryRepository() *repository.BookCategoryRepository {
	if ctrl, exists := c.repositories["book_category"]; exists {
		return ctrl.(*repository.BookCategoryRepository)
	}

	repo := repository.NewBookCategoryRepository(c.db, c.cache, c.config)
	c.repositories["book_category"] = repo
	return repo
}
func (c *Container) GetBookCategoryService() *service.BookCategoryService {
	if ctrl, exists := c.services["book_category"]; exists {
		return ctrl.(*service.BookCategoryService)
	}

	repo := c.GetBookCategoryRepository()
	svc := service.NewBookCategoryService(repo)
	c.services["book_category"] = svc
	return svc
}
func (c *Container) GetBookCategoryController() *controller.BookCategoryController {
	if ctrl, exists := c.controllers["book_category"]; exists {
		return ctrl.(*controller.BookCategoryController)
	}

	svc := c.GetBookCategoryService()
	ctrl := controller.NewBookCategoryController(c.baseCtl, svc)
	c.controllers["book_category"] = ctrl
	return ctrl
}
func (c *Container) GetUserRepository() *repository.UserRepository {
	if ctrl, exists := c.repositories["user"]; exists {
		return ctrl.(*repository.UserRepository)
	}

	repo := repository.NewUserRepository(c.db, c.cache, c.config)
	c.repositories["user"] = repo
	return repo
}
func (c *Container) GetTokenBlacklistRepository() repository.TokenBlacklistRepository {
	if ctrl, exists := c.repositories["token_blacklist"]; exists {
		return ctrl.(repository.TokenBlacklistRepository)
	}

	repo := repository.NewTokenBlacklistRepository(c.db)
	c.repositories["token_blacklist"] = repo
	return repo
}
func (c *Container) GetUserService() *service.UserService {
	if ctrl, exists := c.services["user"]; exists {
		return ctrl.(*service.UserService)
	}

	repo := c.GetUserRepository()
	svc := service.NewUserService(repo)
	c.services["user"] = svc
	return svc
}
func (c *Container) GetAuthService() *service.AuthService {
	if ctrl, exists := c.services["auth"]; exists {
		return ctrl.(*service.AuthService)
	}

	repo := c.GetUserRepository()
	tokenBlacklistRepo := c.GetTokenBlacklistRepository()
	svc := service.NewAuthService(repo,tokenBlacklistRepo)
	c.services["auth"] = svc
	return svc
}
func (c *Container) GetEmailService() service.EmailService {
	if ctrl, exists := c.services["email"]; exists {
		return ctrl.(service.EmailService)
	}

	svc := service.NewEmailService(c.config,c.redis,c.db)
	c.services["email"] = svc
	return svc
}
func (c *Container) GetUserController() *controller.UserController {
	if ctrl, exists := c.controllers["user"]; exists {
		return ctrl.(*controller.UserController)
	}

	svc := c.GetUserService()
	emailSvc := c.GetEmailService()
	ctrl := controller.NewUserController(c.baseCtl, svc,emailSvc)
	c.controllers["user"] = ctrl
	return ctrl
}
func (c *Container) GetBookController() *controller.BookController {
	if ctrl, exists := c.controllers["book"]; exists {
		return ctrl.(*controller.BookController)
	}
	svc := c.GetBookService()
	ctrl := controller.NewBookController(c.baseCtl, svc)

	c.controllers["book"] = ctrl
	return ctrl
}
func (c *Container) GetAuthController() *controller.AuthController {
	if ctrl, exists := c.controllers["auth"]; exists {
		return ctrl.(*controller.AuthController)
	}
	svc := c.GetAuthService()
	userSvc := c.GetUserService()
	emailSvc := c.GetEmailService()
	ctrl := controller.NewAuthController(c.baseCtl, userSvc,svc,emailSvc)

	c.controllers["auth"] = ctrl
	return ctrl
}