package entity

import (
	"context"
	"github.com/google/uuid"
	"log/slog"
	"time"
)

type Round struct {
	log       *slog.Logger
	id        uuid.UUID
	battleId  uuid.UUID
	timer     *time.Timer
	hits      []*Hit
	createdAt time.Time
	closedAt  time.Time
	gamers    []*Hero
	hitsCh    chan *Hit
}

func NewRound(log *slog.Logger, battleId uuid.UUID, timerDuration time.Duration, gamers []*Hero) *Round {
	r := &Round{
		log:      log,
		id:       uuid.New(),
		battleId: battleId,
		timer:    time.NewTimer(timerDuration),
		gamers:   gamers,
		hits:     make([]*Hit, 0),
	}
	log.Info("round created", "id", r.id)
	return r
}

func (r *Round) Start(ctx context.Context) {
	r.log.Info("round started", "id", r.id)
	for {
		select {
		case <-r.timer.C:
			r.log.Info("round timeout", "id", r.id)
			r.CalculateDamage()
			return
		case h := <-r.hitsCh:
			r.log.Info("round got hit", "id", r.id, "hit", h)
			r.hits = append(r.hits, h)
			if len(r.hits) == len(r.gamers) {
				r.log.Info("round finished", "id", r.id, "hits", r.hits)
				r.CalculateDamage()
				return
			}
		case <-ctx.Done():
			r.log.Info("round cancelled", "id", r.id)
			return
		}
	}
}

func (r *Round) CalculateDamage() {
	r.log.Info("round calculated damage", "id", r.id)
	//if len(r.hits) == 0 {
	//	r.log.Info("round hits is empty", "id", r.id)
	//	return
	//}
	for _, hero := range r.gamers {
		var heroHit *Hit
		for _, h := range r.hits {
			if h.whoHit == hero {
				heroHit = h
			}
		}
		if heroHit == nil {
			r.log.Info("the hero did not land a blow in this round", "id", r.id, "hero", hero)
			heroHit = NewHit(r.id, hero, nil, MissHit, MissDefence)
		}

		for _, h := range r.hits {
			if h.toHit == hero {
				if h.IsHitToTarget(heroHit) {
					r.log.Info("hero hit to target", "id", r.id, "hero", hero, "hit", h)
					hero.DealDamage(1)
				}
			}
		}
	}
	return
}
