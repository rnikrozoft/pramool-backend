package config

type AppConfigs struct {
	Database DatabaseConfig
	Jwt      JwtConfig
}

type DatabaseConfig struct {
	Host         string
	Port         string
	Username     string
	Password     string
	DatabaseName string
}

type JwtConfig struct {
	Secret     string
	ExpireTime int
}
