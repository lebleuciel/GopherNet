package config

type Config struct {
	Database Database `mapstructure:"database"`
}

// Database holds the database configuration
type Database struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}
