package ticket

import (
	"context"
	"errors"
)

// errors for package ticket
var (
	ErrInvalidID = errors.New("ticket: invalid ID")
)

// ID uniquely identifies a particular ticket.
type ID string

// Ticket is the central class in the domain model.
// Each ticket is represented in string with each cell seperated by a semicolon(;).
type Ticket struct {
	ID       ID
	Username string
	Cells    string
}

// New creates a new ticket
func New(id ID, username, cells string) *Ticket {
	return &Ticket{
		ID:       id,
		Username: username,
		Cells:    cells,
	}
}

// Repository provides access a game store.
type Repository interface {
	Insert(context.Context, Ticket) (ID, error)
	Find(context.Context, ID) (Ticket, error)
}
