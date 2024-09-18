package config

import (
	"flag"
	"runtime"
	"time"
)

type Config struct {
	Bandwidth       int
	Port            string
	LogsPath        string
	TimeForShutdown time.Duration
}

func Get() Config {
	//todo write usage and docs
	port := flag.String("port", "8080", "")
	logsPath := flag.String("logs_path", "logs.bin", "")
	timeForShutdown := flag.Duration("time_for_shutdown", 5*time.Minute, "")
	bandwidth := flag.Int("bandwidth", 10*runtime.NumCPU(), "")

	flag.Parse()

	return Config{
		*bandwidth,
		*port,
		*logsPath,
		*timeForShutdown,
	}
}
