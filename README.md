# Webhook Dispatcher

[Russian Version](README.ru.md)

[![Build Status](https://img.shields.io/badge/build-passing-brightgreen)](https://github.com/shanth1/hookrelay)
[![Go Report Card](https://goreportcard.com/badge/github.com/shanth1/hookrelay)](https://goreportcard.com/report/github.com/shanth1/hookrelay)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A flexible Go-based service to receive webhooks from multiple sources like GitHub, Kanboard, and custom applications, and dispatch them as formatted notifications to various channels like Telegram and Email.

## Overview

`hookrelay` (`webhook-dispatcher`) is a lightweight, configurable server that acts as a bridge between your services and your notification channels. It listens for webhook events, verifies their authenticity using source-specific methods, formats them into human-readable messages using Go templates, and broadcasts them to a pre-configured list of recipients.

This allows you to create a centralized notification system for your team, keeping everyone updated on activities from various platforms in one place.

## Features

- **Multi-Source Webhook Handling**: Natively supports webhooks from **GitHub**, **Kanboard**, and **Custom** sources.
- **Secure & Flexible Verification**: Each webhook endpoint can have its own secret and verification method:
  - **GitHub**: HMAC signature (`X-Hub-Signature-256`).
  - **Kanboard**: URL query token.
  - **Custom**: Authentication header token (`X-Auth-Token`).
- **Multi-Channel Notifications**: Out-of-the-box support for **Telegram** and **Email (SMTP)**. The architecture is extensible to support more services.
- **Customizable Message Templates**: Uses Go `html/template` to format notifications. Templates are provided for common events and can be easily customized or extended for each source.
- **Advanced Routing**: A powerful YAML configuration allows you to define multiple webhook endpoints and route their notifications to specific recipients or groups of recipients.
- **YAML-based Configuration**: All settings, including server address, webhook endpoints, secrets, senders, and recipients, are managed in a single, easy-to-read YAML file.
- **Graceful Shutdown**: The server handles `SIGINT`/`SIGTERM` for a clean shutdown process.

## Prerequisites

- [Go](https://golang.org/dl/) (version 1.18 or higher is recommended)
- [Make](https://www.gnu.org/software/make/) for using the command shortcuts.
- For testing with the `webhook.sh` script:
  - [curl](https://curl.se/)
  - [jq](https://stedolan.github.io/jq/)
  - [openssl](https://www.openssl.org/)

## Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/your-username/webhook-dispatcher.git
cd webhook-dispatcher
```

### 2. Create Configuration Files

The project uses a `Makefile` to simplify setup. The following command copies the example configuration file into `config/local.yaml` (for development) and `config/production.yaml`.

```bash
make configs
```

### 3. Build the Application

The `Makefile` provides a convenient command to build the binary. The executable (`hookrelay`) will be placed in the `build/` directory.

```bash
make build
```

## Configuration

The application is configured using a YAML file (e.g., `config/production.yaml`). The configuration is split into three main sections: `webhooks`, `senders`, and `recipients`.

**Edit `config/production.yaml`:**

```yaml
# The address and port the server will listen on.
addr: ':8080'

env: 'local'

logger:
  level: 'info'
  app: 'my-app'
  service: 'my-service'
  udp_address: ''

# Define all incoming webhook endpoints here.
webhooks:
  - name: 'github-repo-events'
    path: '/webhook/github'
    type: 'github' # Supported types: github, kanboard, custom
    secret: 'your-super-secret-github-string' # Secret for HMAC signature
    recipients:
      - 'Dev Team (Telegram)' # A list of recipient names (defined below)
      - 'Admin (Email)'

  - name: 'kanboard-project-updates'
    path: '/webhook/kanboard'
    type: 'kanboard'
    secret: 'your-kanboard-secret-token' # Secret used as a URL token
    recipients:
      - 'Dev Team (Telegram)'

  - name: 'custom-service-alerts'
    path: '/webhook/custom'
    type: 'custom'
    secret: 'secret-for-custom-messages' # Secret for X-Auth-Token header
    recipients:
      - 'Admin (Email)'
      - 'On-call Engineer'

# Senders define the "how" and "from where" of notifications.
senders:
  - name: 'project-telegram-bot' # A unique name for this sender
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

# Recipients link a destination (target) with a sender.
recipients:
  - name: 'Dev Team (Telegram)'
    sender: 'project-telegram-bot' # Must match a sender name above
    target: '-100123456789' # Telegram Channel/Group Chat ID

  - name: 'Admin (Email)'
    sender: 'smtp-notifications'
    target: 'admin@example.com' # Email address

  - name: 'On-call Engineer'
    sender: 'project-telegram-bot'
    target: '987654321' # Another Telegram Chat ID
```

## Running the Application

Pass the path to your configuration file using the `--config` flag.

```bash
./build/hookrelay --config config/production.yaml
```

You can also use the `make` shortcut, which runs the application with `config/local.yaml`:

```bash
make run
```

The server will start and log a message like: `{"level":"info","time":"...","message":"staring server on :8080"}`.

## Setting up Webhooks

### GitHub

1.  Navigate to your GitHub repository **Settings > Webhooks > Add webhook**.
2.  **Payload URL**: `http://your-server.com:8080/webhook/github` (or the `path` you configured).
3.  **Content type**: `application/x-www-form-urlencoded` or `application/json`.
4.  **Secret**: Enter the exact same string you set for `secret` in the corresponding `webhooks` block in your config.

### Kanboard

1.  In Kanboard, go to **Project Settings > Webhooks**.
2.  **Webhook URL**: `http://your-server.com:8080/webhook/kanboard?token=your-kanboard-secret-token`. The `token` in the URL must match the `secret` from your config.

### Custom Source

1.  Configure your custom application to send a `POST` request to `http://your-server.com:8080/webhook/custom`.
2.  The request body should be `text/plain` containing the message you want to send.
3.  The request must include an `X-Auth-Token` header with the value of the `secret` from your config.

## Testing

The project includes a comprehensive test script (`webhook.sh`) and standard Go tests.

### Running Go Unit Tests

```bash
make test
```

### Simulating Webhooks

The `webhook.sh` script can simulate sending webhook events to a locally running server. This is extremely useful for testing your configuration and message templates.

**Important**: Make sure the secrets inside `webhook.sh` match the `secret` values in your YAML config.

- **Send a GitHub `push` event:**
  (First, ensure the `hookrelay` server is running in another terminal)

  ```bash
  ./webhook.sh github push
  ```

- **Send a Kanboard `task.create` event:**

  ```bash
  ./webhook.sh kanboard task.create
  ```

- **Send a custom alert:**

  ```bash
  ./webhook.sh custom "ALERT: The service is down!"
  ```

- **Use a custom JSON payload from a file for a GitHub `issues` event:**
  ```bash
  ./webhook.sh github issues ./payloads/my_issue.json
  ```

## Makefile Commands

| Command        | Description                                                         |
| :------------- | :------------------------------------------------------------------ |
| `make help`    | Display a list of available commands.                               |
| `make configs` | Create `local.yaml` and `production.yaml` configs from the example. |
| `make run`     | Run the application with the `config/local.yaml` file.              |
| `make build`   | Build the application binary into the `build/` directory.           |
| `make test`    | Run all Go unit tests.                                              |
| `make clean`   | Remove the `build/` directory and all compiled binaries.            |

## Project Structure

```
.
├── build/                 # Compiled binaries (created by make)
├── cmd/
│   └── main.go            # Main application entry point
├── config/
│   └── example.yaml       # Example configuration file
├── internal/              # Internal Go packages
│   ├── app/               # Application setup and server lifecycle
│   ├── config/            # Configuration struct definitions
│   ├── middleware/        # HTTP middleware (logging, verification)
│   ├── notifier/          # Notification sending logic (Telegram, Email)
│   ├── processor/         # Webhook payload processing (GitHub, Kanboard)
│   ├── server/            # HTTP server and routing
│   ├── templates/         # Go template loading logic
│   └── verifier/          # Webhook signature/token verification
├── templates/             # Go templates for formatting notifications
│   ├── github/*.tmpl      # Templates for GitHub events
│   ├── kanboard/*.tmpl    # Templates for Kanboard events
│   └── embed.go           # Embeds templates into the binary
├── Makefile               # Build, run, and test commands
├── go.mod                 # Go module definitions
└── webhook.sh             # Script for sending test webhooks
```

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change. Please make sure to update tests as appropriate.

## License

This project is licensed under the [MIT License](LICENSE).
