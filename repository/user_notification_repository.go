package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/uptrace/bun"
)

const notificationReadRetention = 3 * 24 * time.Hour

type UserNotificationRow struct {
	NotificationID int64
	UserID         string
	Kind           string
	Title          string
	Body           string
	ReadAt         *time.Time
	ExpiresAt      time.Time
	CreatedAt      time.Time
}

type UserNotificationRepository interface {
	DeleteExpired(ctx context.Context) error
	CountUnread(ctx context.Context, userID string) (int, error)
	List(ctx context.Context, userID string, limit, offset int) ([]UserNotificationRow, int64, error)
	MarkRead(ctx context.Context, userID string, notificationID int64) (*UserNotificationRow, error)
}

type userNotificationRepository struct {
	bun *bun.DB
}

func NewUserNotificationRepository(bun *bun.DB) UserNotificationRepository {
	return userNotificationRepository{bun: bun}
}

func (r userNotificationRepository) DeleteExpired(ctx context.Context) error {
	_, err := r.bun.NewRaw(`
		DELETE FROM user_notifications
		WHERE read_at IS NOT NULL AND expires_at <= NOW()
	`).Exec(ctx)
	return err
}

func (r userNotificationRepository) CountUnread(ctx context.Context, userID string) (int, error) {
	if err := r.DeleteExpired(ctx); err != nil {
		return 0, err
	}
	userID = strings.TrimSpace(userID)
	var n int
	err := r.bun.NewRaw(`
		SELECT COUNT(*)::int
		FROM user_notifications
		WHERE user_id = ? AND read_at IS NULL
	`, userID).Scan(ctx, &n)
	return n, err
}

func (r userNotificationRepository) List(ctx context.Context, userID string, limit, offset int) ([]UserNotificationRow, int64, error) {
	if err := r.DeleteExpired(ctx); err != nil {
		return nil, 0, err
	}
	userID = strings.TrimSpace(userID)
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	var total int64
	if err := r.bun.NewRaw(`
		SELECT COUNT(*)::bigint
		FROM user_notifications
		WHERE user_id = ? AND (read_at IS NULL OR expires_at > NOW())
	`, userID).Scan(ctx, &total); err != nil {
		return nil, 0, err
	}

	rows, err := r.bun.QueryContext(ctx, `
		SELECT notification_id, user_id, kind, title, body, read_at, expires_at, created_at
		FROM user_notifications
		WHERE user_id = ? AND (read_at IS NULL OR expires_at > NOW())
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var out []UserNotificationRow
	for rows.Next() {
		var row UserNotificationRow
		if err := rows.Scan(
			&row.NotificationID, &row.UserID, &row.Kind, &row.Title, &row.Body,
			&row.ReadAt, &row.ExpiresAt, &row.CreatedAt,
		); err != nil {
			return nil, 0, err
		}
		out = append(out, row)
	}
	return out, total, rows.Err()
}

func (r userNotificationRepository) MarkRead(ctx context.Context, userID string, notificationID int64) (*UserNotificationRow, error) {
	if err := r.DeleteExpired(ctx); err != nil {
		return nil, err
	}
	userID = strings.TrimSpace(userID)
	now := time.Now()
	expiresAt := now.Add(notificationReadRetention)

	res, err := r.bun.NewRaw(`
		UPDATE user_notifications
		SET read_at = COALESCE(read_at, ?),
		    expires_at = CASE WHEN read_at IS NULL THEN ? ELSE expires_at END
		WHERE notification_id = ? AND user_id = ? AND (read_at IS NULL OR expires_at > NOW())
	`, now, expiresAt, notificationID, userID).Exec(ctx)
	if err != nil {
		return nil, err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, sql.ErrNoRows
	}

	row := new(UserNotificationRow)
	err = r.bun.NewRaw(`
		SELECT notification_id, user_id, kind, title, body, read_at, expires_at, created_at
		FROM user_notifications
		WHERE notification_id = ? AND user_id = ?
	`, notificationID, userID).Scan(ctx, row)
	if err != nil {
		return nil, err
	}
	return row, nil
}

var ErrNotificationNotFound = errors.New("notification not found")
