# Инструкция по развертыванию демо на Windows

Пошаговое руководство для запуска Asterisk Manager на Windows с помощью Docker.

## Системные требования

- **Windows 10/11** (64-bit)
- **Docker Desktop for Windows** 4.0+
- **4 GB RAM** (минимум)
- **10 GB** свободного места на диске

## 1. Установка необходимого ПО

### 1.1. Установка Docker Desktop

1. Скачайте Docker Desktop с официального сайта:
   - https://www.docker.com/products/docker-desktop/

2. Запустите установщик `Docker Desktop Installer.exe`

3. Следуйте инструкциям установщика:
   - Включите опцию "Use WSL 2 instead of Hyper-V" (рекомендуется)
   - Нажмите "Ok" и дождитесь завершения установки

4. Перезагрузите компьютер

5. После перезагрузки запустите Docker Desktop:
   - Найдите Docker Desktop в меню Пуск
   - Дождитесь, пока Docker Engine запустится (индикатор в трее станет зеленым)

6. Проверьте установку, открыв PowerShell или Command Prompt:
   ```powershell
   docker --version
   docker-compose --version
   ```

   Должны отобразиться версии Docker и Docker Compose.

## 2. Распаковка архива проекта

1. Распакуйте полученный ZIP-архив `asterisk-manager.zip` в удобное место, например:
   - `C:\Projects\asterisk-manager`

2. Откройте PowerShell или Command Prompt

3. Перейдите в директорию проекта:
   ```powershell
   cd C:\Projects\asterisk-manager
   ```

4. Проверьте содержимое:
   ```powershell
   dir
   ```

   Вы должны увидеть папки: `backend`, `frontend`, `docs` и файлы `docker-compose.yml`, `Makefile`, `README.md`

## 3. Настройка переменных окружения (опционально)

Проект работает с настройками по умолчанию, но при необходимости вы можете их изменить:

1. Создайте файл `.env` в корне проекта:
   ```powershell
   copy .env.example .env
   ```

2. Отредактируйте `.env` в текстовом редакторе (например, Notepad):
   ```env
   DB_HOST=postgres
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=postgres
   DB_NAME=asterisk_manager
   APP_PORT=8080
   FRONTEND_PORT=3000
   ```

## 4. Запуск приложения

### Вариант 1: Быстрый запуск с Docker Compose (рекомендуется)

1. Убедитесь, что Docker Desktop запущен

2. В PowerShell или Command Prompt выполните:
   ```powershell
   docker-compose up -d --build
   ```

   Флаги:
   - `-d` - запуск в фоновом режиме
   - `--build` - пересборка образов

3. Дождитесь завершения сборки и запуска (первый запуск может занять 5-10 минут)

4. Проверьте статус контейнеров:
   ```powershell
   docker-compose ps
   ```

   Все три сервиса должны быть в статусе `Up`:
   - `asterisk-postgres`
   - `asterisk-backend`
   - `asterisk-frontend`

5. Заполните базу данных тестовыми данными:
   ```powershell
   docker-compose exec backend /app/seed
   ```

6. Откройте в браузере:
   - **Frontend:** http://localhost:3000
   - **API:** http://localhost:8080/api

### Вариант 2: Использование Make (требуется установка Make для Windows)

Если у вас установлен Make для Windows (можно установить через Chocolatey или MinGW):

```powershell
make demo
```

Эта команда автоматически:
- Соберет все образы
- Запустит сервисы
- Заполнит базу тестовыми данными

## 5. Проверка работоспособности

### 5.1. Проверка Frontend

1. Откройте браузер и перейдите на http://localhost:3000

2. Вы должны увидеть веб-интерфейс с меню:
   - Сотрудники
   - Устройства
   - Локации

3. Перейдите в раздел "Сотрудники" - должна отобразиться таблица с 10 сотрудниками

### 5.2. Проверка API

Используйте PowerShell для проверки API:

```powershell
# Получить список сотрудников
Invoke-WebRequest -Uri "http://localhost:8080/api/profiles?page=1&perPage=10" | Select-Object -Expand Content

# Получить список устройств
Invoke-WebRequest -Uri "http://localhost:8080/api/devices" | Select-Object -Expand Content

# Получить список локаций
Invoke-WebRequest -Uri "http://localhost:8080/api/locations" | Select-Object -Expand Content
```

Или используйте браузер:
- http://localhost:8080/api/profiles
- http://localhost:8080/api/devices
- http://localhost:8080/api/locations

### 5.3. Проверка базы данных

Подключитесь к PostgreSQL:

```powershell
docker-compose exec postgres psql -U postgres -d asterisk_manager
```

Выполните запрос:
```sql
SELECT COUNT(*) FROM sipadmin.profiles;
```

Должно вывести `10` (количество сотрудников).

Выход из psql: `\q`

## 6. Управление сервисами

### Остановка всех сервисов
```powershell
docker-compose down
```

### Запуск остановленных сервисов
```powershell
docker-compose up -d
```

### Перезапуск сервисов
```powershell
docker-compose restart
```

### Просмотр логов

Все сервисы:
```powershell
docker-compose logs -f
```

Конкретный сервис:
```powershell
docker-compose logs -f backend
docker-compose logs -f frontend
docker-compose logs -f postgres
```

Выход из логов: `Ctrl+C`

### Пересборка образов

После изменения кода:
```powershell
docker-compose up -d --build
```

### Полная очистка

Остановка и удаление всех контейнеров, сетей и volumes:
```powershell
docker-compose down -v
```

**Внимание:** Это удалит все данные из базы данных!

## 7. Работа с генератором конфигураций Asterisk

### Запуск генератора

```powershell
docker-compose exec backend /app/generator
```

### Просмотр результатов

Результаты генерации сохраняются в `backend/results/`:

```powershell
# Просмотр структуры
dir backend\results

# Просмотр конфига устройства
type backend\results\tftpboot\<mac-адрес>.cfg

# Просмотр SIP конфигурации
type backend\results\UsersConf\<внутренний-номер>.conf
```

## 8. Устранение неполадок

### Docker Desktop не запускается

1. Убедитесь, что WSL 2 установлен:
   ```powershell
   wsl --update
   ```

2. Включите виртуализацию в BIOS (если отключена)

3. Перезапустите Docker Desktop от имени администратора

### Порты заняты

Если порты 3000, 8080 или 5432 заняты другими приложениями:

1. Измените порты в `.env`:
   ```env
   FRONTEND_PORT=3001
   APP_PORT=8081
   DB_PORT=5433
   ```

2. Перезапустите сервисы:
   ```powershell
   docker-compose down
   docker-compose up -d
   ```

### Ошибка "No space left on device"

Очистите неиспользуемые Docker образы и контейнеры:

```powershell
docker system prune -a
```

### База данных не инициализируется

1. Остановите сервисы и удалите volumes:
   ```powershell
   docker-compose down -v
   ```

2. Запустите заново:
   ```powershell
   docker-compose up -d
   ```

3. Дождитесь инициализации БД (30-60 секунд)

4. Заполните тестовыми данными:
   ```powershell
   docker-compose exec backend /app/seed
   ```

### Проверка здоровья контейнеров

```powershell
# Статус всех контейнеров
docker-compose ps

# Детальная информация о контейнере
docker inspect asterisk-backend
docker inspect asterisk-frontend
docker inspect asterisk-postgres

# Логи конкретного контейнера
docker logs asterisk-backend
```

### Backend не может подключиться к БД

1. Проверьте, что PostgreSQL запущен:
   ```powershell
   docker-compose ps postgres
   ```

2. Проверьте health check:
   ```powershell
   docker inspect asterisk-postgres | Select-String "Health"
   ```

3. Перезапустите backend:
   ```powershell
   docker-compose restart backend
   ```

## 9. Разработка с hot reload

Для разработки frontend с автоматической перезагрузкой при изменениях:

### Терминал 1: Backend в Docker
```powershell
docker-compose up -d backend postgres
```

### Терминал 2: Frontend локально

1. Установите Node.js (если не установлен):
   - Скачайте с https://nodejs.org/ (LTS версия)
   - Установите

2. Перейдите в папку frontend:
   ```powershell
   cd frontend
   ```

3. Установите зависимости:
   ```powershell
   npm install
   ```

4. Запустите dev сервер:
   ```powershell
   npm run dev
   ```

5. Откройте http://localhost:5173 (Vite dev server)

## 10. Дополнительные команды

### Вход в контейнер

Backend:
```powershell
docker-compose exec backend sh
```

Frontend:
```powershell
docker-compose exec frontend sh
```

PostgreSQL:
```powershell
docker-compose exec postgres psql -U postgres -d asterisk_manager
```

### Резервное копирование базы данных

```powershell
docker-compose exec postgres pg_dump -U postgres asterisk_manager > backup.sql
```

### Восстановление из бэкапа

```powershell
Get-Content backup.sql | docker-compose exec -T postgres psql -U postgres -d asterisk_manager
```

## 11. Архитектура и порты

После успешного запуска у вас будут работать:

| Сервис | Порт | URL | Описание |
|--------|------|-----|----------|
| Frontend | 3000 | http://localhost:3000 | Vue 3 веб-интерфейс |
| Backend API | 8080 | http://localhost:8080/api | REST API (Fiber) |
| PostgreSQL | 5432 | localhost:5432 | База данных |

### Сетевая архитектура

Все контейнеры находятся в одной Docker сети `asterisk-network`:
- Frontend общается с Backend через Nginx proxy
- Backend подключается к PostgreSQL по hostname `postgres`
- Внешний доступ только к портам 3000 и 8080

## 12. Тестовые данные

После выполнения `docker-compose exec backend /app/seed` будут созданы:

- **10 сотрудников** с внутренними номерами 1000-1009
- **3 локации** (Zags, Sev, Gag)
- **10 устройств** (Yealink T27G, T23G, Fanvil, Cisco)

Вы можете просмотреть их в веб-интерфейсе или через API.

## 13. Полезные ссылки

- **Основная документация:** README.md
- **Production деплой:** docs/DEPLOYMENT.md
- **Docker Desktop:** https://www.docker.com/products/docker-desktop/
- **Node.js:** https://nodejs.org/ (только для локальной разработки frontend)

## Поддержка

При возникновении проблем:
1. Проверьте логи: `docker-compose logs`
2. Убедитесь, что Docker Desktop запущен
3. Проверьте доступность портов
4. Попробуйте полную пересборку: `docker-compose down -v && docker-compose up -d --build`
