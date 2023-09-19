package models

import "time"

type ServerConfig struct {
	Port               int           `yaml:"Port"`
	ReadTimeout        time.Duration `yaml:"ReadTimeout"`
	WriteTimeout       time.Duration `yaml:"WriteTimeout"`
	MaxConnsPerIP      int           `yaml:"MaxConnsPerIP"`
	MaxRequestsPerConn int           `yaml:"MaxRequestsPerConn"`
	ReduceMemoryUsage  bool          `yaml:"ReduceMemoryUsage"`
	GetOnly            bool          `yaml:"GetOnly"`
}

type ApplicationConfig struct {
	MockFlag  bool `yaml:"MockFlag"`
	DebugMode bool `yaml:"DebugMode"`
	CacheSize int  `yaml:"CacheSize"`
	Workers   int  `yaml:"Workers"`
}

type LimiterConfig struct {
	Limit  int    `yaml:"Limit"`
	Period string `yaml:"Period"`
}

type Config struct {
	Server      ServerConfig      `yaml:"Server"`
	Application ApplicationConfig `yaml:"Application"`
	Limiter     LimiterConfig     `yaml:"Limiter"`
}
