package config

type Database struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

var DefaultDatabase = Database{
	Host:     "localhost",
	Port:     5432,
	User:     "postgres",
	Password: "postgres",
	Database: "gophernet",
}
