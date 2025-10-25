#!/bin/bash

# --- CONFIGURATION ---
WEBHOOK_URL="http://localhost:8080/webhook"
WEBHOOK_SECRET='your-super-secret-webhook-string'
DELAY=1 # Delay between requests in seconds

# --- HELPER FUNCTION FOR URL ENCODING ---
urlencode() {
    local lang=C  # Используем стандартную локаль
    local length="${#1}"
    for (( i = 0; i < length; i++ )); do
        local c="${1:i:1}"
        case $c in
            [a-zA-Z0-9.~_-]) printf "%s" "$c" ;;
            *) printf '%%%02X' "'$c" ;; # Кодируем все остальные символы
        esac
    done
}


# --- FUNCTION TO SEND REQUEST ---
send_request() {
  local EVENT_TYPE=$1
  local PAYLOAD=$2

  # Формируем тело запроса для application/x-www-form-urlencoded
  local FORM_BODY="payload=$(urlencode "$PAYLOAD")"

  # Генерируем подпись на основе ИСХОДНОГО JSON-пейлоада
  local SIGNATURE_256=$(echo -n "$PAYLOAD" | openssl dgst -sha256 -hmac "$WEBHOOK_SECRET" | sed 's/^.* //')

  echo "----------------------------------------"
  echo "Sending event: $EVENT_TYPE"

  # Отправляем CURL-запрос с новым Content-Type и телом
  curl --silent --show-error \
       -X POST \
       -H "Content-Type: application/x-www-form-urlencoded" \
       -H "X-GitHub-Event: $EVENT_TYPE" \
       -H "X-Hub-Signature-256: sha256=$SIGNATURE_256" \
       -d "$FORM_BODY" \
       "$WEBHOOK_URL"

  echo
  echo "Event '$EVENT_TYPE' sent."
}

# --- PAYLOADS FOR DIFFERENT EVENTS (Без изменений) ---

# PING EVENT PAYLOAD
get_ping_payload() {
  read -r -d '' PAYLOAD <<EOF
{
  "zen": "Approachable is better than simple.",
  "hook_id": 12345678,
  "eventName": "ping",
  "repository": { "full_name": "your-username/my-test-repo" },
  "sender": { "login": "your-username" }
}
EOF
  echo "$PAYLOAD"
}

# PUSH EVENT PAYLOAD
get_push_payload() {
  read -r -d '' PAYLOAD <<EOF
{
  "ref": "refs/heads/main",
  "after": "c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5",
  "repository": { "full_name": "your-username/my-test-repo" },
  "pusher": { "name": "your-username" },
  "sender": { "login": "your-username" },
  "commits": [{ "message": "feat: Add new feature" }]
}
EOF
  echo "$PAYLOAD"
}

# ISSUES EVENT PAYLOAD
get_issues_payload() {
  read -r -d '' PAYLOAD <<EOF
{
  "action": "opened",
  "issue": { "number": 13, "title": "Found a bug" },
  "repository": { "full_name": "your-username/my-test-repo" },
  "sender": { "login": "your-username" }
}
EOF
  echo "$PAYLOAD"
}

# PULL_REQUEST EVENT PAYLOAD
get_pull_request_payload() {
  read -r -d '' PAYLOAD <<EOF
{
  "action": "opened",
  "number": 14,
  "pull_request": { "title": "Update documentation", "user": { "login": "your-username" } },
  "repository": { "full_name": "your-username/my-test-repo" },
  "sender": { "login": "your-username" }
}
EOF
  echo "$PAYLOAD"
}

# ISSUE_COMMENT EVENT PAYLOAD
get_issue_comment_payload() {
  read -r -d '' PAYLOAD <<EOF
{
  "action": "created",
  "issue": { "number": 13 },
  "comment": { "body": "Thanks for reporting!" },
  "repository": { "full_name": "your-username/my-test-repo" },
  "sender": { "login": "your-username" }
}
EOF
  echo "$PAYLOAD"
}

# RELEASE EVENT PAYLOAD
get_release_payload() {
  read -r -d '' PAYLOAD <<EOF
{
  "action": "published",
  "release": { "tag_name": "v1.0.0", "name": "Version 1.0.0" },
  "repository": { "full_name": "your-username/my-test-repo" },
  "sender": { "login": "your-username" }
}
EOF
  echo "$PAYLOAD"
}

# FORK EVENT PAYLOAD
get_fork_payload() {
  read -r -d '' PAYLOAD <<EOF
{
  "forkee": { "full_name": "another-user/my-test-repo" },
  "repository": { "full_name": "your-username/my-test-repo" },
  "sender": { "login": "another-user" }
}
EOF
  echo "$PAYLOAD"
}

# STAR EVENT PAYLOAD
get_star_payload() {
  read -r -d '' PAYLOAD <<EOF
{
  "action": "created",
  "repository": { "full_name": "your-username/my-test-repo" },
  "sender": { "login": "new-stargazer" }
}
EOF
  echo "$PAYLOAD"
}

# --- MAIN LOGIC (Без изменений) ---
main() {
  EVENT_TO_SEND=$1

  if [ -n "$EVENT_TO_SEND" ]; then
    # Send a specific event if one is provided
    case $EVENT_TO_SEND in
      ping) send_request "ping" "$(get_ping_payload)" ;;
      push) send_request "push" "$(get_push_payload)" ;;
      issues) send_request "issues" "$(get_issues_payload)" ;;
      pull_request) send_request "pull_request" "$(get_pull_request_payload)" ;;
      issue_comment) send_request "issue_comment" "$(get_issue_comment_payload)" ;;
      release) send_request "release" "$(get_release_payload)" ;;
      fork) send_request "fork" "$(get_fork_payload)" ;;
      star) send_request "star" "$(get_star_payload)" ;;
      *) echo "Error: Unknown event type '$EVENT_TO_SEND'" && exit 1 ;;
    esac
  else
    # Send all event types sequentially
    echo "Starting to send all test events..."
    send_request "ping" "$(get_ping_payload)" && sleep $DELAY
    send_request "push" "$(get_push_payload)" && sleep $DELAY
    send_request "issues" "$(get_issues_payload)" && sleep $DELAY
    send_request "pull_request" "$(get_pull_request_payload)" && sleep $DELAY
    send_request "issue_comment" "$(get_issue_comment_payload)" && sleep $DELAY
    send_request "release" "$(get_release_payload)" && sleep $DELAY
    send_request "fork" "$(get_fork_payload)" && sleep $DELAY
    send_request "star" "$(get_star_payload)"
  fi

  echo "----------------------------------------"
  echo "Testing completed"
}

main "$@"
