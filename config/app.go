package config

type AppConfigs struct {
	Database         DatabaseConfig
	Jwt              JwtConfig
	ThaiBulkSMS      ThaiBulkSMS
	CorsAllowOrigins     string
	TrackingMoreAPIKey   string
	NationalIDEncKey     string
}

type DatabaseConfig struct {
	Host         string
	Port         string
	Username     string
	Password     string
	DatabaseName string
}

type JwtConfig struct {
	Issuer     string
	Secret     string
	ExpireTime int // access token TTL (hours)
	// RefreshExpireTime is refresh token TTL (hours). If 0, defaults to 168 (7d) in cmd/root.
	RefreshExpireTime int
}

type ThaiBulkSMS struct {
	AddressRequest string
	AddressVerify  string
	APIKey         string
	APISecret      string
}
