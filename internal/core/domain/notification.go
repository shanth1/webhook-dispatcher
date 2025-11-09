package domain

// Notification is a general message that needs to be sent
type Notification struct {
	Title     string
	Body      string
	ParseMode string
}
