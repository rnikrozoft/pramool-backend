package dto

const NotificationAutoDeleteNote = "การแจ้งเตือนนี้จะถูกลบออกอัตโนมัติภายใน 3 วันนับจากวันที่คุณเปิดอ่าน"

type UserNotificationItem struct {
	NotificationID int64  `json:"notification_id"`
	Kind           string `json:"kind"`
	Title          string `json:"title"`
	Body           string `json:"body"`
	Read           bool   `json:"read"`
	ReadAt         string `json:"read_at,omitempty"`
	ExpiresAt      string `json:"expires_at"`
	CreatedAt      string `json:"created_at"`
	AutoDeleteNote string `json:"auto_delete_note,omitempty"`
}

type UserNotificationListResponse struct {
	Items []UserNotificationItem `json:"items"`
	Total int64                  `json:"total"`
}

type UnreadNotificationCountResponse struct {
	Count int `json:"count"`
}

type MarkNotificationReadResponse struct {
	ReadAt         string `json:"read_at"`
	ExpiresAt      string `json:"expires_at"`
	AutoDeleteNote string `json:"auto_delete_note"`
}
