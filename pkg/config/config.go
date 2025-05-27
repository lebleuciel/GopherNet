package config

import (
	"time"
)

type Config struct {
	Database  Database  `mapstructure:"database"`
	Scheduler Scheduler `mapstructure:"scheduler"`
	Logger    Logger    `mapstructure:"logger"`
}

type Scheduler struct {
	ReportInterval     time.Duration `mapstructure:"report_interval"`
	UpdateInterval     time.Duration `mapstructure:"update_interval"`
	MaxBurrowAge       int           `mapstructure:"max_burrow_age"`
	DepthIncrementRate float64       `mapstructure:"depth_increment"`
}

type Logger struct {
	Debug bool `mapstructure:"debug"`
}
