# Фикс

«Фикс» - веб-приложение для работы с электронными пропусками и посещаемостью студентов. Студенту доступен QR-код пропуска; преподаватели и администраторы могут работать с группами, занятиями, заявками и статистикой посещений. Для входа предусмотрена интеграция с MAX.

## Хакатон

Проект разработан для хакатона от Сферума «Идея Фикс».

## Участники

- [ka3tri](https://github.com/ka3tri)
- [lissoppw](https://github.com/lissoppw)
- [zxcalesha](https://github.com/zxcalesha)
- [exexeeeex](https://github.com/exexeeeex)
- [aexqodev](https://github.com/aexqodev)
- [afoninigor716-bit](https://github.com/afoninigor716-bit)

## Технологии

- Frontend: React, TypeScript, React Router, Vite.
- Backend: Go, PostgreSQL, Redis.
- Инфраструктура: Docker Compose, Nginx.

## Docker-образы

- [Frontend — feels/fix_bot-frontend](https://hub.docker.com/repository/docker/feels/fix_bot-frontend)
- [Backend — feels/fix_bot-backend](https://hub.docker.com/repository/docker/feels/fix_bot-backend)

Production-команда Docker Compose из этого README собирает frontend и backend из исходного кода. Чтобы запускать опубликованные образы, укажите их в Compose вместо секций `build`; теги образов смотрите на Docker Hub.

## Ограничения сканирования QR-кодов

В связи с техническими ограничениями мессенджера MAX на платформах Android и Windows получение доступа к камере из встроенного браузера может работать нестабильно или быть недоступно.

Поэтому для стабильной работы сканера QR-кодов рекомендуется использовать:

- **iOS (iPhone)** — сканирование QR-кодов непосредственно в приложении.
- **Обычный браузер** — откройте [fiksbot.ru](https://fiksbot.ru) напрямую, чтобы воспользоваться приложением в демо-режиме.

На Android и Windows корректная работа сканера через встроенный браузер MAX не гарантируется. Это ограничение связано с особенностями платформы MAX, а не с работой самого приложения.

## Локальный запуск для разработки

Требования к инструментам и настройкам окружения перечислены в [requirements.txt](requirements.txt). Для запуска нужен Docker Compose и заполненный `backend/.env.dev`; укажите в нём необходимые секреты, например `MAX_BOT_TOKEN`. Compose сам задаст адреса PostgreSQL и Redis внутри сети контейнеров и применит миграции.

Запустите весь стек из корня проекта одной командой:

```powershell
docker compose -f compose.dev.yaml up --build
```

Приложение будет доступно по адресу `http://localhost:8080`; API проксируется через Nginx. Чтобы изменить порт Nginx, задайте `DEV_HTTP_PORT`. Исходники frontend и backend подключены в контейнеры, frontend поддерживает hot reload. PostgreSQL и Redis также доступны с хоста на портах `5432` и `6379`.

Остановить стек можно командой:

```powershell
docker compose -f compose.dev.yaml down
```

## Запуск production-стека

Для развёртывания production-версии проекта не требуется клонировать весь репозиторий или иметь исходный код frontend и backend на сервере. Приложение запускается из опубликованных Docker-образов:

- [Frontend — feels/fix_bot-frontend](https://hub.docker.com/repository/docker/feels/fix_bot-frontend)
- [Backend — feels/fix_bot-backend](https://hub.docker.com/repository/docker/feels/fix_bot-backend)

Docker Compose загружает необходимые образы из Docker Hub и запускает приложение.

### Структура файлов на сервере

Для развёртывания достаточно подготовить следующую структуру:

```text
fix_bot/
├── compose.yaml
├── .environments/
│   ├── psq.prod.env
│   └── prod.env
├── backend/
│   └── migrations/
├── migrations/
└── nginx/
    ├── nginx.conf
    └── certs/
        ├── fullchain.pem
        ├── chain.pem
        ├── cert.pem
        └── privkey.pem
```

В этой структуре:

- `compose.yaml` — конфигурация Docker Compose, использующая опубликованные образы вместо сборки из исходного кода.
- `.environments/prod.env` — переменные окружения PostgreSQL, Redis и backend, включая необходимые секреты.
- `backend/migrations/` и `migrations/` — файлы миграций базы данных, используемые при развёртывании.
- `nginx/` — конфигурация Nginx и TLS-сертификаты.

Наличие исходников frontend и backend на сервере не требуется.

### Подготовка окружения

Перед запуском необходимо:

1. Установить Docker Engine и Docker Compose.
2. Подготовить файл `.environments/prod.env`, указав параметры PostgreSQL, Redis, адреса подключения backend к сервисам внутри Docker Compose и необходимые секреты.
3. Разместить файлы миграций в соответствующих директориях.
4. Подготовить конфигурацию Nginx и TLS-сертификаты:
   - `nginx/certs/fullchain.pem`
   - `nginx/certs/privkey.pem`

Не публикуйте секреты и TLS-ключи в репозитории.

### Запуск

Из директории `fix_bot` выполните:

```bash
docker compose pull
docker compose up -d
```

Первая команда загружает опубликованные Docker-образы, вторая запускает production-стек.

Compose применит миграции перед запуском backend.

По умолчанию Nginx принимает HTTP на порту 80 и HTTPS на 443. Порты можно изменить с помощью переменных `HTTP_PORT` и `HTTPS_PORT`.

### Управление

Посмотреть состояние контейнеров:

```bash
docker compose ps
```

Посмотреть журналы всех сервисов:

```bash
docker compose logs -f
```

Посмотреть журналы конкретного сервиса:

```bash
docker compose logs -f backend
```

Остановить production-стек:

```bash
docker compose down
```

Для обновления приложения после публикации новых версий Docker-образов выполните:

```bash
docker compose pull
docker compose up -d
```

При этом данные PostgreSQL и Redis сохраняются в Docker volumes, поэтому обычное обновление контейнеров не удаляет данные приложения.
