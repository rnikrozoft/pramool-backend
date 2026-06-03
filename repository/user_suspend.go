package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"
)

func (r user) IsUserSuspended(ctx context.Context, subject string) (bool, error) {
	subject = strings.TrimSpace(subject)
	if subject == "" {
		return false, nil
	}
	var banned bool
	var err error
	if SubjectIsUserUUID(subject) {
		err = r.bun.NewRaw(`
			SELECT (suspended_at IS NOT NULL)
			FROM users
			WHERE user_id = ?
			LIMIT 1
		`, subject).Scan(ctx, &banned)
	} else {
		err = r.bun.NewRaw(`
			SELECT (suspended_at IS NOT NULL)
			FROM users
			WHERE TRIM(tel) = ?
			LIMIT 1
		`, subject).Scan(ctx, &banned)
	}
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return banned, err
}
