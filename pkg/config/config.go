package config

import (
	"time"
)

type Config struct {
	Database  Database  `mapstructure:"database"`
	Scheduler Scheduler `mapstructure:"scheduler"`
}

type Scheduler struct {
	ReportInterval time.Duration `mapstructure:"report_interval"`
	UpdateInterval time.Duration `mapstructure:"update_interval"`
	MaxBurrowAge   int           `mapstructure:"max_burrow_age"`
	DepthIncrement float64       `mapstructure:"depth_increment"`
}
