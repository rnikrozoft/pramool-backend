package entity

import (
	"time"

	"github.com/uptrace/bun"
)

type ConsentLog struct {
	bun.BaseModel `bun:"table:consent_log"`
	ID            int64     `bun:"id,pk,autoincrement"`
	UserID        *string   `bun:"user_id"`
	Tel           *string   `bun:"tel"`
	ConsentType   string    `bun:"consent_type,notnull"`
	PolicyVersion string    `bun:"policy_version,notnull"`
	IPAddress     *string   `bun:"ip_address"`
	UserAgent     *string   `bun:"user_agent"`
	CreatedAt     time.Time `bun:"created_at,notnull,default:current_timestamp"`
}

type DSARRequest struct {
	bun.BaseModel      `bun:"table:dsar_requests"`
	ID                 int64      `bun:"id,pk,autoincrement"`
	UserID             string     `bun:"user_id,notnull"`
	RequestType        string     `bun:"request_type,notnull"`
	Status             string     `bun:"status,notnull"`
	UserNote           *string    `bun:"user_note"`
	AdminNote          *string    `bun:"admin_note"`
	HandledByAdminID   *int       `bun:"handled_by_admin_id"`
	DueAt              *time.Time `bun:"due_at"`
	DeletionExecutedAt *time.Time `bun:"deletion_executed_at"`
	CreatedAt          time.Time  `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt          time.Time  `bun:"updated_at,notnull,default:current_timestamp"`
	CompletedAt        *time.Time `bun:"completed_at"`
}
