# github-webhook-dispatcher

[Russian Version](README.ru.md)

[![Build Status](https://img.shields.io/badge/build-passing-brightgreen)](https://github.com/shanth1/gitrelay)
[![Go Report Card](https://goreportcard.com/badge/github.com/shanth1/gitrelay)](https://goreportcard.com/report/github.com/shanth1/gitrelay)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A flexible Go-based service to receive GitHub webhooks and dispatch them as formatted notifications to various channels like Telegram and Email.

## Overview

`github-webhook-dispatcher` (also known as `gitrelay`) is a lightweight, configurable server that acts as a bridge between your GitHub repositories and your notification channels. It listens for webhook events from GitHub, verifies their authenticity, formats them into human-readable messages using Go templates, and broadcasts them to a pre-configured list of recipients via different services (e.g., Telegram, SMTP Email).

This allows you to create a centralized notification system for your team, keeping everyone updated on repository activities like pushes, pull requests, new issues, and more.

## Features

- **Secure Webhook Handling**: Verifies incoming webhooks using the `X-Hub-Signature-256` header to ensure they originate from GitHub.
- **Multi-Channel Notifications**: Out-of-the-box support for **Telegram** and **Email (SMTP)**. The architecture is extensible to support more services.
- **Customizable Message Templates**: Uses Go `html/template` to format notifications. Templates are provided for common events (`push`, `issues`, `pull_request`, etc.) and can be easily customized.
- **Flexible Routing**: Configure multiple "senders" (e.g., different Telegram bots or email accounts) and route notifications to various "recipients" with fine-grained control.
- **YAML-based Configuration**: All settings, including server address, webhook secrets, senders, and recipients, are managed in a single, easy-to-read YAML file.
- **Graceful Shutdown**: The server handles `SIGINT`/`SIGTERM` for a clean shutdown process.
- **Two Versions**:
  - `gitrelay`: The complete, feature-rich application that uses the YAML configuration.
  - `basic-webhook-handler`: A minimal example demonstrating the core webhook handling logic with hardcoded values.

## Prerequisites

- [Go](https://golang.org/dl/) (version 1.18 or higher is recommended)
- [Make](https://www.gnu.org/software/make/) for using the command shortcuts.
- (Optional) [golangci-lint](https://golangci-lint.run/usage/install/) for linting the code.

## Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/your-username/github-webhook-dispatcher.git
cd github-webhook-dispatcher
```

### 2. Install Dependencies

The project uses Go Modules. Tidy up the dependencies:

```bash
go mod tidy
```

Or using the Makefile:

```bash
make tidy
```

### 3. Build the Application

The `Makefile` provides convenient commands to build the binaries. The executables will be placed in the `build/` directory.

- **Build both applications (default):**
  ```bash
  make
  ```
- **Build only the complete application (`gitrelay`):**
  ```bash
  make build-complete
  ```
- **Build only the basic version (`basic-webhook-handler`):**
  ```bash
  make build-basic
  ```

## Configuration

The main application (`gitrelay`) is configured using a YAML file. A template is provided at `config/example.yaml`.

1.  **Copy the example configuration:**

    ```bash
    cp config/example.yaml config/production.yaml
    ```

2.  **Edit `config/production.yaml`:**

    ```yaml
    # The address and port the server will listen on.
    addr: ':8080'

    # The secret key used to verify webhook payloads from GitHub.
    # This MUST match the secret you set in your GitHub webhook settings.
    webhook_secret: 'your-super-secret-webhook-string'

    # Senders define the "how" and "from where" of notifications.
    senders:
      - name: 'project-telegram-bot' # A unique name for this sender
        type: 'telegram'
        settings:
          token: 'YOUR_TELEGRAM_BOT_TOKEN_1'

      - name: 'smtp-notifications'
        type: 'email'
        settings:
          host: 'smtp.example.com'
          port: 587
          username: 'user@example.com'
          password: 'your-smtp-password'
          from: 'GitHub Notifier <no-reply@example.com>'

    # Recipients define "who" receives the notifications.
    recipients:
      - name: 'Dev Team Channel (Telegram)'
        sender: 'project-telegram-bot' # Must match a sender name above
        target: '-100123456789' # Telegram Channel/Group Chat ID

      - name: 'Project Lead (Email)'
        sender: 'smtp-notifications'
        target: 'lead@example.com' # Email address
    ```

### Configuration Details

- `webhook_secret`: A crucial security feature. This string should be a high-entropy, randomly generated secret.
- `senders`: A list of notification services.
  - `name`: A unique identifier you'll use to link recipients to this sender.
  - `type`: Can be `telegram` or `email`.
  - `settings`: A map of key-value pairs specific to the sender type.
- `recipients`: A list of destinations for the notifications.
  - `sender`: The `name` of the sender to use for this recipient.
  - `target`: The destination identifier (e.g., a Telegram Chat ID or an email address).

## Running the Application

### 1. Run the Complete Application (`gitrelay`)

Pass the path to your configuration file using the `--config` flag.

```bash
./build/gitrelay --config config/production.yaml
```

You can also use the `make` shortcut, which runs the application with the example config file:

```bash
make run
```

The server will start and log a message like: `{"level":"info","service":"telehook","time":"...","message":"staring server on :8080"}`.

### 2. Set up the GitHub Webhook

1.  Navigate to your GitHub repository.
2.  Go to **Settings** > **Webhooks**.
3.  Click **Add webhook**.
4.  **Payload URL**: Enter the public URL where your application is running (e.g., `http://your-server-ip:8080/webhook`). For local testing, you might need a tool like [ngrok](https://ngrok.com/).
5.  **Content type**: Select `application/json`.
6.  **Secret**: Enter the exact same string you set for `webhook_secret` in your config file.
7.  **Which events would you like to trigger this webhook?**: Select the events you are interested in (e.g., "Pushes", "Issues", "Pull requests").
8.  Click **Add webhook**.

GitHub will send a `ping` event to your server to verify the connection.

## Testing

The project includes a comprehensive test script (`webhook.sh`) and standard Go tests.

### Running Go Unit Tests

To run the unit tests for the entire project:

```bash
make test
```

### Simulating GitHub Webhooks

The `webhook.sh` script can simulate sending various webhook events to a locally running server. This is extremely useful for testing your configuration and message templates without needing to perform actions on GitHub.

**Important**: Make sure the `WEBHOOK_SECRET` inside `webhook.sh` matches the `webhook_secret` in your YAML config.

- **Send all supported test events:**
  (First, ensure the `gitrelay` server is running in another terminal)

  ```bash
  make test-webhook
  ```

- **Send a specific event (e.g., `push`):**

  ```bash
  make test-webhook EVENT=push
  ```

- **Convenience shortcuts are also available:**
  ```bash
  make test-push
  make test-ping
  make test-issues
  ```

## Makefile Commands

The `Makefile` provides several helpful commands for development and maintenance.

| Command             | Description                                                    |
| ------------------- | -------------------------------------------------------------- |
| `make all` / `make` | Build both the complete and basic applications (default).      |
| `make build`        | Alias for `make all`.                                          |
| `make build-all`    | Explicitly builds both binaries.                               |
| `make run`          | Build and run the main application with `config/example.yaml`. |
| `make run-basic`    | Build and run the basic application.                           |
| `make test`         | Run all Go unit tests.                                         |
| `make test-webhook` | Send all test webhooks using `webhook.sh`.                     |
| `make test-push`    | Send a test `push` event.                                      |
| `make clean`        | Remove the `build/` directory and all compiled binaries.       |
| `make tidy`         | Run `go mod tidy`.                                             |
| `make fmt`          | Format all Go source files with `go fmt`.                      |
| `make vet`          | Run `go vet` to check for common issues.                       |
| `make lint`         | Run the `golangci-lint` linter (if installed).                 |

## Project Structure

```
.
├── build/                 # Compiled binaries (created by make)
├── cmd/
│   ├── basic/             # Source for the simple, hardcoded webhook handler
│   └── complete/          # Source for the main configurable application (gitrelay)
├── config/
│   └── example.yaml       # Example configuration file
├── internal/              # Internal Go packages for the main application
│   ├── app/               # Application setup and server lifecycle
│   ├── config/            # Configuration struct definitions
│   ├── handler/           # HTTP handlers for processing webhooks
│   ├── service/           # Notification services (Telegram, Email, etc.)
│   └── templates/         # Go template loading logic
├── templates/             # Go templates for formatting notification messages
│   ├── *.tmpl             # Template files for specific GitHub events
│   └── embed.go           # Embeds templates into the binary
├── Makefile               # Build, run, and test commands
├── go.mod                 # Go module definitions
└── webhook.sh             # Script for sending test webhooks
```

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change. Please make sure to update tests as appropriate.

## License

This project is licensed under the [MIT License](LICENSE).
