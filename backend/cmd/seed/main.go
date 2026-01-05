package main

import (
	"fmt"
	"log"

	"asterisk-manager/domain"
	"asterisk-manager/repositories"
)

func main() {
	fmt.Println("🌱 SIP Admin Seeds")
	fmt.Println("==================")

	// Подключение к базе данных
	connStr := repositories.ConnectionStringFromEnv()
	repos := repositories.InitRepos(connStr)

	// Выполняем миграции сначала
	fmt.Println("\n🔄 Запуск миграций...")
	if err := repos.MigrateDB(); err != nil {
		log.Fatalf("❌ Ошибка миграции: %v", err)
	}
	fmt.Println("✅ Миграции выполнены успешно")

	// Очищаем таблицы перед seed
	fmt.Println("\n🗑️  Очистка таблиц...")
	if err := cleanTables(repos); err != nil {
		log.Fatalf("❌ Ошибка очистки таблиц: %v", err)
	}
	fmt.Println("✅ Таблицы очищены")

	// Сидим данные
	fmt.Println("\n📥 Заполнение данных...")

	// Локации
	fmt.Println("  → Локации...")
	if err := seedLocations(repos); err != nil {
		log.Fatalf("❌ Ошибка заполнения локаций: %v", err)
	}

	// Устройства
	fmt.Println("  → Устройства...")
	if err := seedDevices(repos); err != nil {
		log.Fatalf("❌ Ошибка заполнения устройств: %v", err)
	}

	// Профили
	fmt.Println("  → Профили...")
	if err := seedProfiles(repos); err != nil {
		log.Fatalf("❌ Ошибка заполнения профилей: %v", err)
	}

	fmt.Println("\n✅ Все данные успешно загружены!")
}

func cleanTables(repos *repositories.Repos) error {
	// Удаляем в правильном порядке (сначала зависимые таблицы)
	if err := repos.DeleteAll(&domain.Profile{}); err != nil {
		return fmt.Errorf("очистка profiles: %w", err)
	}
	if err := repos.DeleteAll(&domain.Device{}); err != nil {
		return fmt.Errorf("очистка devices: %w", err)
	}
	if err := repos.DeleteAll(&domain.Location{}); err != nil {
		return fmt.Errorf("очистка locations: %w", err)
	}

	// Сбрасываем последовательности
	repos.Exec("ALTER SEQUENCE sipadmin.profiles_id_seq RESTART WITH 1")
	repos.Exec("ALTER SEQUENCE sipadmin.locations_id_seq RESTART WITH 1")

	return nil
}

func seedLocations(repos *repositories.Repos) error {
	locations := []domain.Location{
		{Name: "Zags", Server: "10.16.0.102", Subnet: "10.1.191.0/26", VoipVLAN: 5, VLAN: 601},
		{Name: "Sov", Server: "10.16.0.102", Subnet: "10.1.191.0/26", VoipVLAN: 5, VLAN: 6},
		{Name: "Nad3", Server: "10.16.0.102", Subnet: "10.1.17.0/24", VoipVLAN: 65, VLAN: 8},
		{Name: "Ubil1", Server: "10.16.0.102", Subnet: "10.1.80.0/24", VoipVLAN: 4, VLAN: 266},
		{Name: "Ind4", Server: "10.16.0.102", Subnet: "10.1.96.0/25", VoipVLAN: 65, VLAN: 49},
		{Name: "Len15v", Server: "10.16.0.102", Subnet: "10.1.96.128/25", VoipVLAN: 65, VLAN: 10},
	}

	for _, loc := range locations {
		if err := repos.Create(&loc); err != nil {
			return fmt.Errorf("создание локации %s: %w", loc.Name, err)
		}
	}
	return nil
}

func seedDevices(repos *repositories.Repos) error {
	devices := []domain.Device{
		{MAC: "80:5e:c0:b4:42:7c", DeviceModel: domain.DeviceModelYealinkT23G},
		{MAC: "80:5e:c0:18:ab:ac", DeviceModel: domain.DeviceModelYealinkT27G},
		{MAC: "80:5e:c0:e4:a2:fa", DeviceModel: domain.DeviceModelYealinkT27G},
		{MAC: "80:5e:c0:81:d1:6a", DeviceModel: domain.DeviceModelYealinkT27G},
		{MAC: "80:5e:c0:4e:68:c6", DeviceModel: domain.DeviceModelYealinkT27G},
		{MAC: "80:5e:c0:e4:a0:3d", DeviceModel: domain.DeviceModelYealinkT27G},
		{MAC: "80:5e:c0:4e:68:61", DeviceModel: domain.DeviceModelYealinkT27G},
		{MAC: "0c:38:3e:40:3e:52", DeviceModel: domain.DeviceModelCisco},
		{MAC: "1c:e6:c7:99:19:55", DeviceModel: domain.DeviceModelCisco},
		{MAC: "44:db:d2:5e:36:83", DeviceModel: domain.DeviceModelYealinkT27G},
	}

	for _, dev := range devices {
		if err := repos.Create(&dev); err != nil {
			return fmt.Errorf("создание устройства %s: %w", dev.MAC, err)
		}
	}
	return nil
}

func seedProfiles(repos *repositories.Repos) error {
	// Хелперы для указателей
	uintPtr := func(i uint) *uint { return &i }
	intPtr := func(i int) *int { return &i }
	strPtr := func(s string) *string { return &s }

	profiles := []domain.Profile{
		{
			Name:           "Лычкина Елена Владимировна",
			Device:         strPtr("80:5e:c0:b4:42:7c"),
			LocationID:     uintPtr(1),
			InternalNumber: 1119,
			ExternalNumber: "244842",
			RingGroup:      intPtr(6008),
			IsActive:       true,
		},
		{
			Name:           "Кривошеев Олег Викторович",
			Device:         strPtr("80:5e:c0:18:ab:ac"),
			LocationID:     uintPtr(3),
			InternalNumber: 1058,
			ExternalNumber: "947719",
			RingGroup:      intPtr(6018),
			IsActive:       true,
		},
		{
			Name:           "Степанов Александр Александрович",
			Device:         strPtr("80:5e:c0:e4:a2:fa"),
			LocationID:     uintPtr(3),
			InternalNumber: 1048,
			ExternalNumber: "947862",
			RingGroup:      intPtr(6161),
			IsActive:       true,
		},
		{
			Name:           "Дудко Юлия Владимировна",
			Device:         strPtr("80:5e:c0:81:d1:6a"),
			LocationID:     uintPtr(3),
			InternalNumber: 1059,
			ExternalNumber: "947740",
			RingGroup:      intPtr(6039),
			PickupGroup:    intPtr(5),
			IsActive:       true,
		},
		{
			Name:           "Кожокарь Ольга Рафаиловна",
			Device:         strPtr("80:5e:c0:4e:68:c6"),
			LocationID:     uintPtr(3),
			InternalNumber: 1026,
			ExternalNumber: "947740",
			RingGroup:      intPtr(6039),
			PickupGroup:    intPtr(5),
			IsActive:       true,
		},
		{
			Name:           "Забелкин Александр Иванович",
			Device:         strPtr("80:5e:c0:e4:a0:3d"),
			LocationID:     uintPtr(3),
			InternalNumber: 1086,
			ExternalNumber: "947795",
			RingGroup:      intPtr(6094),
			IsActive:       true,
		},
		{
			Name:           "Корсунова Елизавета Васильевна",
			Device:         strPtr("80:5e:c0:4e:68:61"),
			LocationID:     uintPtr(3),
			InternalNumber: 1076,
			ExternalNumber: "947794",
			RingGroup:      intPtr(6093),
			IsActive:       true,
		},
		{
			Name:           "Головина Наталья Ивановна",
			LocationID:     uintPtr(2),
			InternalNumber: 1100,
			ExternalNumber: "947855",
			RingGroup:      intPtr(6154),
			IsActive:       true,
		},
		{
			Name:           "Котенко Анна Владимировна",
			Device:         strPtr("1c:e6:c7:99:19:55"),
			LocationID:     uintPtr(2),
			InternalNumber: 1108,
			ExternalNumber: "947857",
			RingGroup:      intPtr(6156),
			IsActive:       true,
		},
		{
			Name:           "Воловенко Анатолий Сергеевич",
			Device:         strPtr("44:db:d2:5e:36:83"),
			LocationID:     uintPtr(5),
			InternalNumber: 1666,
			ExternalNumber: "930778",
			RingGroup:      intPtr(6392),
			IsActive:       true,
		},
	}

	for _, profile := range profiles {
		if err := repos.Create(&profile); err != nil {
			return fmt.Errorf("создание профиля %s: %w", profile.Name, err)
		}
	}
	return nil
}
