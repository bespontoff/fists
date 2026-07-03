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
	log.Info("Welcome to Fists Game")
	g := &Game{
		cfg: cfg,
		log: log,
	}
	return g
}

func (g *Game) Run(ctx context.Context) {
	g.log.Info("game started")
	bg := battleground.NewBattleGround(g.log)
	go bg.StartEventLoop(ctx)

	bob := entity.NewHero("Bob")
	alice := entity.NewHero("Alice")
	g.log.Debug("created heroes", "bob", bob, "alice", alice)

	tid := bg.RegisterDuel(bob, 1, time.Minute)
	duel, _ := bg.GetTicket(tid)
	duel.RegisterGamer(alice)

	time.Sleep(time.Second)
	battle, _ := bg.GetBattle(bob.CurrentBattleId)
	g.log.Debug("current battle", "battle", battle)

	// round 1
	currentRound := battle.CurrentRound()
	hitsCh := currentRound.GetHitsChannel()
	hitsCh <- entity.NewHit(currentRound.GetId(), bob, alice, entity.HeadHit, entity.HeadDefence)
	hitsCh <- entity.NewHit(currentRound.GetId(), alice, bob, entity.ChestHit, entity.ChestDefence)
	time.Sleep(time.Second)
	g.log.Debug("round heroes", "bob", bob, "alice", alice)

	select {
	case <-ctx.Done():
		g.log.Info("game context canceled")
		time.Sleep(time.Second)
	}
}
