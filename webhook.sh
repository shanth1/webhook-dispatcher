#!/bin/bash

# ==============================================================================
# Webhook Test Script
#
# A script to test the webhook dispatcher service by emulating
# webhook requests from different providers.
#
# Dependencies: curl, openssl, jq
#
# Usage:
#   ./webhook.sh <mode> <arg2> [arg3]
#
# Modes:
#   - github: Emulates a GitHub webhook request.
#     - arg2: The GitHub event name (e.g., 'push', 'star', 'issues').
#     - arg3 (optional): Path to a JSON file containing the payload.
#                       If not provided, a default payload for the event is used.
#
#   - kanboard: Emulates a Kanboard webhook request.
#     - arg2: The Kanboard event name (e.g., 'task.create', 'task.move').
#     - arg3 (optional): Path to a JSON file containing the payload.
#
#   - custom: Sends a simple text-based webhook.
#     - arg2: The raw text message to send as the request body.
#
# Examples:
#   # Send a default GitHub 'push' event
#   ./webhook.sh github push
#
#   # Send a Kanboard 'task.create' event
#   ./webhook.sh kanboard task.create
#
#   # Send a GitHub 'issues' event using a custom payload from a file
#   ./webhook.sh github issues ./payloads/my_issue.json
#
#   # Send a custom alert message
#   ./webhook.sh custom "ALERT: Service is down!"
#
# ==============================================================================

# --- Script Configuration ---
# Ensure the script exits on any error
set -e
set -o pipefail

# --- Service Configuration ---
# These should match the values in your 'config/example.yaml'

# Base URL of the running service
BASE_URL="http://localhost:8080"

# GitHub Webhook Configuration
GITHUB_PATH="/webhook/github"
GITHUB_SECRET="your-super-secret-github-string"

# Kanboard Webhook Configuration
KANBOARD_PATH="/webhook/kanboard"
KANBOARD_SECRET="your-kanboard-secret-token"

# Custom Webhook Configuration
CUSTOM_PATH="/webhook/alerts"
CUSTOM_SECRET="secret-for-custom-alerts"

# --- Helper Functions ---

# Function to print usage instructions
print_usage() {
    echo "Usage: $0 <mode> <arg2> [arg3]"
    echo ""
    echo "Modes:"
    echo "  github <event_name> [json_file]      - Emulate a GitHub webhook."
    echo "  kanboard <event_name> [json_file]    - Emulate a Kanboard webhook."
    echo "  custom <'message'>                   - Send a custom plain text webhook."
    echo ""
    echo "Examples:"
    echo "  $0 github push"
    echo "  $0 kanboard task.create"
    echo "  $0 github issues ./payload.json"
    echo "  $0 custom 'This is a test alert'"
    exit 1
}

# Function to provide default JSON payloads for GitHub events
get_default_github_payload() {
    local event_name=$1
    case "$event_name" in
        "push")
            jq -n \
              '{ "ref": "refs/heads/main", "repository": { "full_name": "test/repo", "html_url": "https://github.com/test/repo" }, "sender": { "login": "testuser", "html_url": "https://github.com/testuser" }, "commits": [ { "id": "abc1234", "message": "feat: Add new feature", "author": { "name": "Test User" }, "url": "https://github.com/test/repo/commit/abc1234" } ] }'
            ;;
        "star")
            jq -n \
              '{ "action": "created", "repository": { "full_name": "test/repo", "stargazers_count": 101, "html_url": "https://github.com/test/repo" }, "sender": { "login": "another-user", "html_url": "https://github.com/another-user" } }'
            ;;
        *)
            echo "{\"message\":\"Default payload for event '$event_name' not found.\"}" | jq '.'
            ;;
    esac
}

# Function to provide default JSON payloads for Kanboard events
get_default_kanboard_payload() {
    local event_name=$1
    case "$event_name" in
        "task.create")
            jq -n \
              --arg event_name "$event_name" \
              '{ "event_name": $event_name, "event_data": { "task": { "id": "101", "title": "Implement new API endpoint", "project_name": "My Awesome Project", "creator_username": "admin", "url": "http://kanboard.local/?controller=TaskViewController&action=show&task_id=101" } } }'
            ;;
        "task.move")
            jq -n \
              --arg event_name "$event_name" \
              '{ "event_name": $event_name, "event_data": { "task": { "id": "102", "title": "Refactor database module", "project_name": "My Awesome Project", "column_name": "In Progress", "url": "http://kanboard.local/?controller=TaskViewController&action=show&task_id=102" } } }'
            ;;
        *)
            jq -n \
              --arg event_name "$event_name" \
              '{ "event_name": $event_name, "event_data": { "message": "Default payload" } }'
            ;;
    esac
}

# --- Main Logic ---

# Check for required dependencies
for cmd in curl openssl jq; do
    if ! command -v $cmd &> /dev/null; then
        echo "Error: Required command '$cmd' is not installed." >&2
        exit 1
    fi
done

# Check for minimum number of arguments
if [ "$#" -lt 2 ]; then
    echo "Error: Not enough arguments provided." >&2
    print_usage
fi

MODE=$1
shift

case "$MODE" in
    "github")
        EVENT_NAME=$1
        JSON_FILE=$2
        TARGET_URL="${BASE_URL}${GITHUB_PATH}"

        if [ -z "$EVENT_NAME" ]; then
            echo "Error: GitHub event name is required for 'github' mode." >&2
            print_usage
        fi

        echo "--- GitHub Mode ---"
        echo "Event: $EVENT_NAME"
        echo "URL:   $TARGET_URL"

        # Determine the JSON payload
        if [ -n "$JSON_FILE" ]; then
            if [ ! -f "$JSON_FILE" ]; then
                echo "Error: JSON file not found at '$JSON_FILE'" >&2
                exit 1
            fi
            echo "Payload: from file '$JSON_FILE'"
            JSON_PAYLOAD=$(cat "$JSON_FILE")
        else
            echo "Payload: using default for '$EVENT_NAME'"
            JSON_PAYLOAD=$(get_default_github_payload "$EVENT_NAME")
        fi

        # GitHub sends the payload as a URL-encoded form field `payload={...}`
        # We need to URL-encode the JSON string first.
        URL_ENCODED_JSON=$(echo -n "$JSON_PAYLOAD" | jq -sRr @uri)
        REQUEST_BODY="payload=${URL_ENCODED_JSON}"

        # Calculate the HMAC-SHA256 signature
        # The signature is calculated on the raw request body, not just the JSON.
        SIGNATURE="sha256=$(echo -n "$REQUEST_BODY" | openssl dgst -sha256 -hmac "$GITHUB_SECRET" -binary | xxd -p -c 256)"
        echo "Signature: $SIGNATURE"

        # Construct and execute the curl command
        CURL_CMD=(
            curl -v -X POST \
            -H "Content-Type: application/x-www-form-urlencoded" \
            -H "X-GitHub-Event: ${EVENT_NAME}" \
            -H "X-Hub-Signature-256: ${SIGNATURE}" \
            -d "${REQUEST_BODY}" \
            "${TARGET_URL}"
        )
        ;;

    "kanboard")
        EVENT_NAME=$1
        JSON_FILE=$2
        TARGET_URL="${BASE_URL}${KANBOARD_PATH}"

        if [ -z "$EVENT_NAME" ]; then
            echo "Error: Kanboard event name is required for 'kanboard' mode." >&2
            print_usage
        fi

        echo "--- Kanboard Mode ---"
        echo "Event: $EVENT_NAME"
        echo "URL:   $TARGET_URL"

        # Determine the JSON payload
        if [ -n "$JSON_FILE" ]; then
            if [ ! -f "$JSON_FILE" ]; then
                echo "Error: JSON file not found at '$JSON_FILE'" >&2
                exit 1
            fi
            echo "Payload: from file '$JSON_FILE'"
            JSON_PAYLOAD=$(cat "$JSON_FILE")
        else
            echo "Payload: using default for '$EVENT_NAME'"
            JSON_PAYLOAD=$(get_default_kanboard_payload "$EVENT_NAME")
        fi

        REQUEST_BODY="$JSON_PAYLOAD"
        SIGNATURE="$KANBOARD_SECRET"
        echo "Auth Token: $SIGNATURE"

        # Construct and execute the curl command
        CURL_CMD=(
            curl -v -X POST \
            -H "Content-Type: application/json" \
            -H "X-Kanboard-Token: ${SIGNATURE}" \
            -d "${REQUEST_BODY}" \
            "${TARGET_URL}"
        )
        ;;

    "custom")
        MESSAGE=$1
        TARGET_URL="${BASE_URL}${CUSTOM_PATH}"

        if [ -z "$MESSAGE" ]; then
            echo "Error: A message is required for 'custom' mode." >&2
            print_usage
        fi

        echo "--- Custom Mode ---"
        echo "URL:     $TARGET_URL"
        echo "Message: '$MESSAGE'"

        # The custom verifier uses a simple secret in the X-Auth-Token header
        REQUEST_BODY="$MESSAGE"
        SIGNATURE="$CUSTOM_SECRET"
        echo "Auth Header: $SIGNATURE"

        # Construct and execute the curl command
        CURL_CMD=(
            curl -v -X POST \
            -H "Content-Type: text/plain" \
            -H "X-Auth-Token: ${SIGNATURE}" \
            -d "${REQUEST_BODY}" \
            "${TARGET_URL}"
        )
        ;;

    *)
        echo "Error: Invalid mode '$MODE'." >&2
        print_usage
        ;;
esac

echo ""
echo "Executing command:"
# The next line prints the command in a copy-paste friendly format
printf '%q ' "${CURL_CMD[@]}"
echo ""
echo ""
echo "--- Server Response ---"

# Execute the command
"${CURL_CMD[@]}"
