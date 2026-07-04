package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"

	"github.com/feynmaz/farm/internal/config"
	"github.com/feynmaz/farm/internal/logger"
	"github.com/feynmaz/farm/internal/server"
)

// Tag is git tag set from Dockerfile
var Tag string

// Commit is git commit set from Dockerfile
var Commit string

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	version := getVersion()

	cfg, err := config.GetDefault()
	if err != nil {
		return fmt.Errorf("failed to get config: %w", err)
	}
	cfg.App.Version = version
	cfgContent, _ := json.Marshal(cfg)

	l, err := logger.New(cfg.App.Name, cfg.App.Env, cfg.Log.Level)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}

	l.Debug().RawJSON("config", cfgContent).Send()

	srv := server.New(cfg, l)

	signalCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	g, ctx := errgroup.WithContext(signalCtx)

	g.Go(func() error {
		if err := srv.Run(ctx); err != nil {
			return fmt.Errorf("server error: %w", err)
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		l.Error().Err(err).Msg("application error")
		return err
	}

	l.Info().Msg("server stopped")
	return nil
}

func getVersion() string {
	tag, commit := Tag, Commit

	if Tag == "" {
		tag = "tag"
	}
	if Commit == "" {
		commit = "commit"
	}
	return tag + "-" + commit
}
