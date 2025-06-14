package config

type AppConfigs struct {
	Database    DatabaseConfig
	Jwt         JwtConfig
	ThaiBulkSMS ThaiBulkSMS
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
	ExpireTime int
}

type ThaiBulkSMS struct {
	AddressRequest string
	AddressVerify  string
	APIKey         string
	APISecret      string
}
