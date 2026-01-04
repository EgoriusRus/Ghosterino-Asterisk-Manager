package main

import (
	"fmt"
	"log"

	"asterisk-manager/repositories"
	"asterisk-manager/services"
)

func main() {
	fmt.Println("🚀 Asterisk Configuration Generator")
	fmt.Println("====================================")

	// Подключаемся к базе данных
	fmt.Println("\n📡 Подключение к базе данных...")
	connStr := repositories.ConnectionStringFromEnv()
	repos := repositories.InitRepos(connStr)

	// Создаём генератор с выходной папкой results
	generator := services.NewAsteriskGenerator("results")

	// Загружаем данные из БД
	fmt.Println("\n📂 Загрузка данных из базы данных...")
	if err := generator.LoadFromDatabase(repos); err != nil {
		log.Fatalf("❌ Ошибка загрузки данных из БД: %v", err)
	}

	// Генерируем конфигурационные файлы
	fmt.Println("\n⚙️  Генерация конфигурационных файлов...")
	if err := generator.Generate(); err != nil {
		log.Fatalf("❌ Ошибка генерации: %v", err)
	}

	// Выводим статистику
	fmt.Println("\n📊 Статистика:")
	stats := generator.GetStats()
	fmt.Printf("  Всего записей:        %d\n", stats["total"])
	fmt.Printf("  Активных:             %d\n", stats["active"])
	fmt.Printf("  Телефонов T27G:       %d\n", stats["t27"])
	fmt.Printf("  Телефонов T23G:       %d\n", stats["t23"])
	fmt.Printf("  Телефонов Fanvil:     %d\n", stats["fanvil"])
	fmt.Printf("  Cisco/Fax устройств:  %d\n", stats["cisco"])
	fmt.Printf("  С MAC адресами:       %d\n", stats["withMAC"])
	fmt.Printf("  Только локальные:     %d\n", stats["local"])

	fmt.Println("\n✅ Генерация завершена успешно!")
	fmt.Println("📁 Результаты сохранены в папке: results/")
}
