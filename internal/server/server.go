package server

import (
	AuthHandler "eduVix_backend/internal/auth/handler"
	AuthRepository "eduVix_backend/internal/auth/repository"
	AuthRouter "eduVix_backend/internal/auth/router"
	AuthService "eduVix_backend/internal/auth/service"

	ProductHandler "eduVix_backend/internal/product/handler"
	ProductRepository "eduVix_backend/internal/product/repository"
	ProductRouter "eduVix_backend/internal/product/router"
	ProductService "eduVix_backend/internal/product/service"

	UserHandler "eduVix_backend/internal/user/handler"
	UserRepository "eduVix_backend/internal/user/repository"
	UserRouter "eduVix_backend/internal/user/router"
	UserService "eduVix_backend/internal/user/service"

	CartHandler "eduVix_backend/internal/cart/handler"
	CartRepository "eduVix_backend/internal/cart/repository"
	CartRouter "eduVix_backend/internal/cart/router"
	CartService "eduVix_backend/internal/cart/service"

	"eduVix_backend/internal/storage"
	"log"

	"github.com/gofiber/fiber/v2"
)

func Start() {
	app := fiber.New()
	storage, err := storage.NewMinioStorage()
	if err != nil {
		log.Fatal(err)
	}

	procRepo := ProductRepository.NewProcessorRepository()
	procService := ProductService.NewProcessorService(procRepo, storage)
	procHandler := ProductHandler.NewProcessorHandler(procService)

	ProductRouter.RegisterProcessorRouter(app, procHandler)

	cartRepo := CartRepository.NewCartRepository()
	cartService := CartService.NewCartService(cartRepo, procRepo)
	cartHandler := CartHandler.NewCartHandler(cartService)

	CartRouter.RegisterCartRouter(app, cartHandler)

	userRepo := UserRepository.NewUserRepository()
	userService := UserService.NewUserService(userRepo)
	userHandler := UserHandler.NewUserHandler(userService)

	UserRouter.RegisterUserRoutes(app, userHandler)

	authRepo := AuthRepository.NewAuthRepository()
	authService := AuthService.NewAuthService(authRepo, cartRepo)
	authHandler := AuthHandler.NewAuthHandler(authService)

	AuthRouter.RegisterAuthRouter(app, authHandler)

	app.Listen(":3000")
}
