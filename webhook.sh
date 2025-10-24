#!/bin/bash

WEBHOOK_URL="http://localhost:8080/webhook"

WEBHOOK_SECRET='your-super-secret-webhook-string'
REQUEST_COUNT=1
DELAY=1

# --- (PAYLOAD) ---
read -r -d '' PAYLOAD <<EOF
{
  "ref": "refs/heads/main",
  "before": "0000000000000000000000000000000000000000",
  "after": "c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5",
  "repository": {
    "id": 12345678,
    "name": "my-test-repo",
    "full_name": "your-username/my-test-repo"
  },
  "pusher": {
    "name": "your-username",
    "email": "your-email@example.com"
  },
  "sender": {
    "login": "your-username",
    "id": 12345
  }
}
EOF

# --- SIGNATURE ---
SIGNATURE_256=$(echo -n "$PAYLOAD" | openssl dgst -sha256 -hmac "$WEBHOOK_SECRET" | sed 's/^.* //')

echo "Starting sending $REQUEST_COUNT test requests..."

# --- SENDING REQUESTS ---
for (( i=1; i<=$REQUEST_COUNT; i++ ))
do
  echo "----------------------------------------"
  echo "Sending request #$i..."

  curl --silent --show-error \
       -X POST \
       -H "Content-Type: application/json" \
       -H "X-GitHub-Event: push" \
       -H "X-Hub-Signature-256: sha256=$SIGNATURE_256" \
       -d "$PAYLOAD" \
       "$WEBHOOK_URL"

  echo
  echo "Request #$i sent"

  if [ $i -lt $REQUEST_COUNT ]; then
    sleep $DELAY
  fi
done

echo "----------------------------------------"
echo "Testing completed"
