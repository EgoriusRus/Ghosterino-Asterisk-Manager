# Production Deployment Guide

Руководство по развертыванию Asterisk Manager на продакшн сервере.

## Архитектура

```
Производственный сервер
├── Nginx (reverse proxy + SSL)
│   └── https://api.yourdomain.com → localhost:8080
├── Docker Compose
│   ├── PostgreSQL (изолированный, без внешнего порта)
│   └── API Backend (bind на 127.0.0.1:8080)
├── Asterisk PBX (на хосте)
│   ├── /srv/tftp/ ← конфиги телефонов
│   ├── /etc/asterisk/users/ ← SIP peers
│   └── /etc/asterisk/extensions/ ← Dialplan
└── Systemd (автозапуск сервисов)
```

## Требования

- Ubuntu 20.04+ / CentOS 8+ / Debian 11+
- Docker 20.10+
- Docker Compose 2.0+
- Nginx
- Asterisk PBX (установлен на хосте)
- Минимум 2GB RAM, 20GB disk

## Быстрый старт

### 1. Подготовка сервера

```bash
# Обновление системы
sudo apt update && sudo apt upgrade -y

# Установка Docker
curl -fsSL https://get.docker.com | sh
sudo usermod -aG docker $USER
newgrp docker

# Установка Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Установка Nginx
sudo apt install -y nginx certbot python3-certbot-nginx

# Создание рабочей директории
sudo mkdir -p /opt/asterisk-manager
sudo chown -R $USER:$USER /opt/asterisk-manager
```

### 2. Распаковка проекта

```bash
# Загрузите asterisk-manager.zip на сервер
unzip asterisk-manager.zip
sudo mv asterisk-manager /opt/
cd /opt/asterisk-manager
```

### 3. Настройка переменных окружения

```bash
# Копирование примера
cp .env.prod.example .env.prod

# Редактирование (установите сильные пароли!)
nano .env.prod
```

Обязательные изменения в `.env.prod`:
```bash
DB_USER=asterisk_prod
DB_PASSWORD=YOUR_STRONG_PASSWORD_HERE  # ИЗМЕНИТЕ!
DB_NAME=asterisk_manager_prod
APP_PORT=8080

# Пути к директориям Asterisk
ASTERISK_TFTP_DIR=/srv/tftp
ASTERISK_USERS_DIR=/etc/asterisk/users
ASTERISK_EXT_DIR=/etc/asterisk/extensions
```

### 4. Создание директорий Asterisk

```bash
# Создание директорий
sudo mkdir -p /srv/tftp
sudo mkdir -p /etc/asterisk/users
sudo mkdir -p /etc/asterisk/extensions
sudo mkdir -p /opt/asterisk-manager/results
sudo mkdir -p /opt/asterisk-manager/backups

# Установка прав
sudo chown -R $USER:$USER /srv/tftp
sudo chown -R $USER:$USER /etc/asterisk/users
sudo chown -R $USER:$USER /etc/asterisk/extensions
```

### 5. Запуск Docker Compose

```bash
cd /opt/asterisk-manager

# Запуск в продакшн режиме
docker-compose -f docker-compose.prod.yml up -d

# Проверка статуса
docker-compose -f docker-compose.prod.yml ps

# Просмотр логов
docker-compose -f docker-compose.prod.yml logs -f
```

### 6. Заполнение базы данных

```bash
# Вход в контейнер
docker-compose -f docker-compose.prod.yml exec backend sh

# Запуск seed (внутри контейнера)
cd /app
go run cmd/seed/main.go

# Выход
exit
```

### 7. Настройка Nginx

```bash
# Копирование конфигурации
sudo cp deploy/nginx.conf /etc/nginx/sites-available/asterisk-api

# Редактирование (замените yourdomain.com на ваш домен)
sudo nano /etc/nginx/sites-available/asterisk-api

# Включение сайта
sudo ln -s /etc/nginx/sites-available/asterisk-api /etc/nginx/sites-enabled/

# Проверка конфигурации
sudo nginx -t

# Перезапуск Nginx
sudo systemctl restart nginx
```

### 8. Получение SSL сертификата

```bash
# Let's Encrypt
sudo certbot --nginx -d api.yourdomain.com

# Автообновление сертификата
sudo systemctl enable certbot.timer
```

### 9. Настройка автозапуска (systemd)

```bash
# Копирование service файла
sudo cp deploy/systemd/asterisk-manager.service /etc/systemd/system/

# Редактирование путей если нужно
sudo nano /etc/systemd/system/asterisk-manager.service

# Перезагрузка systemd
sudo systemctl daemon-reload

# Включение автозапуска
sudo systemctl enable asterisk-manager.service

# Запуск сервиса
sudo systemctl start asterisk-manager.service

# Проверка статуса
sudo systemctl status asterisk-manager.service
```

## Управление сервисом

### Systemd команды

```bash
# Статус
sudo systemctl status asterisk-manager

# Запуск
sudo systemctl start asterisk-manager

# Остановка
sudo systemctl stop asterisk-manager

# Перезапуск
sudo systemctl restart asterisk-manager

# Логи
sudo journalctl -u asterisk-manager -f
```

### Docker Compose команды

```bash
cd /opt/asterisk-manager

# Статус контейнеров
docker-compose -f docker-compose.prod.yml ps

# Логи
docker-compose -f docker-compose.prod.yml logs -f backend
docker-compose -f docker-compose.prod.yml logs -f postgres

# Перезапуск отдельного сервиса
docker-compose -f docker-compose.prod.yml restart backend

# Остановка всех контейнеров
docker-compose -f docker-compose.prod.yml down

# Пересборка и запуск
docker-compose -f docker-compose.prod.yml up -d --build
```

## Генерация конфигов Asterisk

### Ручная генерация

```bash
# Через make
make generator

# Или напрямую через Docker
docker-compose -f docker-compose.prod.yml exec backend go run cmd/generator/main.go

# Перезагрузка Asterisk
sudo asterisk -rx "reload"
```

### Автоматическая генерация (cron)

Создайте cron задачу для регулярной генерации конфигов:

```bash
crontab -e
```

Добавьте:
```cron
# Генерация конфигов Asterisk каждые 5 минут
*/5 * * * * cd /opt/asterisk-manager && docker-compose -f docker-compose.prod.yml exec -T backend go run cmd/generator/main.go && asterisk -rx "reload"

# Бэкап базы каждую ночь в 2:00
0 2 * * * cd /opt/asterisk-manager && make backup
```

## Бэкап и восстановление

### Создание бэкапа

```bash
# Ручной бэкап
make backup

# Или напрямую
docker-compose -f docker-compose.prod.yml exec -T postgres pg_dump -U asterisk_prod asterisk_manager_prod | gzip > backup_$(date +%Y%m%d).sql.gz
```

### Восстановление из бэкапа

```bash
# Распаковка и восстановление
gunzip < backup_20260104.sql.gz | docker-compose -f docker-compose.prod.yml exec -T postgres psql -U asterisk_prod asterisk_manager_prod
```

## Мониторинг

### Проверка здоровья API

```bash
# Health check
curl http://localhost:8080/

# Список профилей
curl http://localhost:8080/api/profiles

# Через HTTPS
curl https://api.yourdomain.com/api/profiles
```

### Логи

```bash
# Backend логи
docker-compose -f docker-compose.prod.yml logs -f backend

# PostgreSQL логи
docker-compose -f docker-compose.prod.yml logs -f postgres

# Nginx access log
sudo tail -f /var/log/nginx/asterisk-api-access.log

# Nginx error log
sudo tail -f /var/log/nginx/asterisk-api-error.log
```

### Метрики

```bash
# Использование ресурсов контейнерами
docker stats

# Размер базы данных
docker-compose -f docker-compose.prod.yml exec postgres psql -U asterisk_prod -d asterisk_manager_prod -c "SELECT pg_size_pretty(pg_database_size('asterisk_manager_prod'));"

# Количество записей
docker-compose -f docker-compose.prod.yml exec postgres psql -U asterisk_prod -d asterisk_manager_prod -c "SELECT 'profiles', COUNT(*) FROM sipadmin.profiles UNION SELECT 'devices', COUNT(*) FROM sipadmin.devices UNION SELECT 'locations', COUNT(*) FROM sipadmin.locations;"
```

## Обновление приложения

```bash
cd /opt/asterisk-manager

# Остановка сервисов
docker-compose -f docker-compose.prod.yml down

# Распаковка новой версии
unzip asterisk-manager-new.zip
cp -r asterisk-manager-new/* .
rm -rf asterisk-manager-new

# Пересборка и запуск
docker-compose -f docker-compose.prod.yml up -d --build

# Или через systemd
sudo systemctl restart asterisk-manager
```

## Безопасность

### Firewall настройки

```bash
# UFW (Ubuntu)
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw allow 22/tcp
sudo ufw enable

# Asterisk порты (если нужен внешний доступ)
sudo ufw allow 5060/udp  # SIP
sudo ufw allow 10000:10100/udp  # RTP
```

### Рекомендации

- ✅ Используйте сильные пароли в `.env.prod`
- ✅ Регулярно обновляйте систему и Docker образы
- ✅ Настройте автоматические бэкапы
- ✅ Используйте HTTPS (Let's Encrypt)
- ✅ Ограничьте доступ к API через Nginx
- ✅ Включите rate limiting в Nginx
- ✅ Мониторьте логи на подозрительную активность

## Troubleshooting

### Контейнеры не запускаются

```bash
# Проверка логов
docker-compose -f docker-compose.prod.yml logs

# Проверка конфигурации
docker-compose -f docker-compose.prod.yml config

# Пересоздание контейнеров
docker-compose -f docker-compose.prod.yml down
docker-compose -f docker-compose.prod.yml up -d --force-recreate
```

### API не отвечает

```bash
# Проверка что backend запущен
docker-compose -f docker-compose.prod.yml ps backend

# Проверка портов
sudo netstat -tlnp | grep 8080

# Проверка внутри контейнера
docker-compose -f docker-compose.prod.yml exec backend wget -O- http://localhost:8080/
```

### База данных недоступна

```bash
# Проверка статуса PostgreSQL
docker-compose -f docker-compose.prod.yml ps postgres

# Подключение к БД
docker-compose -f docker-compose.prod.yml exec postgres psql -U asterisk_prod -d asterisk_manager_prod

# Проверка логов
docker-compose -f docker-compose.prod.yml logs postgres
```

### Asterisk не видит конфиги

```bash
# Проверка что файлы создались
ls -la /srv/tftp/
ls -la /etc/asterisk/users/
ls -la /etc/asterisk/extensions/

# Проверка прав доступа
ls -ld /srv/tftp/

# Проверка что Asterisk читает правильные директории
asterisk -rx "core show config"
```

## Поддержка

При возникновении проблем:
1. Проверьте логи контейнеров
2. Проверьте конфигурацию Nginx
3. Обратитесь к документации в этом файле
