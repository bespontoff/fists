package entity

import (
	"github.com/google/uuid"
	"time"
)

type Round struct {
	id        uuid.UUID
	battleId  uuid.UUID
	timer     *time.Timer
	hits      []*Hit
	createdAt time.Time
	closedAt  time.Time
	gamers    []*Hero
	hitsCh    chan *Hit
}

func NewRound(battleId uuid.UUID, timerDuration time.Duration, gamers []*Hero) *Round {
	r := &Round{
		id:       uuid.New(),
		battleId: battleId,
		timer:    time.NewTimer(timerDuration),
		gamers:   gamers,
		hits:     make([]*Hit, 0),
	}
	return r
}

func (r *Round) Start() {
	for {
		select {
		case <-r.timer.C:
			r.CalculateDamage()
			return
		case h := <-r.hitsCh:
			r.hits = append(r.hits, h)
			if len(r.hits) == len(r.gamers) {
				r.CalculateDamage()
				return
			}
		}
	}
}

func (r *Round) CalculateDamage() {
	if len(r.hits) == 0 {
		return
	}
	for _, hero := range r.gamers {
		var heroHit *Hit
		for _, h := range r.hits {
			if h.whoHit == hero {
				heroHit = h
			}
		}
		if heroHit == nil {
			heroHit = NewHit(r.id, hero, nil, MissHit, MissDefence)
		}

		for _, h := range r.hits {
			if h.toHit == hero {
				if h.IsHitToTarget(heroHit) {
					hero.DealDamage(1)
				}
			}
		}
	}
	return
}
