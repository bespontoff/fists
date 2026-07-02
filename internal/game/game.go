package game

import (
	"context"
	"fists/config"
	"log/slog"
)

type Game struct {
	cfg *config.Config
	log *slog.Logger
}

func New(cfg *config.Config, log *slog.Logger) *Game {
	g := &Game{
		cfg: cfg,
		log: log,
	}
	return g
}

func (g *Game) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			g.log.Info("stopping game")
			return
		}
	}
}
