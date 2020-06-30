package bingo

import (
	"context"

	"github.com/divyanshgaba/bingo/game"
	"github.com/divyanshgaba/bingo/ticket"
	"github.com/go-kit/kit/endpoint"
)

type createGameRequest struct{}

type createGameResponse struct {
	GameID game.ID `json:"game_id,omitempty"`
	Err    error   `json:"error,omitempty"`
}

func (r createGameResponse) error() error { return r.Err }

func makeCreateGameEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		_ = request.(createGameRequest)
		id, err := s.CreateGame(ctx)
		return createGameResponse{GameID: id, Err: err}, nil
	}
}

type createTicketRequest struct {
	GameID   game.ID
	Username string
}
type createTicketResponse struct {
	TicketID ticket.ID `json:"ticket_id,omitempty"`
	Err      error     `json:"error,omitempty"`
}

func (r createTicketResponse) error() error { return r.Err }

func makeCreateTicketEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createTicketRequest)
		id, err := s.CreateTicket(ctx, req.GameID, req.Username)
		return createTicketResponse{TicketID: id, Err: err}, nil
	}
}

type generateNumberRequest struct {
	GameID game.ID
}
type generateNumberResponse struct {
	Number int64 `json:"number,omitempty"`
	Err    error `json:"error,omitempty"`
}

func (r generateNumberResponse) error() error { return r.Err }

func makeGenerateNumberEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(generateNumberRequest)
		num, err := s.GenerateNumber(ctx, req.GameID)
		return generateNumberResponse{Number: num, Err: err}, nil
	}
}

type showTicketRequest struct {
	TicketID ticket.ID
}

type showTicketResponse struct {
	Ticket Ticket `json:"ticket,omitempty"`
	Err    error  `json:"error,omitempty"`
}

func (r showTicketResponse) error() error { return r.Err }

func makeShowTicketEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(showTicketRequest)
		ticket, err := s.ShowTicket(ctx, req.TicketID)
		return showTicketResponse{Ticket: ticket, Err: err}, nil
	}
}

type getAllNumbersRequest struct {
	GameID game.ID
}
type getAllNumbersResponse struct {
	Numbers []int64 `json:"numbers,omitempty"`
	Err     error   `json:"error,omitempty"`
}

func (r getAllNumbersResponse) error() error { return r.Err }

func makeGetAllNumbersEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getAllNumbersRequest)
		numbers, err := s.GetAllNumbers(ctx, req.GameID)
		return getAllNumbersResponse{Numbers: numbers, Err: err}, nil
	}
}

type getStatsRequest struct {
	GameID game.ID
}
type getStatsResponse struct {
	NumbersDrawn     int   `json:"numbers_drawn,omitempty"`
	TicketsGenerated int   `json:"tickets_generated,omitempty"`
	Err              error `json:"error,omitempty"`
}

func (r getStatsResponse) error() error { return r.Err }

func makeGetStatsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getStatsRequest)
		numbersDrawn, ticketsGenerated, err := s.GetStats(ctx, req.GameID)
		return getStatsResponse{NumbersDrawn: numbersDrawn, TicketsGenerated: ticketsGenerated, Err: err}, nil
	}
}
