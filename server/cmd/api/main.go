package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"runtime/debug"
	"time"

	"peterweightman.com/runda/internal/database"
	"peterweightman.com/runda/internal/env"

	"github.com/lmittmann/tint"
)

const version = "0.0.1"

type config struct {
	httpPort int
	env      string
	baseURL  string
	db       struct {
		dsn          string
		automigrate  bool
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  time.Duration
		maxLifetime  time.Duration
	}
}

type application struct {
	config config
	logger *slog.Logger
	models database.Models
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

	flag.StringVar(&cfg.baseURL, "base-url", "http://localhost:4000", "Base URL")
	flag.IntVar(&cfg.httpPort, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", env.GetString("ENV", "development"), "Environment (development|staging|production) [env var: ENV]")

	flag.StringVar(&cfg.db.dsn, "db-dsn", env.GetString("DB_DSN", ""), "PostgreSQL DSN [env var: DB_DSN]")
	flag.BoolVar(&cfg.db.automigrate, "db-automigrate", env.GetBool("DB_AUTOMIGRATE", true), "Automigrate database [env var: DB_AUTOMIGRATE]")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", env.GetInt("DB_MAX_OPEN_CONNS", 25), "PostgreSQL max open connections [env var: DB_MAX_OPEN_CONNS]")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", env.GetInt("DB_MAX_IDLE_CONNS", 25), "PostgreSQL max idle connections [env var: DB_MAX_IDLE_CONNS]")
	flag.DurationVar(&cfg.db.maxIdleTime, "db-max-idle-time", env.GetDuration("DB_MAX_IDLE_TIME", time.Minute, 15), "PostgreSQL max connection idle time (mins) [env var: DB_MAX_IDLE_TIME]")
	flag.DurationVar(&cfg.db.maxLifetime, "db-max-lifetime", env.GetDuration("DB_MAX_LIFETIME", time.Hour, 2), "PostgreSQL max connection lifetime (hours) [env var: DB_MAX_IDLE_TIME]")

	showVersion := flag.Bool("version", false, "display version and exit")

	flag.Parse()

	if *showVersion {
		fmt.Printf("Version: %s\n", version)
		return nil
	}

	if cfg.db.dsn == "" {
		panic("DB_DSN environment variable is not set")
	}

	db, err := database.New(cfg.db.dsn, cfg.db.automigrate, database.DbPoolConfig{
		MaxOpenConns: cfg.db.maxOpenConns,
		MaxIdleConns: cfg.db.maxIdleConns,
		MaxIdleTime:  cfg.db.maxIdleTime,
		MaxLifetime:  cfg.db.maxLifetime,
	}, logger)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()
	logger.Info("database connection pool established")

	app := &application{
		config: cfg,
		logger: logger,
		models: database.NewModels(db),
	}

	return app.serveHTTP()
}
