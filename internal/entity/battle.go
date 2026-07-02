package entity

import (
	"github.com/google/uuid"
	"time"
)

type Status string

const (
	InProcessStatus Status = "InProcess"
	ClosedStatus    Status = "Closed"
)

type Battle struct {
	id            uuid.UUID
	status        Status
	rounds        []*Round
	currentRound  *Round
	createdAt     time.Time
	closedAt      time.Time
	gamers        []*Hero
	winner        *Hero
	timerDuration time.Duration
	bid           int
}

func (b *Battle) Status() Status {
	return b.status
}

func NewBattle(gamers []*Hero, timerDuration time.Duration, bid int) *Battle {
	b := &Battle{
		id:            uuid.New(),
		status:        InProcessStatus,
		gamers:        gamers,
		rounds:        make([]*Round, 0),
		timerDuration: timerDuration,
		bid:           bid,
		createdAt:     time.Now(),
	}
	return b
}

func (b *Battle) SetCurrentRound(round *Round) {
	b.currentRound = round
}

func (b *Battle) CurrentRound() *Round {
	return b.currentRound
}

func (b *Battle) Start() {
	for b.status == InProcessStatus {
		round := NewRound(b.id, b.timerDuration, b.gamers)
		b.rounds = append(b.rounds, round)
		b.SetCurrentRound(round)
		round.Start()
		var liveHeroes int
		var lastAliveHero *Hero
		for _, g := range b.gamers {
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
