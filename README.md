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

Для production нужны Docker Compose, настроенный файл `.environments/prod.env` и TLS-сертификаты Nginx: `nginx/certs/fullchain.pem` и `nginx/certs/privkey.pem`. Укажите в env-файле параметры PostgreSQL и Redis, адреса подключения backend к сервисам внутри Compose, а также необходимые секреты. Не публикуйте секреты в репозитории.

Из корня проекта выполните:

```powershell
docker compose up --build -d
```

Compose применит миграции перед запуском backend. По умолчанию Nginx принимает HTTP на порту 80 и HTTPS на 443; порты можно изменить переменными `HTTP_PORT` и `HTTPS_PORT`. Посмотреть состояние сервисов и журналы можно командами:

```powershell
docker compose ps
docker compose logs -f
```
