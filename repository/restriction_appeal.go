package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"
)

var (
	ErrAppealNotRestricted = errors.New("account is not restricted")
	ErrAppealDuplicate     = errors.New("pending appeal already exists")
)

type RestrictionAppealRow struct {
	AppealID   int64
	UserID     string
	Reason     string
	Status     string
	AdminNote  *string
	CreatedAt  time.Time
	ResolvedAt *time.Time
}

func (r user) IsUserRestricted(ctx context.Context, userID string) (bool, error) {
	var restricted bool
	err := r.bun.NewRaw(`
		SELECT (
			(restricted_until IS NOT NULL AND restricted_until > NOW())
			OR (posting_restricted_until IS NOT NULL AND posting_restricted_until > NOW())
		)
		FROM users WHERE user_id = ?
	`, userID).Scan(ctx, &restricted)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return restricted, err
}

func (r user) InsertRestrictionAppeal(ctx context.Context, userID, reason string) (int64, error) {
	restricted, err := r.IsUserRestricted(ctx, userID)
	if err != nil {
		return 0, err
	}
	if !restricted {
		return 0, ErrAppealNotRestricted
	}
	var pending bool
	err = r.bun.NewRaw(`
		SELECT EXISTS(
			SELECT 1 FROM user_restriction_appeals
			WHERE user_id = ? AND status = 'pending'
		)
	`, userID).Scan(ctx, &pending)
	if err != nil {
		return 0, err
	}
	if pending {
		return 0, ErrAppealDuplicate
	}
	var appealID int64
	err = r.bun.NewRaw(`
		INSERT INTO user_restriction_appeals (user_id, reason, status)
		VALUES (?, ?, 'pending')
		RETURNING appeal_id
	`, userID, reason).Scan(ctx, &appealID)
	return appealID, err
}

func (r user) GetRestrictionAppealForUser(ctx context.Context, userID string) (*RestrictionAppealRow, error) {
	row := new(RestrictionAppealRow)
	err := r.bun.NewRaw(`
		SELECT appeal_id, user_id, reason, status, admin_note, created_at, resolved_at
		FROM user_restriction_appeals
		WHERE user_id = ?
		ORDER BY
			CASE WHEN status = 'pending' THEN 0 ELSE 1 END,
			created_at DESC
		LIMIT 1
	`, userID).Scan(ctx, &row.AppealID, &row.UserID, &row.Reason, &row.Status, &row.AdminNote, &row.CreatedAt, &row.ResolvedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return row, nil
}

func (r user) HasPendingRestrictionAppeal(ctx context.Context, userID string) (bool, error) {
	var pending bool
	err := r.bun.NewRaw(`
		SELECT EXISTS(
			SELECT 1 FROM user_restriction_appeals
			WHERE user_id = ? AND status = 'pending'
		)
	`, userID).Scan(ctx, &pending)
	return pending, err
}

func scanAppealNote(note *string) string {
	if note == nil {
		return ""
	}
	return strings.TrimSpace(*note)
}

func AppealNoteString(note *string) string {
	return scanAppealNote(note)
}
