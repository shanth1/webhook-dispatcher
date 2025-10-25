package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

const (
	telegramBotToken = "YOUR_TELEGRAM_BOT_TOKEN"
	telegramChatID   = "YOUR_TELEGRAM_CHAT_ID"
)

type telegramMessage struct {
	ChatID string `json:"chat_id"`
	Text   string `json:"text"`
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	eventName := r.Header.Get("X-GitHub-Event")
	messageText := "Received a webhook event: " + eventName

	err := sendTelegramNotification(messageText)
	if err != nil {
		log.Printf("Failed to send Telegram notification: %v", err)
		http.Error(w, "Failed to send notification", http.StatusInternalServerError)
		return
	}

	log.Printf("Event '%s' successfully sent to Telegram.", eventName)
	w.WriteHeader(http.StatusOK)
}

func sendTelegramNotification(text string) error {
	apiURL := "https://api.telegram.org/bot" + telegramBotToken + "/sendMessage"

	msg := telegramMessage{
		ChatID: telegramChatID,
		Text:   text,
	}
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	_, err = http.Post(apiURL, "application/json", bytes.NewBuffer(body))
	return err
}

func main() {
	http.HandleFunc("/webhook", webhookHandler)
	log.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
