package server

import (
	AuthHandler "Market_backend/internal/auth/handler"
	AuthRepository "Market_backend/internal/auth/repository"
	AuthRouter "Market_backend/internal/auth/router"
	AuthService "Market_backend/internal/auth/service"

	PaymentHandler "Market_backend/internal/payment/handler"
	PaymentRepo "Market_backend/internal/payment/repository"
	PaymentRouter "Market_backend/internal/payment/router"
	PaymentService "Market_backend/internal/payment/service"

	ProductHandler "Market_backend/internal/product/handler"
	ProductRepository "Market_backend/internal/product/repository"
	ProductRouter "Market_backend/internal/product/router"
	ProductService "Market_backend/internal/product/service"

	UserHandler "Market_backend/internal/user/handler"
	UserRepository "Market_backend/internal/user/repository"
	UserRouter "Market_backend/internal/user/router"
	UserService "Market_backend/internal/user/service"

	CartHandler "Market_backend/internal/cart/handler"
	CartRepository "Market_backend/internal/cart/repository"
	CartRouter "Market_backend/internal/cart/router"
	CartService "Market_backend/internal/cart/service"

	OrderHandler "Market_backend/internal/order/handler"
	OrderRepository "Market_backend/internal/order/repository"
	OrderRouter "Market_backend/internal/order/router"
	OrderService "Market_backend/internal/order/service"

	MessageHandler "Market_backend/internal/messages/handler"
	MessageRepository "Market_backend/internal/messages/repository"
	MessageRouter "Market_backend/internal/messages/router"
	MessageService "Market_backend/internal/messages/service"

	"Market_backend/internal/storage"
	"log"

	"Market_backend/internal/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func Start() {
	app := fiber.New()
	miniStorage, err := storage.NewMinioStorage()
	if err != nil {
		log.Fatal(err)
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins:     config.AllowedOrigins, // Обрати внимание на запятую и пробел
		AllowMethods:     "GET, POST, PUT, PATCH, DELETE, OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
		ExposeHeaders:    "Set-Cookie", // Важно для cookies!
	}))

	procRepo := ProductRepository.NewProcessorRepository()
	procService := ProductService.NewProcessorService(procRepo, miniStorage)
	procHandler := ProductHandler.NewProcessorHandler(procService)

	ProductRouter.RegisterProcessorRouter(app, procHandler)

	flashdriveRepo := ProductRepository.NewFlashDriveRepository()
	flashdriveService := ProductService.NewFlashDriveService(flashdriveRepo, miniStorage)
	flashdriveHandler := ProductHandler.NewFlashDriveHandler(flashdriveService)

	ProductRouter.RegisterFlashDriverRouter(app, flashdriveHandler)

	cartRepo := CartRepository.NewCartRepository()
	cartService := CartService.NewCartService(cartRepo, procRepo, flashdriveRepo)
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

	orderRepo := OrderRepository.NewOrderRepository()
	orderService := OrderService.NewOrderService(orderRepo, cartRepo, cartService, procService, flashdriveService)
	orderHandler := OrderHandler.NewOrderHandler(orderService)

	OrderRouter.RegisterOrderRouter(app, orderHandler)

	paymentRepo := PaymentRepo.NewPaymentRepository()
	paymentService := PaymentService.NewPaymentService(paymentRepo, orderRepo)
	paymentHandler := PaymentHandler.NewPaymentHandler(paymentService, orderRepo)

	PaymentRouter.RegisterPaymentRouter(app, paymentHandler)

	messageRepo := MessageRepository.NewMessageRepository()
	messageService := MessageService.NewMessageService(messageRepo)
	messageHandler := MessageHandler.NewMessageHandler(messageService)

	MessageRouter.RegisterMessageRoutes(app, messageHandler)

	if config.AppPort != "" {
		err = app.Listen(":" + config.AppPort)
	} else {
		err = app.Listen(":3000")
	}
	if err != nil {
		log.Fatal(err)
	}
}
