# HookRelay (Диспетчер Вебхуков)

[English Version](README.md)

Гибкий сервис на Go для приема вебхуков из различных источников (таких как GitHub, Kanboard или собственные приложения) и отправки отформатированных уведомлений в различные каналы (Telegram, Email).

## Обзор

`hookrelay` — это легковесный, настраиваемый сервер, который служит мостом между вашими сервисами и каналами уведомлений. Он слушает события вебхуков, проверяет их подлинность, форматирует в читаемые сообщения с использованием Go-шаблонов и рассылает их настроенному списку получателей.

## Возможности

- **Обработка вебхуков из разных источников**: Встроенная поддержка **GitHub**, **Kanboard** и **Custom** (произвольных) источников.
- **Безопасная проверка**:
  - **GitHub**: Подпись HMAC (`X-Hub-Signature-256`).
  - **Kanboard**: Токен в параметрах URL.
  - **Custom**: Токен в заголовке авторизации (`X-Auth-Token`).
- **Уведомления в разные каналы**: Поддержка **Telegram** и **Email (SMTP)** через конфигурацию `notifiers`.
- **Шаблонизация сообщений**: Использование встроенных Go-шаблонов (`html/template`).
  - Поддержка фоллбэка (стандартного шаблона) для неизвестных событий.
  - Специфичные шаблоны для сложных событий (например, GitHub Push, создание задачи в Kanboard).
- **API и диагностика**: Эндпоинты для проверки здоровья (health check) и получения информации о конфигурации.
- **Graceful Shutdown**: Корректное завершение работы при получении сигналов остановки.

## Требования

- [Go](https://golang.org/dl/) (версия 1.23 или выше)
- [Make](https://www.gnu.org/software/make/) (для сборки)

## Быстрый старт

### 1. Клонирование репозитория

```bash
git clone https://github.com/shanth1/hookrelay.git
cd hookrelay
```

### 2. Создание конфигурации

Команда ниже создаст файлы конфигурации на основе примера:

```bash
make configs
# Создает config/local.yaml и config/production.yaml
```

### 3. Сборка

```bash
make build
# Бинарный файл будет создан в ./build/hookrelay
```

## Конфигурация

Настройка выполняется через YAML файл.

**Пример `config/production.yaml`:**

```yaml
# Адрес и порт сервера
addr: ':8080'
env: 'production'

# Если true, события без специального шаблона будут игнорироваться,
# вместо использования шаблона по умолчанию.
disable_unknown_templates: false

logger:
  level: 'info'
  app: 'hookrelay'
  service: 'webhook-service'

# 1. Определение входящих вебхуков (Источники)
webhooks:
  - name: 'github-repo'
    path: '/webhook/github'
    type: 'github'
    secret: 'YOUR_GITHUB_WEBHOOK_SECRET' # Секрет для HMAC подписи
    recipients:
      - 'Dev Team (Telegram)'
      - 'Admin (Email)'

  - name: 'kanboard-project'
    path: '/webhook/kanboard'
    type: 'kanboard'
    secret: 'YOUR_KANBOARD_TOKEN' # Токен для URL
    base_url: 'https://kanboard.example.com' # Обязательно для генерации ссылок на задачи
    recipients:
      - 'Dev Team (Telegram)'

  - name: 'custom-alert'
    path: '/webhook/custom'
    type: 'custom'
    secret: 'YOUR_CUSTOM_AUTH_TOKEN'
    recipients:
      - 'Dev Team (Telegram)'

# 2. Определение Уведомлений (Каналы отправки)
notifiers:
  - name: 'telegram-bot'
    type: 'telegram'
    settings:
      token: 'YOUR_TELEGRAM_BOT_TOKEN'

  - name: 'smtp-mail'
    type: 'email'
    settings:
      host: 'smtp.yandex.ru'
      port: 465
      username: 'notifier@yandex.ru'
      password: 'app-password'
      from: 'Notifier <notifier@yandex.ru>'

# 3. Определение Получателей (Маршрутизация)
recipients:
  - name: 'Dev Team (Telegram)'
    notifier: 'telegram-bot' # Должно совпадать с именем в notifiers
    target: '-100123456789' # ID чата Telegram

  - name: 'Admin (Email)'
    notifier: 'smtp-mail'
    target: 'admin@example.com' # Email адрес
```

## Запуск

**Через Make (для разработки):**

```bash
make run
```

**Бинарный файл (продакшн):**

```bash
./build/hookrelay --config config/production.yaml
```

## Настройка источников

### GitHub

1. В настройках репозитория: Settings -> Webhooks.
2. Payload URL: `http://ваш-сервер/webhook/github`.
3. Content type: `application/json` или `application/x-www-form-urlencoded`.
4. Secret: Должен совпадать с `secret` в конфиге `webhooks`.

### Kanboard

1. В настройках проекта: Webhooks.
2. Webhook URL: `http://ваш-сервер/webhook/kanboard?token=YOUR_KANBOARD_TOKEN`.
3. Параметр `token` в URL должен совпадать с `secret` в конфиге.

### Custom (Произвольный)

1. Отправьте `POST` запрос на `http://ваш-сервер/webhook/custom`.
2. Заголовок `X-Auth-Token`: Должен совпадать с `secret` в конфиге.
3. Тело запроса: Текст или JSON.

## Тестирование

Используйте скрипт `webhook.sh` для эмуляции событий:

```bash
# GitHub Push Event
./webhook.sh github push

# Kanboard Task Creation
./webhook.sh kanboard task.create

# Custom Message
./webhook.sh custom "Внимание: Сервис недоступен!"
```

## Структура проекта

```
.
├── cmd/                # Точка входа (main.go)
├── config/             # Файлы конфигурации и Go-структуры конфига
├── internal/
│   ├── adapters/       # Адаптеры: Входящие (GitHub/Kanboard) и Исходящие (TG/Email)
│   ├── app/            # Инициализация приложения
│   ├── core/           # Доменная логика и интерфейсы (порты)
│   ├── service/        # Бизнес-логика (Маршрутизация уведомлений)
│   └── transport/      # HTTP сервер, роутинг и middleware
└── templates/          # Встроенные шаблоны уведомлений
```

## Лицензия

MIT License. См. файл [LICENSE](LICENSE).
