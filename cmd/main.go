package main

import (
	"BankSystem/internal/config"
	"BankSystem/internal/db"
	"BankSystem/internal/handlers"
	"BankSystem/internal/middleware"
	repositories "BankSystem/internal/repositories"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"

	"BankSystem/internal/services"
	account_service "BankSystem/internal/services/account"
)

func main() {
	logger := logrus.New()

	dbCfg := config.LoadDB()
	dsn := db.BuildDSN(dbCfg)
	crypto := config.LoadCrypto()
	runMigrations(dsn)

	ctx := context.Background()

	pool, err := db.New(ctx, dbCfg)
	if err != nil {
		logger.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer pool.Close()
	logger.Info("Подключение к БД успешно установлено")

	dbConnect, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	userRepository := repositories.NewUserRepository(dbConnect)
	accountRepository := repositories.NewAccountRepository(dbConnect)
	cardRepository := repositories.NewCardRepository(dbConnect)

	accountService := account_service.NewAccountService(accountRepository)
	userService := services.NewUserService(userRepository, accountService)
	authService := services.NewAuthService(userRepository)
	cardService := services.NewCardService(dbConnect, cardRepository, accountRepository, accountService, crypto.HMACKey, logger)

	// Инициализация хендлера
	authHandler := handlers.NewAuthHandler(userService)

	// Настройка маршрутов
	r := gin.Default()
	auth := r.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	userHandler := handlers.NewUserHandler(userRepository, authService)
	user := r.Group("/user")
	{
		user.GET("/profile", middleware.AuthMiddleware(), userHandler.GetCurrentUser)
	}

	accountHandler := handlers.NewAccountHandler(accountService, userService, authService)
	account := r.Group("/account")
	{
		account.POST("/create", middleware.AuthMiddleware(), accountHandler.CreateAccount)
		account.GET("/all", middleware.AuthMiddleware(), accountHandler.GetAllAccounts)
		account.POST("/deposit", middleware.AuthMiddleware(), accountHandler.Deposit)
		account.POST("/withdraw", middleware.AuthMiddleware(), accountHandler.Withdraw)
	}

	transfer := r.Group("/transfer")
	{
		transfer.POST("/create", middleware.AuthMiddleware(), accountHandler.Transfer)
	}

	cardHandler := handlers.NewCardHandler(userService, accountService, cardService, authService, accountRepository)
	card := r.Group("/card")
	{
		card.POST("/create", middleware.AuthMiddleware(), cardHandler.CreateCard)
		card.GET("/all", middleware.AuthMiddleware(), cardHandler.GetCards)
		card.POST("/payment", middleware.AuthMiddleware(), cardHandler.PayWithCard)
	}

	// Swagger
	// r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Println("Server is running on :8080")
	err = r.Run(":8080")
	if err != nil {
		log.Println("Server NOT running")
	}
}

func runMigrations(dsn string) {
	m, err := migrate.New("file://migrations", dsn)
	if err != nil {
		logrus.Fatalf("Migration error : %v", err)
		panic(err)
	}

	err = m.Up()

	if err := m.Up(); err != nil {
		if err.Error() == "no change" {
			fmt.Println("No migrations to apply")
			return
		}
		logrus.Error(err)
		panic(err)
	}

	logrus.Info("Migrations applied successfully")
}
