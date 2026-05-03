package model

import "github.com/golang-jwt/jwt/v5"

// TokenUseAccess / TokenUseRefresh distinguish short-lived access vs long-lived refresh JWTs.
const (
	TokenUseAccess  = "access"
	TokenUseRefresh = "refresh"
)

type CustomClaims struct {
	UserID   string `json:"user_id,omitempty"`
	LoggedIn bool   `json:"logged_in,omitempty"`
	TokenUse string `json:"token_use,omitempty"` // "access" | "refresh"; omit => treated as access (legacy)
	jwt.RegisteredClaims
}
