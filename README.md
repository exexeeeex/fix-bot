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

Требования к инструментам, сервисам и настройкам окружения перечислены в [requirements.txt](requirements.txt). Нужны Docker Compose, Go версии из `backend/go.mod`, Node.js и Yarn.

1. Запустите PostgreSQL и Redis из корня проекта:

    ```powershell
    docker compose -f backend/compose.yaml up -d
    ```

2. Убедитесь, что в `backend/.env.dev` настроено подключение к локальным сервисам. Backend читает этот файл при `APP_ENV=development`; обязательны `DATABASE_URL` и доступ к Redis. `MAX_BOT_TOKEN` нужен для проверки авторизации через MAX.

3. Примените миграции из папки `backend`. Для этого нужна установленная утилита `migrate`:

    ```powershell
    migrate -path .\migrations -database "<DATABASE_URL из backend/.env.dev>" up
    ```

    Замените значение в угловых скобках строкой подключения к локальной PostgreSQL.

4. В первом терминале запустите backend:

    ```powershell
    cd backend
    go run ./cmd/api
    ```

    API будет доступно на `http://localhost:8080`.

5. Во втором терминале установите зависимости и запустите frontend:

    ```powershell
    cd frontend
    corepack enable
    yarn install --frozen-lockfile
    yarn dev
    ```

    Приложение будет доступно на `http://localhost:3000`. В режиме разработки frontend обращается к API на `http://localhost:8080`.

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
