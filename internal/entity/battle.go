package entity

import (
	"github.com/google/uuid"
	"time"
)

type Status string
type HitDirection string
type DefenceDirection string

const (
	InProcessStatus Status = "InProcess"
	ClosedStatus    Status = "Closed"

	HeadHit  HitDirection = "head"
	ChestHit HitDirection = "chest"
	LegsHit  HitDirection = "legs"
	ArmsHit  HitDirection = "arms"

	HeadDefence  DefenceDirection = "head"
	ChestDefence DefenceDirection = "chest"
	LegsDefence  DefenceDirection = "legs"
	ArmsDefence  DefenceDirection = "arms"
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

type Round struct {
	id          uuid.UUID
	battleId    uuid.UUID
	timer       *time.Timer
	hits        []*hit
	createdAt   time.Time
	closedAt    time.Time
	gamers      []*Hero
	hitsCh      chan *hit
	pendingHits []*hit
}

type hit struct {
	id           uuid.UUID
	roundId      uuid.UUID
	whoHit       *Hero
	toHit        *Hero
	whereHit     HitDirection
	whereDefence DefenceDirection
}

func (b *Battle) Status() Status {
	return b.status
}

func NewBattle(gamers []*Hero, timerDuration time.Duration, bid int) *Battle {
	b := &Battle{
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
		round := NewRound(b.timerDuration, b.gamers)
		b.SetCurrentRound(round)
		round.Start()
	}
}

func NewRound(timerDuration time.Duration, gamers []*Hero) *Round {
	r := &Round{
		timer:       time.NewTimer(timerDuration),
		gamers:      gamers,
		hits:        make([]*hit, 0),
		pendingHits: make([]*hit, 0),
	}
	return r
}

func (r *Round) Start() {
	for {
		select {
		case <-r.timer.C:
			r.CalculateResult()
			return
		case h := <-r.hitsCh:
			r.pendingHits = append(r.pendingHits, h)
			if len(r.pendingHits) == len(r.gamers) {
				r.CalculateResult()
				return
			}
		}
	}
}

func (r *Round) CalculateResult() {
	return
}
