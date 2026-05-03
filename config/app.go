package config

type AppConfigs struct {
	Database         DatabaseConfig
	Jwt              JwtConfig
	ThaiBulkSMS      ThaiBulkSMS
	CorsAllowOrigins string // comma-separated; e.g. http://localhost:3000,http://192.168.1.5:3000
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
