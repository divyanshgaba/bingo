package game

import (
	"context"
	"errors"

	"github.com/divyanshgaba/bingo/ticket"
)

// errors for package game
var (
	ErrInvalidID = errors.New("game: invalid ID")
)

// ID uniquely identifies a particular game.
type ID string

// Game is the central class in the domain model.
type Game struct {
	ID      ID
	Tickets []ticket.ID
	Numbers []int64
}

// New creates a new game.
func New(id ID, tickets []ticket.ID, numbers []int64) *Game {
	return &Game{
		ID:      id,
		Tickets: tickets,
		Numbers: numbers,
	}
}

// Repository provides access a game store.
type Repository interface {
	Insert(context.Context, Game) (ID, error)
	Find(context.Context, ID) (Game, error)
	AddTicket(context.Context, ID, ticket.ID) error
	AddNumber(context.Context, ID, int64) error
}
