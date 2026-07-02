package entity

import "github.com/google/uuid"

type HitDirection string
type DefenceDirection string

const (
	HeadHit  HitDirection = "head"
	ChestHit HitDirection = "chest"
	LegsHit  HitDirection = "legs"
	ArmsHit  HitDirection = "arms"
	MissHit  HitDirection = "miss"

	HeadDefence  DefenceDirection = "head"
	ChestDefence DefenceDirection = "chest"
	LegsDefence  DefenceDirection = "legs"
	ArmsDefence  DefenceDirection = "arms"
	MissDefence  DefenceDirection = "miss"
)

type Hit struct {
	id           uuid.UUID
	roundId      uuid.UUID
	whoHit       *Hero
	toHit        *Hero
	whereHit     HitDirection
	whereDefence DefenceDirection
}

func NewHit(roundId uuid.UUID, whoHit, toHit *Hero, hitDirection HitDirection, defenceDirection DefenceDirection) *Hit {
	hit := &Hit{
		id:           uuid.New(),
		roundId:      roundId,
		whoHit:       whoHit,
		toHit:        toHit,
		whereHit:     hitDirection,
		whereDefence: defenceDirection,
	}
	return hit
}

func (hit *Hit) IsHitToTarget(other *Hit) bool {
	if string(hit.whereHit) == string(other.whereDefence) {
		return false
	} else {
		return true
	}
}
