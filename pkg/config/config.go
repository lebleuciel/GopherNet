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

var DefaultScheduler = Scheduler{
	ReportInterval: 2 * time.Minute,
	UpdateInterval: 1 * time.Minute,
	MaxBurrowAge:   1440, // 24 hours in minutes
	DepthIncrement: 0.009,
}
