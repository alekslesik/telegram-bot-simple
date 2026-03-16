# telegram-bot-simple

https://t.me/TYygGIBZQUTFmqoA_bot

Простой шаблон Telegram-бота на Go (long polling, уровень 1) с Docker, рассчитанный на деплой на VPS.

## Требования

- Go 1.25+
- Docker (для контейнеризации)
- Аккаунт Telegram и созданный бот через `@BotFather`

## Настройка окружения

1. Создай бота в Telegram через `@BotFather` и получи токен.
2. В файле `.env` должны быть переменные (файл не коммитится в git):

   ```env
   TOKEN=ваш_токен_бота
   USERNAME=ваш_username_бота_без_@
   ```

3. Локальный запуск:

   ```bash
   export $(grep -v '^#' .env | xargs)
   go run main.go
   ```

## Запуск в Docker

Собрать образ:

```bash
docker build -t telegram-bot-simple .
```

Запустить контейнер (используя `.env`):

```bash
docker run --rm \
  --env-file .env \
  --name telegram-bot-simple \
  telegram-bot-simple
```

## Деплой на VPS (общая схема)

1. Создать VPS у любого провайдера (Hetzner, DO, Contabo и т.д.) с установленным Docker.
2. Скопировать файлы проекта на сервер (например, через `scp` или `git clone`).
3. На сервере:

   ```bash
   cd telegram-bot-simple
   docker build -t telegram-bot-simple .
   docker run -d \
     --env-file .env \
     --name telegram-bot-simple \
     --restart unless-stopped \
     telegram-bot-simple
   ```

Бот использует long polling, поэтому не требуется публичный HTTPS и настройка webhook — достаточно, чтобы VPS имел доступ в интернет.
