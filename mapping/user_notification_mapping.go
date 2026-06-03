package mapping

import (
	"time"

	"github.com/rnikrozoft/pramool-core/model/dto"
	"github.com/rnikrozoft/pramool-core/repository"
)

func ToUserNotificationItem(row repository.UserNotificationRow) dto.UserNotificationItem {
	item := dto.UserNotificationItem{
		NotificationID: row.NotificationID,
		Kind:           row.Kind,
		Title:          row.Title,
		Body:           row.Body,
		Read:           row.ReadAt != nil,
		CreatedAt:      row.CreatedAt.Format(time.RFC3339),
	}
	if row.ReadAt != nil {
		item.ReadAt = row.ReadAt.Format(time.RFC3339)
		item.ExpiresAt = row.ExpiresAt.Format(time.RFC3339)
		item.AutoDeleteNote = dto.NotificationAutoDeleteNote
	}
	return item
}
