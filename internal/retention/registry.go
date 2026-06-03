package retention

import (
	"context"
	"time"

	"github.com/rnikrozoft/pramool-core/internal/privacy"
	"github.com/rnikrozoft/pramool-core/repository"
)

const (
	RunnerInline = "inline"
	RunnerFuture = "future"
)

type Job struct {
	ID            string
	NameTH        string
	Description   string
	RetentionDays int
	BatchRunner   string
	Enabled       bool
	Run           func(ctx context.Context, deps *Deps) (int64, error)
}

type Deps struct {
	Privacy repository.Privacy
}

func Registry(deps *Deps) []Job {
	return []Job{
		{
			ID:            "tel_verify_stale",
			NameTH:        "ลบ tel_verify ค้าง",
			Description:   "ลบข้อมูลสมัครค้างที่ยังไม่สมัครครบและไม่มี user",
			RetentionDays: privacy.TelVerifyRetention,
			BatchRunner:   RunnerInline,
			Enabled:       true,
			Run: func(ctx context.Context, d *Deps) (int64, error) {
				if d == nil || d.Privacy == nil {
					return 0, nil
				}
				return d.Privacy.PurgeStaleTelVerify(ctx, privacy.TelVerifyRetention)
			},
		},
		{
			ID:            "user_notifications_expired",
			NameTH:        "ลบการแจ้งเตือนหมดอายุ",
			Description:   "ลบ user_notifications ที่ expires_at เก่ากว่า 90 วัน — รอ batch runner",
			RetentionDays: 90,
			BatchRunner:   RunnerFuture,
			Enabled:       true,
			Run:           nil,
		},
		{
			ID:            "consent_log_archive",
			NameTH:        "จัดเก็บ consent log เก่า",
			Description:   "ย้าย/ลบ consent_log เก่า — รอ batch runner",
			RetentionDays: 2555,
			BatchRunner:   RunnerFuture,
			Enabled:       false,
			Run:           nil,
		},
		{
			ID:            "access_log_files",
			NameTH:        "ลบ access log ไฟล์",
			Description:   "structured access log เก่ากว่า 12 เดือน — รอ batch runner / log infra",
			RetentionDays: 365,
			BatchRunner:   RunnerFuture,
			Enabled:       false,
			Run:           nil,
		},
		{
			ID:            "data_disclosure_log",
			NameTH:        "ลบ disclosure log เก่า",
			Description:   "ลบ data_disclosure_log เก่ากว่า 24 เดือน — รอ batch runner",
			RetentionDays: 730,
			BatchRunner:   RunnerFuture,
			Enabled:       false,
			Run:           nil,
		},
	}
}

func InlineJobs(deps *Deps) []Job {
	var out []Job
	for _, job := range Registry(deps) {
		if job.BatchRunner == RunnerInline && job.Enabled && job.Run != nil {
			out = append(out, job)
		}
	}
	return out
}

func FutureBatchJobs(deps *Deps) []Job {
	var out []Job
	for _, job := range Registry(deps) {
		if job.BatchRunner == RunnerFuture {
			out = append(out, job)
		}
	}
	return out
}

func RunInlineJob(ctx context.Context, jobID string, deps *Deps) (int64, error) {
	for _, job := range InlineJobs(deps) {
		if job.ID == jobID {
			return job.Run(ctx, deps)
		}
	}
	return 0, nil
}

func DefaultInterval() time.Duration {
	return 24 * time.Hour
}
