package battleground

import (
	"fists/internal/entity"
	"github.com/google/uuid"
	"time"
)

type BattleType string

type TicketStatus string

const (
	DuelType BattleType = "duel"

	OpenStatus   TicketStatus = "open"
	ClosedStatus TicketStatus = "closed"
	DoneStatus   TicketStatus = "done"
)

type Ticket struct {
	id               uuid.UUID
	createdAt        time.Time
	creator          *entity.Hero
	duration         time.Duration
	battleType       BattleType
	gamers           []*entity.Hero
	status           TicketStatus
	timer            *time.Timer
	bid              int
	ticketEventsChan chan uuid.UUID
}

func NewTicket(creator *entity.Hero, duration time.Duration, battleType BattleType, bid int, ticketEventsCh chan uuid.UUID) *Ticket {
	t := &Ticket{
		id:               uuid.New(),
		createdAt:        time.Now(),
		creator:          creator,
		duration:         duration,
		battleType:       battleType,
		gamers:           make([]*entity.Hero, 0),
		status:           OpenStatus,
		timer:            time.NewTimer(duration),
		bid:              bid,
		ticketEventsChan: ticketEventsCh,
	}
	t.gamers = append(t.gamers, creator)
	return t
}

func (t *Ticket) ID() uuid.UUID {
	return t.id
}

func (t *Ticket) RegisterGamer(gamer *entity.Hero) {
	t.gamers = append(t.gamers, gamer)

	switch t.battleType {
	case DuelType:
		if len(t.gamers) == 2 {
			t.status = DoneStatus
			t.timer.Stop()
		}
	}
}

func (t *Ticket) Start() {
	for range t.timer.C {
		t.status = ClosedStatus
	}
}
