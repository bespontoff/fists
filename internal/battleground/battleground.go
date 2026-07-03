package battleground

import (
	"context"
	"fists/internal/entity"
	"fmt"
	"github.com/google/uuid"
	"log/slog"
	"time"
)

type BattleGround struct {
	log               *slog.Logger
	registeredBattles map[uuid.UUID]*entity.Battle
	tickets           map[uuid.UUID]*Ticket
	ticketEventsChan  chan uuid.UUID
}

func NewBattleGround(log *slog.Logger) *BattleGround {
	return &BattleGround{
		log:               log,
		registeredBattles: make(map[uuid.UUID]*entity.Battle),
		tickets:           make(map[uuid.UUID]*Ticket),
		ticketEventsChan:  make(chan uuid.UUID, 10),
	}
}

// RegisterDuel create a ticket and return ticket id
func (b *BattleGround) RegisterDuel(creator *entity.Hero, bid int, timer time.Duration) uuid.UUID {
	ticket := NewTicket(b.log, creator, timer, DuelType, bid, b.ticketEventsChan)
	b.tickets[ticket.id] = ticket
	go ticket.Start()
	return ticket.id
}

func (b *BattleGround) GetTicket(id uuid.UUID) (*Ticket, error) {
	ticket, ok := b.tickets[id]
	if !ok {
		b.log.Error("Could not find ticket", "id", id)
		return nil, fmt.Errorf("no ticket found with id %v", id)
	}
	return ticket, nil
}

func (b *BattleGround) StartEventLoop(ctx context.Context) {
	b.log.Info("Starting event loop")
	for {
		select {
		case ticketId := <-b.ticketEventsChan:
			switch b.tickets[ticketId].status {
			case ClosedStatus:
				delete(b.tickets, ticketId)
				b.log.Info("ticket closed", "id", ticketId)

			case DoneStatus:
				ticket := b.tickets[ticketId]
				battle := entity.NewBattle(b.log, ticket.players, time.Second*90, ticket.bid)
				b.registeredBattles[battle.GetId()] = battle
				go battle.Start(ctx)
				delete(b.tickets, ticketId)
				b.log.Info("ticket done", "id", ticketId)
			}
		case <-ctx.Done():
			b.log.Info("Stopping event loop")
			return
		}
	}
}
