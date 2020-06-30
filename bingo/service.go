package bingo

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"

	"github.com/divyanshgaba/bingo/game"
	"github.com/divyanshgaba/bingo/ticket"
)

// ErrInvalidArgument is returned when one or more arguments are invalid.
var (
	ErrInvalidArgument     = errors.New("bingo: invalid argument")
	ErrMaxNumbersGenerated = errors.New("bingo: max number of numbers generated for the game")
)

// Service is the interface that provides bingo methods.
type Service interface {
	// create a new bingo game and return its ID
	CreateGame(context.Context) (game.ID, error)

	// create a new ticket for gameID with username
	CreateTicket(context.Context, game.ID, string) (ticket.ID, error)

	// display ticket as HTML
	ShowTicket(context.Context, ticket.ID) (Ticket, error)

	// generate a random number, which has not been picked earlier for this game
	GenerateNumber(context.Context, game.ID) (int64, error)

	// returns all numbers generated so far for a game
	GetAllNumbers(context.Context, game.ID) ([]int64, error)

	// returns stats for a game.
	GetStats(context.Context, game.ID) (numbersDrawn, ticketsGenerated int, err error)
}

type service struct {
	games   game.Repository
	tickets ticket.Repository
	mutex   *sync.Mutex
}

// NewService returns implementation of `Service`.
func NewService(games game.Repository, tickets ticket.Repository) Service {
	s := &service{
		games:   games,
		tickets: tickets,
		mutex:   &sync.Mutex{},
	}
	return s
}

func (s service) CreateGame(ctx context.Context) (game.ID, error) {
	return s.games.Insert(ctx, game.Game{Tickets: []ticket.ID{}, Numbers: []int64{}})
}

func (s service) CreateTicket(ctx context.Context, gameID game.ID, username string) (ticket.ID, error) {
	_, err := s.games.Find(ctx, gameID)
	if err != nil {
		return "", err
	}
	ticket := ticket.Ticket{
		Username: username,
		Cells:    generateRandomCells(),
	}
	ticketID, err := s.tickets.Insert(ctx, ticket)
	if err != nil {
		return "", err
	}
	err = s.games.AddTicket(ctx, gameID, ticketID)
	return ticketID, err
}

func (s service) ShowTicket(ctx context.Context, ticketID ticket.ID) (Ticket, error) {
	t, err := s.tickets.Find(ctx, ticketID)
	if err != nil {
		return Ticket{}, err
	}
	return Ticket{Username: t.Username, Cells: parseCellString(t.Cells)}, nil
}

func (s service) GenerateNumber(ctx context.Context, gameID game.ID) (int64, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	game, err := s.games.Find(ctx, gameID)
	if err != nil {
		return -1, err
	}
	prev := game.Numbers
	if len(prev) == 100 {
		return -1, ErrMaxNumbersGenerated
	}
	set := make(map[int64]bool)
	for _, n := range prev {
		set[n] = true
	}
	perm := rand.Perm(100)
	var number int64
	for _, n := range perm {
		if !set[int64(n)] { // not drawn before
			number = int64(n)
			break
		}
	}
	err = s.games.AddNumber(ctx, gameID, number)
	return number, err
}
func (s service) GetAllNumbers(ctx context.Context, gameID game.ID) ([]int64, error) {
	g, err := s.games.Find(ctx, gameID)
	if err != nil {
		return []int64{}, err
	}
	return g.Numbers, nil
}

func (s service) GetStats(ctx context.Context, gameID game.ID) (numbersDrawn, ticketsGenerated int, err error) {
	g, err := s.games.Find(ctx, gameID)
	if err != nil {
		return 0, 0, err
	}
	return len(g.Numbers), len(g.Tickets), nil
}

// Ticket is a read model for tickets.
type Ticket struct {
	Username string  `json:"username,omitempty"`
	Cells    []int64 `json:"cells,omitempty"`
}

func parseCellString(cells string) []int64 {
	strVals := strings.Split(cells, ";")
	intVals := make([]int64, len(strVals))
	for i, sv := range strVals {
		nv, _ := strconv.ParseInt(sv, 10, 32)
		intVals[i] = nv
	}
	return intVals
}

func generateRandomCells() string {
	perm := rand.Perm(100)[0:15] // ticket will have exactly 15 numbers, in range [0,100)
	var cells []int
	cells = append(cells, createRow(9, perm[:5])...)
	cells = append(cells, createRow(9, perm[5:10])...)
	cells = append(cells, createRow(9, perm[10:])...)
	return arrayToString(cells, ";")
}

func arrayToString(a []int, delim string) string {
	return strings.Trim(strings.Replace(fmt.Sprint(a), " ", delim, -1), "[]")
}

// create row of `size` with `values` randomly populated. Lenght of `values` must be less than or equal to size.
func createRow(size int, values []int) []int {
	if len(values) >= size { // limit
		values = values[:size]
	}
	perm := rand.Perm(size)[:len(values)]
	row := make([]int, size)
	for i := range row {
		row[i] = -1
	}
	for i := range values {
		row[perm[i]] = values[i]
	}
	return row
}
