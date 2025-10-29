# Webhook Dispatcher

[English Version](README.md)

[![Build Status](https://img.shields.io/badge/build-passing-brightgreen)](https://github.com/shanth1/hookrelay)
[![Go Report Card](https://goreportcard.com/badge/github.com/shanth1/hookrelay)](https://goreportcard.com/report/github.com/shanth1/hookrelay)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Гибкий сервис на Go для получения вебхуков из разных источников, таких как GitHub, Kanboard и кастомные приложения, и их отправки в виде форматированных уведомлений в различные каналы, например Telegram и Email.

## Обзор

`hookrelay` (`webhook-dispatcher`) — это легковесный, настраиваемый сервер, который работает как мост между вашими сервисами и каналами уведомлений. Он принимает события вебхуков, проверяет их подлинность, используя специфичные для каждого источника методы, форматирует их в человекочитаемые сообщения с помощью Go-шаблонов и рассылает по заранее настроенному списку получателей.

Это позволяет создать централизованную систему уведомлений для вашей команды, чтобы все были в курсе активностей с разных платформ в одном месте.

## Возможности

- **Обработка вебхуков из разных источников**: Встроенная поддержка вебхуков от **GitHub**, **Kanboard** и **пользовательских (Custom)** источников.
- **Безопасная и гибкая верификация**: Каждая точка приема вебхуков может иметь свой секретный ключ и метод проверки:
  - **GitHub**: HMAC-подпись (`X-Hub-Signature-256`).
  - **Kanboard**: Токен в URL-запросе.
  - **Custom**: Токен в заголовке аутентификации (`X-Auth-Token`).
- **Многоканальные уведомления**: Поддержка **Telegram** и **Email (SMTP)** «из коробки». Архитектура позволяет легко добавлять новые сервисы.
- **Настраиваемые шаблоны сообщений**: Использует Go `html/template` для форматирования уведомлений. Предоставляются шаблоны для популярных событий, которые легко можно изменить или расширить для каждого источника.
- **Продвинутая маршрутизация**: Мощная конфигурация в формате YAML позволяет определять множество эндпоинтов для вебхуков и направлять уведомления от них конкретным получателям или группам.
- **Конфигурация в формате YAML**: Все настройки, включая адрес сервера, эндпоинты, секретные ключи, отправителей и получателей, управляются через один простой для чтения YAML-файл.
- **Корректное завершение работы**: Сервер обрабатывает сигналы `SIGINT`/`SIGTERM` для чистого и безопасного выключения.

## Требования

- [Go](https://golang.org/dl/) (рекомендуется версия 1.18 или выше)
- [Make](https://www.gnu.org/software/make/) для использования удобных команд.
- Для тестирования с помощью скрипта `webhook.sh`:
  - [curl](https://curl.se/)
  - [jq](https://stedolan.github.io/jq/)
  - [openssl](https://www.openssl.org/)

## Начало работы

### 1. Клонируйте репозиторий

```bash
git clone https://github.com/your-username/webhook-dispatcher.git
cd webhook-dispatcher
```

### 2. Создайте файлы конфигурации

Проект использует `Makefile` для упрощения настройки. Следующая команда скопирует пример конфигурации в `config/local.yaml` (для разработки) и `config/production.yaml`.

```bash
make configs
```

### 3. Соберите приложение

`Makefile` предоставляет удобную команду для сборки. Исполняемый файл (`hookrelay`) будет помещен в директорию `build/`.

```bash
make build
```

## Конфигурация

Приложение настраивается с помощью YAML-файла (например, `config/production.yaml`). Конфигурация разделена на три основные секции: `webhooks`, `senders` и `recipients`.

**Отредактируйте `config/production.yaml`:**

```yaml
# Адрес и порт, который будет слушать сервер.
addr: ':8080'

# Здесь определяются все входящие эндпоинты для вебхуков.
webhooks:
  - name: 'github-repo-events'
    path: '/webhook/github'
    type: 'github' # Поддерживаемые типы: github, kanboard, custom
    secret: 'your-super-secret-github-string' # Секрет для HMAC-подписи
    recipients:
      - 'Dev Team (Telegram)' # Список имен получателей (определены ниже)
      - 'Admin (Email)'

  - name: 'kanboard-project-updates'
    path: '/webhook/kanboard'
    type: 'kanboard'
    secret: 'your-kanboard-secret-token' # Секрет, используемый как токен в URL
    recipients:
      - 'Dev Team (Telegram)'

  - name: 'custom-service-alerts'
    path: '/webhook/custom'
    type: 'custom'
    secret: 'secret-for-custom-messages' # Секрет для заголовка X-Auth-Token
    recipients:
      - 'Admin (Email)'
      - 'On-call Engineer'

# "Отправители" (senders) определяют, "как" и "откуда" отправлять уведомления.
senders:
  - name: 'project-telegram-bot' # Уникальное имя для этого отправителя
    type: 'telegram'
    settings:
      token: 'YOUR_TELEGRAM_BOT_TOKEN'

  - name: 'smtp-notifications'
    type: 'email'
    settings:
      host: 'smtp.example.com'
      port: 587
      username: 'user@example.com'
      password: 'your-smtp-password'
      from: 'Notifier <no-reply@example.com>'

# "Получатели" (recipients) связывают адресата (target) с отправителем (sender).
recipients:
  - name: 'Dev Team (Telegram)'
    sender: 'project-telegram-bot' # Должно совпадать с именем отправителя выше
    target: '-100123456789' # ID чата/канала в Telegram

  - name: 'Admin (Email)'
    sender: 'smtp-notifications'
    target: 'admin@example.com' # Адрес электронной почты

  - name: 'On-call Engineer'
    sender: 'project-telegram-bot'
    target: '987654321' # ID другого чата в Telegram
```

## Запуск приложения

Передайте путь к вашему файлу конфигурации через флаг `--config`.

```bash
./build/hookrelay --config config/production.yaml
```

Вы также можете использовать команду `make`, которая запускает приложение с файлом `config/local.yaml`:

```bash
make run
```

Сервер запустится и выведет в лог сообщение: `{"level":"info","time":"...","message":"staring server on :8080"}`.

## Настройка вебхуков

### GitHub

1.  Перейдите в настройки вашего репозитория на GitHub: **Settings > Webhooks > Add webhook**.
2.  **Payload URL**: `http://your-server.com:8080/webhook/github` (или `path`, который вы указали в конфиге).
3.  **Content type**: `application/x-www-form-urlencoded` или `application/json`.
4.  **Secret**: Введите ту же самую строку, что и в поле `secret` для соответствующего блока `webhooks` в вашем конфиге.

### Kanboard

1.  В Kanboard перейдите в **Настройки проекта > Вебхуки**.
2.  **Webhook URL**: `http://your-server.com:8080/webhook/kanboard?token=your-kanboard-secret-token`. `token` в URL должен совпадать с полем `secret` из вашего конфига.

### Пользовательский источник (Custom)

1.  Настройте ваше приложение на отправку `POST`-запроса на `http://your-server.com:8080/webhook/custom`.
2.  Тело запроса должно быть в формате `text/plain` и содержать сообщение, которое вы хотите отправить.
3.  Запрос должен включать заголовок `X-Auth-Token` со значением `secret` из вашего конфига.

## Тестирование

Проект включает в себя скрипт для тестирования (`webhook.sh`) и стандартные тесты Go.

### Запуск юнит-тестов Go

```bash
make test
```

### Симуляция вебхуков

Скрипт `webhook.sh` может симулировать отправку событий вебхуков на локально запущенный сервер. Это очень полезно для проверки конфигурации и шаблонов сообщений.

**Важно**: Убедитесь, что секретные ключи внутри `webhook.sh` совпадают со значениями `secret` в вашем YAML-конфиге.

- **Отправить событие `push` от GitHub:**
  (Сначала убедитесь, что сервер `hookrelay` запущен в другом терминале)

  ```bash
  ./webhook.sh github push
  ```

- **Отправить событие `task.create` от Kanboard:**

  ```bash
  ./webhook.sh kanboard task.create
  ```

- **Отправить кастомное оповещение:**

  ```bash
  ./webhook.sh custom "ТРЕВОГА: Сервис недоступен!"
  ```

- **Использовать собственный JSON из файла для события `issues` от GitHub:**
  ```bash
  ./webhook.sh github issues ./payloads/my_issue.json
  ```

## Команды Makefile

| Команда        | Описание                                                                |
| :------------- | :---------------------------------------------------------------------- |
| `make help`    | Показать список доступных команд.                                       |
| `make configs` | Создать файлы конфигурации `local.yaml` и `production.yaml` из примера. |
| `make run`     | Запустить приложение с файлом `config/local.yaml`.                      |
| `make build`   | Собрать бинарный файл приложения в директорию `build/`.                 |
| `make test`    | Запустить все юнит-тесты Go.                                            |
| `make clean`   | Удалить директорию `build/` и все скомпилированные файлы.               |

## Структура проекта

```
.
├── build/                 # Скомпилированные бинарные файлы (создаются через make)
├── cmd/
│   └── main.go            # Главная точка входа приложения
├── config/
│   └── example.yaml       # Пример файла конфигурации
├── internal/              # Внутренние пакеты Go
│   ├── app/               # Настройка приложения и жизненный цикл сервера
│   ├── config/            # Структуры для конфигурации
│   ├── middleware/        # HTTP middleware (логирование, верификация)
│   ├── notifier/          # Логика отправки уведомлений (Telegram, Email)
│   ├── processor/         # Обработка данных вебхуков (GitHub, Kanboard)
│   ├── server/            # HTTP-сервер и маршрутизация
│   ├── templates/         # Логика загрузки Go-шаблонов
│   └── verifier/          # Проверка подписей/токенов вебхуков
├── templates/             # Go-шаблоны для форматирования уведомлений
│   ├── github/*.tmpl      # Шаблоны для событий GitHub
│   ├── kanboard/*.tmpl    # Шаблоны для событий Kanboard
│   └── embed.go           # Встраивает шаблоны в бинарный файл
├── Makefile               # Команды для сборки, запуска и тестирования
├── go.mod                 # Определения Go-модуля
└── webhook.sh             # Скрипт для отправки тестовых вебхуков
```

## Участие в разработке

Pull request'ы приветствуются. Для крупных изменений, пожалуйста, сначала создайте issue, чтобы обсудить, что вы хотели бы изменить. Пожалуйста, не забудьте обновить тесты при необходимости.

## Лицензия

Проект распространяется под [лицензией MIT](LICENSE).
