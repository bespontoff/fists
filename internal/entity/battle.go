package entity

import (
	"context"
	"github.com/google/uuid"
	"log/slog"
	"time"
)

type Status string

const (
	InProcessStatus Status = "InProcess"
	ClosedStatus    Status = "Closed"
)

type Battle struct {
	log           *slog.Logger
	id            uuid.UUID
	status        Status
	rounds        []*Round
	currentRound  *Round
	createdAt     time.Time
	closedAt      time.Time
	players       []*Hero
	winner        *Hero
	timerDuration time.Duration
	bid           int
}

func (b *Battle) Status() Status {
	return b.status
}

func NewBattle(log *slog.Logger, players []*Hero, timerDuration time.Duration, bid int) *Battle {
	b := &Battle{
		log:           log,
		id:            uuid.New(),
		status:        InProcessStatus,
		players:       players,
		rounds:        make([]*Round, 0),
		timerDuration: timerDuration,
		bid:           bid,
		createdAt:     time.Now(),
	}
	log.Info("battle created", "id", b.id)
	return b
}

func (b *Battle) SetCurrentRound(round *Round) {
	b.currentRound = round
	b.log.Info("current round has been set", "id", round.id)
}

func (b *Battle) CurrentRound() *Round {
	return b.currentRound
}

func (b *Battle) Start(ctx context.Context) {
	b.log.Info("starting battle", "id", b.id)
	for b.status == InProcessStatus {
		select {
		case <-ctx.Done():
			b.log.Info("battle context canceled", "id", b.id)
			return
		default:

		}

		round := NewRound(b.log, b.id, b.timerDuration, b.players)
		b.rounds = append(b.rounds, round)
		b.SetCurrentRound(round)
		round.Start(ctx)
		var liveHeroes int
		var lastAliveHero *Hero
		for _, g := range b.players {
			if !g.IsDead() {
				liveHeroes++
				lastAliveHero = g
			}
		}
		if liveHeroes == 0 {
			b.status = ClosedStatus
		}
		if liveHeroes == 1 {
			b.status = ClosedStatus
			b.winner = lastAliveHero
		}
	}
}

func (b *Battle) GetId() uuid.UUID {
	return b.id
}

func (b *Battle) GetIdString() string {
	return b.id.String()
}
