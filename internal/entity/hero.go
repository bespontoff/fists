package entity

import "github.com/google/uuid"

type Hero struct {
	Name            string
	Health          int
	maxHealth       int
	recoveryRate    int
	CurrentBattleId uuid.UUID
}

func (h *Hero) DealDamage(i int) {
	h.Health -= i
	if h.Health < 0 {
		h.Health = 0
	}
}

func (h *Hero) IsDead() bool {
	return h.Health == 0
}

func NewHero(name string) *Hero {
	h := &Hero{
		Name:         name,
		Health:       10,
		maxHealth:    10,
		recoveryRate: 1,
	}
	return h
}
