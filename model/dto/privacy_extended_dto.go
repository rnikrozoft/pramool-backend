package dto

type CookieConsentRequest struct {
	CookiePolicyVersion string `json:"cookie_policy_version" validate:"required"`
	AcceptEssential     bool   `json:"accept_essential" validate:"required,eq=true"`
	AcceptAnalytics     bool   `json:"accept_analytics"`
}

type MarketingConsentRequest struct {
	MarketingOptIn       bool   `json:"marketing_opt_in"`
	PrivacyPolicyVersion string `json:"privacy_policy_version" validate:"required"`
}

type MarketingConsentResponse struct {
	MarketingOptIn bool `json:"marketing_opt_in"`
}

type DataProcessorItem struct {
	Name           string `json:"name"`
	Purpose        string `json:"purpose"`
	DataCategories string `json:"data_categories"`
	Location       string `json:"location"`
	PrivacyURL     string `json:"privacy_url,omitempty"`
	DPAStatus      string `json:"dpa_status"`
}

type DataProcessorListResponse struct {
	Items []DataProcessorItem `json:"items"`
}

type AccountDeletionBlockersResponse struct {
	CanDelete              bool     `json:"can_delete"`
	AlreadyDeleted         bool     `json:"already_deleted"`
	CreditBalance          int64    `json:"credit_balance"`
	ActiveSellerAuctions   int64    `json:"active_seller_auctions"`
	ActiveBidAuctions      int64    `json:"active_bid_auctions"`
	PendingSellerShip      int64    `json:"pending_seller_ship"`
	PendingBuyerConfirm    int64    `json:"pending_buyer_confirm"`
	PendingWithdrawals     int64    `json:"pending_withdrawals"`
	Blockers               []string `json:"blockers"`
}

type ExecuteAccountDeletionRequest struct {
	DSARRequestID int64 `json:"dsar_request_id,omitempty"`
}

type ExecuteAccountDeletionResponse struct {
	OK      bool   `json:"ok"`
	UserID  string `json:"user_id"`
	Message string `json:"message"`
}

type RetentionJobItem struct {
	JobID            string `json:"job_id"`
	NameTH           string `json:"name_th"`
	Description      string `json:"description"`
	RetentionDays    int    `json:"retention_days"`
	BatchRunner      string `json:"batch_runner"`
	IsEnabled        bool   `json:"is_enabled"`
	LastRunAt        string `json:"last_run_at,omitempty"`
	LastDeletedCount int64  `json:"last_deleted_count,omitempty"`
	Implemented      bool   `json:"implemented"`
}

type RetentionJobListResponse struct {
	InlineJobs  []RetentionJobItem `json:"inline_jobs"`
	FutureJobs  []RetentionJobItem `json:"future_jobs"`
	BatchNote   string             `json:"batch_note"`
}
