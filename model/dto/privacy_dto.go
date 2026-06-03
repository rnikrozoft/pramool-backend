package dto

type ConsentPayload struct {
	PrivacyPolicyVersion string `json:"privacy_policy_version" validate:"required"`
	TermsVersion         string `json:"terms_version" validate:"required"`
	AcceptPrivacy        bool   `json:"accept_privacy" validate:"required,eq=true"`
	AcceptTerms          bool   `json:"accept_terms" validate:"required,eq=true"`
}

type PrivacyPolicyResponse struct {
	PrivacyPolicyVersion string `json:"privacy_policy_version"`
	TermsVersion         string `json:"terms_version"`
	UpdatedAt            string `json:"updated_at"`
	DPOEmail             string `json:"dpo_email"`
}

type CreateDSARRequest struct {
	RequestType string `json:"request_type" validate:"required,oneof=access delete correct"`
	Note        string `json:"note" validate:"omitempty,max=2000"`
}

type DSARRequestItem struct {
	ID          int64  `json:"id"`
	RequestType string `json:"request_type"`
	Status      string `json:"status"`
	UserNote    string `json:"user_note,omitempty"`
	AdminNote   string `json:"admin_note,omitempty"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	DueAt               string `json:"due_at,omitempty"`
	DeletionExecutedAt  string `json:"deletion_executed_at,omitempty"`
	ExportReady        bool   `json:"export_ready,omitempty"`
	CompletedAt         string `json:"completed_at,omitempty"`
}

type DSARRequestListResponse struct {
	Items []DSARRequestItem `json:"items"`
}
