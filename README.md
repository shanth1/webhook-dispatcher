# HookRelay (Webhook Dispatcher)

[Russian Version](README.ru.md)

[![Go Report Card](https://goreportcard.com/badge/github.com/shanth1/hookrelay)](https://goreportcard.com/report/github.com/shanth1/hookrelay)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A flexible Go-based service to receive webhooks from multiple sources like GitHub, Kanboard, and custom applications, and dispatch them as formatted notifications to various channels like Telegram and Email.

## Overview

`hookrelay` is a lightweight, configurable server that acts as a bridge between your services and your notification channels. It listens for webhook events, verifies their authenticity, formats them into human-readable messages using Go templates, and broadcasts them to a pre-configured list of recipients.

## Features

- **Multi-Source Webhook Handling**: Natively supports webhooks from **GitHub**, **Kanboard**, and **Custom** sources.
- **Secure Verification**:
  - **GitHub**: HMAC signature (`X-Hub-Signature-256`).
  - **Kanboard**: URL query token.
  - **Custom**: Authentication header token (`X-Auth-Token`).
- **Multi-Channel Notifications**: Support for **Telegram** and **Email (SMTP)** via the `notifiers` configuration.
- **Message Templating**: Uses embedded Go `html/template` files to format notifications.
  - Supports custom fallback for unknown events.
  - Specific templates for complex events (e.g., GitHub Push, Kanboard Task Create).
- **Service Discovery & Health**: Exposes endpoints for health checks and configuration discovery.
- **Graceful Shutdown**: Handles `SIGINT`/`SIGTERM` for clean shutdown.
- **YAML Configuration**: Single file configuration for endpoints, credentials, and routing.

## Prerequisites

- [Go](https://golang.org/dl/) (version 1.23 or higher)
- [Make](https://www.gnu.org/software/make/) (optional, for build commands)

## Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/shanth1/hookrelay.git
cd hookrelay
```

### 2. Create Configuration

Copy the example configuration to create your local config:

```bash
make configs
# Creates config/local.yaml and config/production.yaml
```

### 3. Build

```bash
make build
# Binary will be created at ./build/hookrelay
```

## Configuration

The application is configured using a YAML file.

**Example `config/production.yaml`:**

```yaml
# Address to listen on
addr: ':8080'
env: 'production'

# If true, events without a specific template file will be ignored
# instead of using the default template.
disable_unknown_templates: false

logger:
  level: 'info'
  app: 'hookrelay'
  service: 'webhook-service'

# 1. Define incoming Webhooks (Sources)
webhooks:
  - name: 'github-main'
    path: '/webhook/github'
    type: 'github'
    secret: 'YOUR_GITHUB_WEBHOOK_SECRET'
    recipients:
      - 'Dev Team (Telegram)'
      - 'Tech Lead (Email)'

  - name: 'kanboard-tasks'
    path: '/webhook/kanboard'
    type: 'kanboard'
    secret: 'YOUR_KANBOARD_TOKEN'
    base_url: 'https://kanboard.your-domain.com' # Required for generating links
    recipients:
      - 'Dev Team (Telegram)'

  - name: 'monitoring'
    path: '/webhook/custom'
    type: 'custom'
    secret: 'YOUR_CUSTOM_AUTH_TOKEN'
    recipients:
      - 'Dev Team (Telegram)'

# 2. Define Notifiers (Channels)
notifiers:
  - name: 'telegram-bot'
    type: 'telegram'
    settings:
      token: '123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11'

  - name: 'smtp-server'
    type: 'email'
    settings:
      host: 'smtp.gmail.com'
      port: 587
      username: 'notifier@example.com'
      password: 'app-specific-password'
      from: 'HookRelay <notifier@example.com>'

# 3. Define Recipients (Routing)
recipients:
  - name: 'Dev Team (Telegram)'
    notifier: 'telegram-bot' # Must match a notifier name above
    target: '-100123456789' # Telegram Chat ID

  - name: 'Tech Lead (Email)'
    notifier: 'smtp-server'
    target: 'lead@example.com' # Email Address
```

## API Endpoints

The server exposes the following endpoints:

- `POST /webhook/{path}`: Endpoints defined in your config for receiving events.
- `GET /health`: Health check (returns `200 OK`).
- `GET /webhooks`: Returns a list of supported webhook types.
- `GET /notifiers`: Returns a list of supported notifier types.

## Running the Application

**Using Make (Development):**

```bash
make run
```

**Using Binary (Production):**

```bash
./build/hookrelay --config config/production.yaml
```

## Setting up Sources

### GitHub

1. Go to Repo Settings -> Webhooks.
2. Payload URL: `http://your-server/webhook/github`.
3. Content type: `application/json` or `application/x-www-form-urlencoded`.
4. Secret: Must match the `secret` in `webhooks` config.

### Kanboard

1. Go to Project Settings -> Webhooks.
2. Webhook URL: `http://your-server/webhook/kanboard?token=YOUR_KANBOARD_TOKEN`.
3. Note: The token in the URL must match the `secret` in the config.

### Custom

1. Send `POST` to `http://your-server/webhook/custom`.
2. Header `X-Auth-Token`: Must match `secret` in config.
3. Body: Plain text or JSON.

## Testing

Use the provided `webhook.sh` script to simulate events:

```bash
# GitHub Push Event
./webhook.sh github push

# Kanboard Task Creation
./webhook.sh kanboard task.create

# Custom Message
./webhook.sh custom "Production server is unreachable!"
```

Run unit tests:

```bash
make test
```

## Project Structure

```
.
├── cmd/                # Entry point
├── config/             # Configuration files and structs
├── internal/
│   ├── adapters/       # Inbound (GitHub/Kanboard) & Outbound (TG/Email) logic
│   ├── app/            # App lifecycle and initialization
│   ├── core/           # Domain models and interfaces (ports)
│   ├── service/        # Business logic (Routing)
│   └── transport/      # HTTP server and middleware
└── templates/          # Embedded template files
```

## License

MIT License. See [LICENSE](LICENSE) for details.
