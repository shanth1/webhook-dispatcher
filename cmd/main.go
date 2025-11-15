package main

import (
	"flag"
	"time"

	"github.com/shanth1/gotools/conf"
	"github.com/shanth1/gotools/consts"
	"github.com/shanth1/gotools/ctx"
	"github.com/shanth1/gotools/flags"
	"github.com/shanth1/gotools/log"
	"github.com/shanth1/hookrelay/internal/app"
	"github.com/shanth1/hookrelay/internal/config"
)

type Flags struct {
	ConfigPath string `flag:"config" usage:"Path to the YAML config file"`
}

func main() {
	ctx, shutdownCtx, cancel, shutdownCancel := ctx.WithGracefulShutdown(10 * time.Second)
	defer cancel()
	defer shutdownCancel()

	logger := log.New()

	flagCfg := &Flags{}
	if err := flags.RegisterFromStruct(flagCfg); err != nil {
		logger.Fatal().Err(err).Msg("register flags")
	}
	flag.Parse()

	cfg := &config.Config{}
	if err := conf.Load(flagCfg.ConfigPath, cfg); err != nil {
		logger.Fatal().Err(err).Msg("load config")
	}

	logger = logger.WithOptions(log.WithConfig(log.Config{
		Level:        cfg.Logger.Level,
		App:          cfg.Logger.App,
		Service:      cfg.Logger.Service,
		UDPAddress:   cfg.Logger.UDPAddress,
		EnableCaller: cfg.Env != consts.EnvProd,
		Console:      cfg.Env != consts.EnvProd,
		JSONOutput:   cfg.Env == consts.EnvProd,
	}))

	ctx = log.NewContext(ctx, logger)
	app.Run(ctx, shutdownCtx, cfg)
}
