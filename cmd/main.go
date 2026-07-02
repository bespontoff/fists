package main

import (
	"context"
	"fists/config"
	"fists/internal/game"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.New()
	logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.LogLevel})
	log := slog.New(logHandler)

	gracefulCh := make(chan os.Signal, 1)
	signal.Notify(gracefulCh, syscall.SIGTERM, syscall.SIGINT)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		sig := <-gracefulCh
		log.Info("graceful shutdown", "signal", sig)
		cancel()
	}()

	g := game.New(cfg, log)
	g.Run(ctx)
}
