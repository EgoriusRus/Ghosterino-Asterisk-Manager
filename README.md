# Asterisk Manager - Генератор конфигураций для Asterisk

Система управления и генерации конфигурационных файлов для IP-телефонии на базе Asterisk PBX с REST API и веб-интерфейсом управления.

## 🚀 Быстрый старт

```bash
# Полный E2E запуск (frontend + backend + database)
make e2e
```

После запуска:
- 🌐 **Frontend:** http://localhost:3000
- 🔌 **REST API:** http://localhost:8080/api
- 📊 10 профилей сотрудников в базе

## Возможности

- 🖥️ **Веб-интерфейс** на Vue 3 + TypeScript + Tailwind CSS
  - Сотрудники - таблица с пагинацией
  - Устройства - список IP-телефонов
  - Локации - управление сетевыми настройками
- 🌐 **REST API** для управления профилями, устройствами и локациями
- 📊 Загрузка данных из PostgreSQL базы данных
- 🔧 Генерация конфигураций для IP-телефонов (Yealink T27G, T23G, Fanvil, Cisco)
- 📞 Создание SIP-конфигов пользователей
- 🎯 Формирование диалпланов Asterisk
- 🌐 Настройка Cisco маршрутизации
- 💾 Полная интеграция с PostgreSQL через GORM
- 🌱 Database seeder для тестовых данных

## Технологии

- **Frontend:** Vue 3, TypeScript, Tailwind CSS, Vite
- **Backend:** Go 1.24, Fiber v2
- **Database:** PostgreSQL 16 (с network типами: inet, cidr, macaddr)
- **ORM:** GORM
- **Deployment:** Docker, Docker Compose, Nginx

## Архитектура

```
├── frontend/                 # Vue 3 SPA
│   ├── src/
│   │   ├── api/             # API клиент
│   │   ├── views/           # Страницы (Profiles, Devices, Locations)
│   │   ├── layouts/         # Layout с сайдбаром
│   │   ├── router/          # Vue Router
│   │   └── types/           # TypeScript типы
│   ├── Dockerfile           # Multi-stage build (node + nginx)
│   └── nginx.conf           # Nginx с proxy на backend
│
├── backend/
│   ├── cmd/
│   │   ├── generator/       # Генератор конфигов Asterisk
│   │   └── seed/            # Заполнение БД тестовыми данными
│   ├── domain/              # Модели данных (Profile, Device, Location)
│   ├── repositories/        # Слой доступа к данным
│   ├── services/            # Бизнес-логика (генератор конфигов)
│   ├── handlers/            # HTTP handlers для REST API
│   ├── main.go              # Fiber REST API сервер
│   └── Dockerfile           # Multi-stage build
│
├── docker-compose.yml       # Dev окружение
├── docker-compose.prod.yml  # Production окружение
└── Makefile                 # Команды управления
```

## Makefile команды

### Основные команды

```bash
make help           # Показать все доступные команды
make e2e            # 🧪 Полный E2E тест (поднять всё + проверить)
make demo           # 🚀 Полное развертывание (Docker + seed)
make dev            # 💻 Режим разработки (backend в Docker, frontend локально)
```

### Docker команды

```bash
make up             # Поднять все сервисы
make down           # Остановить сервисы
make restart        # Перезапустить сервисы
make status         # Показать статус сервисов
make logs           # Логи всех сервисов
make logs-frontend  # Логи frontend
make logs-backend   # Логи backend
make logs-postgres  # Логи postgres
```

### Сборка и тестирование

```bash
make build          # Пересобрать все образы
make build-frontend # Пересобрать только frontend
make build-backend  # Пересобрать только backend
make api-test       # Протестировать API endpoints
make test           # Запустить Go тесты
make clean          # Полная очистка (контейнеры + volumes)
```

### Утилиты

```bash
make seed           # Заполнить базу тестовыми данными
make generator      # Запустить генератор конфигов Asterisk
make shell-frontend # Войти в shell frontend контейнера
make shell-backend  # Войти в shell backend контейнера
make shell-postgres # Войти в psql консоль
```

### Production команды

```bash
make prod-up        # Запустить продакшн сервисы
make prod-down      # Остановить продакшн сервисы
make prod-restart   # Перезапустить сервисы
make prod-logs      # Показать логи
make prod-build     # Пересобрать образы
make backup         # Создать бэкап БД
```

## Локальная разработка

### Вариант 1: Всё в Docker

```bash
make e2e
# Frontend: http://localhost:3000
# API: http://localhost:8080/api
```

### Вариант 2: Frontend локально (hot reload)

```bash
# Терминал 1: Backend в Docker
make dev

# Терминал 2: Frontend с hot reload
cd frontend
npm run dev
# Frontend: http://localhost:3000 (с proxy на backend)
```

## REST API Endpoints

### Профили (Сотрудники)
- `GET /api/profiles` - Список с пагинацией (`?page=1&perPage=10`)
- `GET /api/profiles/:id` - Один профиль по ID
- `POST /api/profiles` - Создать профиль
- `PUT /api/profiles/:id` - Обновить профиль
- `DELETE /api/profiles/:id` - Удалить профиль

### Устройства
- `GET /api/devices` - Список всех устройств
- `GET /api/devices/:mac` - Устройство по MAC
- `POST /api/devices` - Создать устройство
- `PUT /api/devices/:mac` - Обновить устройство
- `DELETE /api/devices/:mac` - Удалить устройство

### Локации
- `GET /api/locations` - Список всех локаций
- `GET /api/locations/:id` - Локация по ID
- `POST /api/locations` - Создать локацию
- `PUT /api/locations/:id` - Обновить локацию
- `DELETE /api/locations/:id` - Удалить локацию

### Примеры использования API

**Получить профили с пагинацией:**
```bash
curl "http://localhost:8080/api/profiles?page=1&perPage=10"
```

**Ответ:**
```json
{
  "data": [
    {
      "id": 1,
      "name": "Иванов Иван Иванович",
      "email": "ivanov@example.com",
      "internalNumber": 1234,
      "location": {
        "id": 1,
        "name": "Zags",
        "server": "10.16.0.102"
      }
    }
  ],
  "pagination": {
    "total": 10,
    "page": 1,
    "perPage": 10,
    "pages": 1
  }
}
```

**Создать профиль:**
```bash
curl -X POST http://localhost:8080/api/profiles \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Иванов Иван Иванович",
    "email": "ivanov@example.com",
    "device": "80:5e:c0:11:22:33",
    "locationId": 1,
    "internalNumber": 1234
  }'
```

## База данных

### Схема: sipadmin

**locations** - Локации/офисы
| Поле | Тип | Описание |
|------|-----|----------|
| id | serial | Primary Key |
| name | varchar | Название локации |
| server | inet | IP адрес SIP сервера |
| subnet | cidr | Подсеть телефонов |
| voip_vlan | int | VLAN для VoIP |
| vlan | int | VLAN для LAN |

**devices** - IP-телефоны
| Поле | Тип | Описание |
|------|-----|----------|
| mac | macaddr | MAC адрес (Primary Key) |
| device_model | varchar | Модель (Yealink T27G, T23G, Fanvil, Cisco) |

**profiles** - Сотрудники
| Поле | Тип | Описание |
|------|-----|----------|
| id | serial | Primary Key |
| name | varchar | ФИО сотрудника |
| email | varchar | Email |
| device | macaddr | MAC адрес телефона |
| location_id | int | FK на locations |
| internal_number | int | Внутренний номер (уникальный) |
| external_number | varchar | Внешний номер |
| ring_group | int | Группа входящих |
| pickup_group | int | Группа перехвата |
| is_active | boolean | Активность |

## Генератор конфигов Asterisk

```bash
make generator
```

Результаты в `backend/results/`:

- **tftpboot/** - конфиги автопровижининга (по MAC адресу)
- **UsersConf/** - SIP конфигурации пользователей
- **ExtConf/** - файлы диалплана Asterisk
- **CiscoConf.txt** - dial-peer для Cisco

## Переменные окружения

| Переменная | Описание | По умолчанию |
|-----------|----------|--------------|
| `DB_HOST` | Хост PostgreSQL | `localhost` |
| `DB_PORT` | Порт PostgreSQL | `5432` |
| `DB_USER` | Пользователь БД | `postgres` |
| `DB_PASSWORD` | Пароль БД | `postgres` |
| `DB_NAME` | Имя базы данных | `asterisk_manager` |
| `APP_PORT` | Порт API сервера | `8080` |
| `FRONTEND_PORT` | Порт Frontend | `3000` |

## Production Deployment

Подробное руководство: **[docs/DEPLOYMENT.md](docs/DEPLOYMENT.md)**

```bash
# 1. Настройка переменных
cp .env.prod.example .env.prod
nano .env.prod

# 2. Запуск
make prod-up

# 3. Проверка
make prod-logs
```

## Лицензия

MIT
