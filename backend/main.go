package main

import (
	"fmt"
	"log"
	"os"

	"asterisk-manager/handlers"
	"asterisk-manager/repositories"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	fmt.Println("🚀 Asterisk Manager API")
	fmt.Println("=======================")

	// Подключение к базе данных
	fmt.Println("\n📡 Подключение к базе данных...")
	connStr := repositories.ConnectionStringFromEnv()
	repos := repositories.InitRepos(connStr)

	// Выполняем миграции
	fmt.Println("\n🔄 Запуск миграций...")
	if err := repos.MigrateDB(); err != nil {
		log.Fatalf("❌ Ошибка миграции: %v", err)
	}
	fmt.Println("✅ Миграции выполнены успешно")

	// Создаем дефолтного админа
	fmt.Println("\n👤 Проверка пользователя admin...")
	if err := repos.CreateDefaultAdmin(); err != nil {
		log.Fatalf("❌ Ошибка создания админа: %v", err)
	}
	fmt.Println("✅ Пользователь admin готов")

	// Создаём handler
	h := handlers.NewHandler(repos)
	authHandler := handlers.NewAuthHandler(h)

	// Создаём Fiber приложение
	app := fiber.New(fiber.Config{
		ErrorHandler: h.ErrorHandler,
		AppName:      "Asterisk Manager API",
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "${time} | ${status} | ${latency} | ${method} ${path}\n",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Инициализируем роуты
	initRoutes(app, h, authHandler)

	// Запускаем сервер
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("\n🌐 Сервер запущен на порту %s\n", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("❌ Ошибка запуска сервера: %v", err)
	}
}
