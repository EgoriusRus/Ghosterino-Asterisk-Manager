package repositories

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"time"

	"asterisk-manager/domain"

	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// applyPagination applies pagination (LIMIT and OFFSET) to a Gorm query
func applyPagination(query *gorm.DB, pagination *domain.PaginationInput) *gorm.DB {
	if pagination == nil {
		return query
	}

	offset := (pagination.Page - 1) * pagination.PerPage
	return query.Limit(pagination.PerPage).Offset(offset)
}

type Repos struct {
	db *gorm.DB
	// Здесь будут добавлены репозитории для конкретных сущностей
	// Users    *UsersRepo
	// Clients  *ClientsRepo
}

// InitRepos инициализирует подключение к базе данных и создаёт репозитории
func InitRepos(connection string) *Repos {
	db, err := initPostgresConnection(connection)
	if err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	if err = sqlDB.Ping(); err != nil {
		panic(err)
	}

	// Настройка пула соединений
	sqlDB.SetMaxOpenConns(5)
	sqlDB.SetMaxIdleConns(2)
	sqlDB.SetConnMaxLifetime(10 * time.Minute)

	log.Printf("Database connected (driver: postgres), stats: %+v", sqlDB.Stats())

	return &Repos{
		db: db,
		// Инициализация репозиториев будет добавлена позже
		// Users:   NewUsersRepo(db),
		// Clients: NewClientsRepo(db),
	}
}

// initPostgresConnection создаёт подключение к PostgreSQL
func initPostgresConnection(connection string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(connection), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to open Postgres connection")
	}

	return db, nil
}

// MigrateDB выполняет автоматическую миграцию схемы базы данных
func (rs *Repos) MigrateDB() error {
	// Создаём схему sipadmin если её нет
	err := rs.db.Exec("CREATE SCHEMA IF NOT EXISTS sipadmin").Error
	if err != nil {
		return errors.Wrap(err, "failed to create schema")
	}

	// Автомиграция таблиц
	err = rs.db.AutoMigrate(
		&domain.Location{},
		&domain.Device{},
		&domain.Profile{},
	)
	if err != nil {
		return errors.WithStack(err)
	}

	// Создаём индексы
	err = rs.createIndexes()
	if err != nil {
		return errors.Wrap(err, "failed to create indexes")
	}

	return nil
}

// createIndexes создаёт дополнительные индексы
func (rs *Repos) createIndexes() error {
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_profiles_location ON sipadmin.profiles(location_id)",
		"CREATE INDEX IF NOT EXISTS idx_profiles_device ON sipadmin.profiles(device)",
		"CREATE INDEX IF NOT EXISTS idx_profiles_internal ON sipadmin.profiles(internal_number)",
		"CREATE INDEX IF NOT EXISTS idx_profiles_external ON sipadmin.profiles(external_number)",
	}

	for _, idx := range indexes {
		if err := rs.db.Exec(idx).Error; err != nil {
			return errors.Wrapf(err, "failed to create index: %s", idx)
		}
	}

	return nil
}

// Save сохраняет объект (создаёт или обновляет)
func (rs *Repos) Save(object interface{}) error {
	rt := reflect.TypeOf(object)
	switch rt.Kind() {
	case reflect.Slice:
		vo := reflect.ValueOf(object)
		if vo.Len() == 0 {
			return nil
		}
	}

	return rs.db.Save(object).Error
}

// Delete удаляет объект из базы данных
func (rs *Repos) Delete(object interface{}) error {
	rt := reflect.TypeOf(object)
	switch rt.Kind() {
	case reflect.Slice:
		vo := reflect.ValueOf(object)
		if vo.Len() == 0 {
			return nil
		}
	}

	return rs.db.Delete(object).Error
}

// FindAll находит все записи
func (rs *Repos) FindAll(dest interface{}) error {
	return rs.db.Find(dest).Error
}

// FindByID находит запись по ID
func (rs *Repos) FindByID(dest interface{}, id interface{}) error {
	return rs.db.First(dest, id).Error
}

// FindOne находит одну запись по условию
func (rs *Repos) FindOne(dest interface{}, condition string, args ...interface{}) error {
	return rs.db.Where(condition, args...).First(dest).Error
}

// FindProfilesWithLocations находит профили с джойном к локациям
func (rs *Repos) FindProfilesWithLocations(isActive *bool, pagination *domain.PaginationInput) ([]domain.ProfileWithLocation, int64, error) {
	var profiles []domain.ProfileWithLocation
	var total int64

	query := rs.db.Table("sipadmin.profiles AS p").
		Select(`
			p.id, p.name, p.email, p.device, p.location_id, p.internal_number,
			p.external_number, p.ring_group, p.pickup_group, p.is_active,
			p.created_at, p.updated_at,
			l.name AS location_name, l.server, l.subnet, l.voip_vlan, l.vlan
		`).
		Joins("LEFT JOIN sipadmin.locations AS l ON p.location_id = l.id")

	if isActive != nil {
		query = query.Where("p.is_active = ?", *isActive)
	}

	// Get total count before pagination
	countQuery := rs.db.Table("sipadmin.profiles AS p")
	if isActive != nil {
		countQuery = countQuery.Where("p.is_active = ?", *isActive)
	}
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	query = applyPagination(query, pagination)

	err := query.Scan(&profiles).Error
	return profiles, total, err
}

// Exec выполняет raw SQL запрос
func (rs *Repos) Exec(sql string) error {
	return rs.db.Exec(sql).Error
}

// ConnectionString формирует строку подключения к PostgreSQL
func ConnectionString(host, port, user, password, dbname string) string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)
}

// ConnectionStringFromEnv формирует строку подключения из переменных окружения
func ConnectionStringFromEnv() string {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "postgres")
	dbname := getEnv("DB_NAME", "asterisk_manager")

	return ConnectionString(host, port, user, password, dbname)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
