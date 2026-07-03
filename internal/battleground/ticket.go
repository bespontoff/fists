package battleground

import (
	"fists/internal/entity"
	"github.com/google/uuid"
	"log/slog"
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
	log              *slog.Logger
	id               uuid.UUID
	createdAt        time.Time
	creator          *entity.Hero
	duration         time.Duration
	battleType       BattleType
	players          []*entity.Hero
	status           TicketStatus
	timer            *time.Timer
	bid              int
	ticketEventsChan chan uuid.UUID
}

func NewTicket(log *slog.Logger, creator *entity.Hero, duration time.Duration, battleType BattleType, bid int, ticketEventsCh chan uuid.UUID) *Ticket {
	t := &Ticket{
		log:              log,
		id:               uuid.New(),
		createdAt:        time.Now(),
		creator:          creator,
		duration:         duration,
		battleType:       battleType,
		players:          make([]*entity.Hero, 0),
		status:           OpenStatus,
		timer:            time.NewTimer(duration),
		bid:              bid,
		ticketEventsChan: ticketEventsCh,
	}
	log.Info("created new ticket", "id", t.id)

	t.players = append(t.players, creator)
	return t
}

func (t *Ticket) ID() uuid.UUID {
	return t.id
}

func (t *Ticket) RegisterGamer(player *entity.Hero) {
	t.log.Info("registering new player", "name", player.Name)
	t.players = append(t.players, player)

	switch t.battleType {
	case DuelType:
		if len(t.players) == 2 {
			t.status = DoneStatus
			t.timer.Stop()
			t.log.Info("the duel ticket has gathered enough players", "players", len(t.players), "id", t.id)

			t.ticketEventsChan <- t.id
		}
	}
}

func (t *Ticket) Start() {
	t.log.Info("starting ticket", "id", t.id)
	for range t.timer.C {
		t.status = ClosedStatus
		t.log.Info("ticket is closed", "id", t.id)
	}
}
