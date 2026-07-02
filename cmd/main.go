package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/feynmaz/farm/internal/config"
	"github.com/rs/zerolog/log"
)

// Tag is git tag set from Dockerfile
var Tag string

// Commit is git commit set from Dockerfile
var Commit string

func main() {
	version := getVersion()

	cfg, err := config.GetDefault()
	if err != nil {
		panic(fmt.Errorf("failed to get config: %w", err))
	}
	cfg.App.Version = version
	cfgContent, _ := json.Marshal(cfg)

	log.Debug().RawJSON("config", cfgContent).Send()

	_, initCancel := context.WithTimeout(context.Background(), cfg.App.InitTimeout)
	defer initCancel()

	// setup server

	// setup server end
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
