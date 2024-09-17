package config

import (
	"flag"
	"runtime"
	"time"
)

type Config struct {
	Bandwidth      int
	Port           string
	LogsPath       string
	TimeToShutdown time.Duration
}

func Get() Config {
	port := flag.String("port", "8080", "")
	logsPath := flag.String("logs_path", "logs.bin", "")
	timeToShutdown := flag.Duration("time_to_shutdown", 5*time.Minute, "")
	bandwidth := flag.Int("bandwidth", 10*runtime.NumCPU(), "")

	flag.Parse()

	return Config{
		*bandwidth,
		*port,
		*logsPath,
		*timeToShutdown,
	}
}
