package game

import (
	"context"
	"fists/config"
	"fists/internal/battleground"
	"fists/internal/entity"
	"log/slog"
	"time"
)

type Game struct {
	cfg      *config.Config
	log      *slog.Logger
	battleCh chan *entity.Battle
}

func New(cfg *config.Config, log *slog.Logger) *Game {
	g := &Game{
		cfg: cfg,
		log: log,
	}
	return g
}

func (g *Game) Run(ctx context.Context) {
	bg := battleground.NewBattleGround(g.log)
	go bg.StartEventLoop(ctx)

	bob := entity.NewHero("Bob")
	alice := entity.NewHero("Alice")

	tid := bg.RegisterDuel(bob, 1, time.Minute)
	duel, _ := bg.GetTicket(tid)
	duel.RegisterGamer(alice)
}
