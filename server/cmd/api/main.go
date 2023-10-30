package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"runtime/debug"

	"peterweightman.com/dgstats/internal/env"

	"github.com/lmittmann/tint"
)

const version = "0.0.1"

type config struct {
	httpPort int
	env      string
	baseURL  string
}

type application struct {
	config config
	logger *slog.Logger
}

func main() {
	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{Level: slog.LevelDebug}))

	err := run(logger)
	if err != nil {
		trace := string(debug.Stack())
		logger.Error(err.Error(), "trace", trace)
		os.Exit(1)
	}
}

func run(logger *slog.Logger) error {
	var cfg config

	cfg.baseURL = env.GetString("BASE_URL", "http://localhost:4444")
	cfg.httpPort = env.GetInt("HTTP_PORT", 4444)
	cfg.env = env.GetString("ENV", "development")

	showVersion := flag.Bool("version", false, "display version and exit")

	flag.Parse()

	if *showVersion {
		fmt.Printf("Version: %s\n", version)
		return nil
	}

	app := &application{
		config: cfg,
		logger: logger,
	}

	return app.serveHTTP()
}
