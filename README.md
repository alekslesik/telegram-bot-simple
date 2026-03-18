# telegram-bot-simple

Демо-бот в Telegram: `https://t.me/TYygGIBZQUTFmqoA_bot`

Это демонстрационный Telegram‑бот, который можно показывать потенциальным заказчикам как «живой пример продукта».

## Что умеет бот (для заказчика)

- **Базовое общение**: бот принимает сообщения и отвечает (пример обработки обращений клиентов).
- **Информативные команды**:
  - **`/start`** — короткое приветствие и что делает бот
  - **`/help`** — список возможностей
  - **`/about`** — зачем такой бот бизнесу
  - **`/usecases`** — примеры задач (услуги, обучение, малый бизнес)
  - **`/features`** — какие функции можно добавить в такого бота (заявки, меню, запись, опросы и т.д.)
  - **`/ping`** — проверка, что бот онлайн
  - **`/echo <текст>`** — пример простой команды (повторяет текст)

## Для каких задач лучше выбрать Telegram-бота

- **Поддержка клиентов 24/7**: ответы на частые вопросы, выдача инструкций.
- **Сбор заявок и контактов**: «оставьте телефон / комментарий» прямо в чате.
- **Автоматизация типовых сценариев**: прайсы, расписание, анкеты, опросы, напоминания.
- **Лёгкий вход в продукт**: Telegram не требует установки приложения — клиент уже там.

## Быстрый старт (для разработки)

### Требования

- Go 1.25+
- Docker + Docker Compose (опционально)
- Аккаунт Telegram и созданный бот через `@BotFather`

### Настройка `.env`

1. Скопируй пример:

```bash
cp .env.example .env
```

2. Заполни `TOKEN` и `USERNAME`.

Переменные:

- **`TOKEN`**: токен бота от `@BotFather`
- **`USERNAME`**: username бота (без `@`)
- **`LOG_LEVEL`**: `debug` или `info` (по умолчанию `info`)
- **`LOG_FORMAT`**: `json` или `text` (по умолчанию `text`)

### Запуск локально

```bash
make run
```

### Тесты

```bash
make test
```

## Запуск в Docker

```bash
make docker-run
```

## Запуск через Docker Compose

```bash
make docker-compose-up
```

Остановить:

```bash
make docker-compose-down
```

## CI/CD: GHCR + автодеплой на VPS по SSH

В репозитории настроены workflow’ы GitHub Actions:

- **CI**: запускает проверки как в `Makefile` (форматирование, линтеры, тесты, vuln, docker build).
- **Release**: собирает Docker-образ и пушит в GitHub Container Registry (GHCR) `ghcr.io/<owner>/<repo>`.
- **Deploy**: после успешного релиза автоматически деплоит на VPS по SSH.

### Подготовка VPS (один раз)

1. Установи Docker и Docker Compose на сервер.
2. Создай директорию для бота, например:

```bash
sudo mkdir -p /opt/bots/telegram-bot-simple
sudo chown -R $USER:$USER /opt/bots/telegram-bot-simple
cd /opt/bots/telegram-bot-simple
```

3. Скопируй на сервер файлы `docker-compose.prod.yaml` и `.env` (секреты храним на сервере, не в git):

```bash
scp docker-compose.prod.yaml user@server:/opt/bots/telegram-bot-simple/
scp .env user@server:/opt/bots/telegram-bot-simple/
```

### Секреты GitHub для деплоя

В Settings → Secrets and variables → Actions добавь:

- **`VPS_HOST`**: IP/домен сервера
- **`VPS_USER`**: пользователь (например `root` или `deploy`)
- **`VPS_SSH_KEY`**: приватный SSH ключ (PEM), которым GitHub будет подключаться
- **`VPS_APP_PATH`**: путь на сервере, например `/opt/bots/telegram-bot-simple`
- **`GHCR_READ_USER`**: логин GitHub (владелец пакета)
- **`GHCR_READ_TOKEN`**: токен с правами `read:packages` (для `docker login ghcr.io` на сервере)

### Как работает деплой

- Пуш в ветку **`dev`** публикует образ `:dev` и деплоит его на VPS.
- Пуш в **`main`/`master`** публикует образ `:main` и деплоит его на VPS.

На сервере используется `docker-compose.prod.yaml`, который запускает образ из GHCR:

```yaml
services:
  bot:
    image: ghcr.io/alekslesik/telegram-bot-simple:${IMAGE_TAG:-latest}
```

## Деплой на VPS (общая схема)

1. Создать VPS у любого провайдера с установленным Docker.
2. Скопировать проект на сервер (`git clone` или `scp`).
3. На сервере:

```bash
cd telegram-bot-simple
cp .env.example .env
nano .env   # заполнить TOKEN/USERNAME и при желании LOG_LEVEL/LOG_FORMAT
docker compose up -d --build
```

Бот использует long polling, поэтому обычно не требуется публичный HTTPS и настройка webhook — достаточно, чтобы сервер имел доступ в интернет.
