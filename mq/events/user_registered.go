package events

const UserRegisteredVersion = "v1"

type UserRegistered struct {
	Version   string `json:"version"`
	EventType string `json:"event_type"`
	UserID    int64  `json:"user_id"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	Timestamp int64  `json:"timestamp"`
}
