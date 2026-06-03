package repository

import (
	"strings"

	"github.com/google/uuid"
)

// SubjectIsUserUUID is true when JWT subject is a persisted users.user_id (UUID), not a phone onboarding subject.
func SubjectIsUserUUID(subject string) bool {
	_, err := uuid.Parse(strings.TrimSpace(subject))
	return err == nil
}
