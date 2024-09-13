package config

import (
	"flag"
	"time"
)

type Config struct {
	Port           string
	LogsPath       string
	TimeToShutdown time.Duration
}

func Get() Config {
	port := flag.String("port", "8080", "")
	logsPath := flag.String("logs_path", "logs.bin", "")
	timeToShutdown := flag.Duration("time_to_shutdown", 10*time.Second, "")

	flag.Parse()

	return Config{*port, *logsPath, *timeToShutdown}
}
