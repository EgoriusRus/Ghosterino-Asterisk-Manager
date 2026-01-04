package main

import (
	"fmt"
	"log"

	"romchek-asteriska/repositories"
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
	queries := []string{
		"DELETE FROM sipadmin.profiles",
		"DELETE FROM sipadmin.devices",
		"DELETE FROM sipadmin.locations",
		"ALTER SEQUENCE sipadmin.profiles_id_seq RESTART WITH 1",
		"ALTER SEQUENCE sipadmin.locations_id_seq RESTART WITH 1",
	}

	for _, query := range queries {
		if err := repos.Exec(query); err != nil {
			return err
		}
	}

	return nil
}

func seedLocations(repos *repositories.Repos) error {
	sql := `
		INSERT INTO sipadmin.locations (name, server, subnet, voip_vlan, vlan) VALUES
		('Zags', '10.16.0.102', '10.1.191.0/26', 5, 601),
		('Sov', '10.16.0.102', '10.1.191.0/26', 5, 6),
		('Nad3', '10.16.0.102', '10.1.17.0/24', 65, 8),
		('Ubil1', '10.16.0.102', '10.1.80.0/24', 4, 266),
		('Ind4', '10.16.0.102', '10.1.96.0/25', 65, 49),
		('Len15v', '10.16.0.102', '10.1.96.128/25', 65, 10)
	`
	return repos.Exec(sql)
}

func seedDevices(repos *repositories.Repos) error {
	sql := `
		INSERT INTO sipadmin.devices (mac, device_model) VALUES
		('80:5e:c0:b4:42:7c', 'Yealink T23G'),
		('80:5e:c0:18:ab:ac', 'Yealink T27G'),
		('80:5e:c0:e4:a2:fa', 'Yealink T27G'),
		('80:5e:c0:81:d1:6a', 'Yealink T27G'),
		('80:5e:c0:4e:68:c6', 'Yealink T27G'),
		('80:5e:c0:e4:a0:3d', 'Yealink T27G'),
		('80:5e:c0:4e:68:61', 'Yealink T27G'),
		('0c:38:3e:40:3e:52', 'Cisco'),
		('1c:e6:c7:99:19:55', 'Cisco'),
		('44:db:d2:5e:36:83', 'Yealink T27G')
	`
	return repos.Exec(sql)
}

func seedProfiles(repos *repositories.Repos) error {
	sql := `
		INSERT INTO sipadmin.profiles (name, email, device, location_id, internal_number, external_number, ring_group, pickup_group, is_active) VALUES
		('Лычкина Елена Владимировна', NULL, '80:5e:c0:b4:42:7c', 1, 1119, '244842', 6008, NULL, true),
		('Кривошеев Олег Викторович', NULL, '80:5e:c0:18:ab:ac', 3, 1058, '947719', 6018, NULL, true),
		('Степанов Александр Александрович', NULL, '80:5e:c0:e4:a2:fa', 3, 1048, '947862', 6161, NULL, true),
		('Дудко Юлия Владимировна', NULL, '80:5e:c0:81:d1:6a', 3, 1059, '947740', 6039, 5, true),
		('Кожокарь Ольга Рафаиловна', NULL, '80:5e:c0:4e:68:c6', 3, 1026, '947740', 6039, 5, true),
		('Забелкин Александр Иванович', NULL, '80:5e:c0:e4:a0:3d', 3, 1086, '947795', 6094, NULL, true),
		('Корсунова Елизавета Васильевна', NULL, '80:5e:c0:4e:68:61', 3, 1076, '947794', 6093, NULL, true),
		('Головина Наталья Ивановна', NULL, NULL, 2, 1100, '947855', 6154, NULL, true),
		('Котенко Анна Владимировна', NULL, '1c:e6:c7:99:19:55', 2, 1108, '947857', 6156, NULL, true),
		('Воловенко Анатолий Сергеевич', NULL, '44:db:d2:5e:36:83', 5, 1666, '930778', 6392, NULL, true)
	`
	return repos.Exec(sql)
}
