package bingo

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/divyanshgaba/bingo/game"
	"github.com/divyanshgaba/bingo/ticket"
	"github.com/gorilla/mux"

	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
)

var errBadRoute = errors.New("bingo: bad route")

// MakeHandler returns a handler for the bingo service.
func MakeHandler(bs Service, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encodeError),
	}
	createGameHandler := kithttp.NewServer(
		makeCreateGameEndpoint(bs),
		decodeCreateGameRequest,
		encodeResponse,
		opts...,
	)
	createTicketHandler := kithttp.NewServer(
		makeCreateTicketEndpoint(bs),
		decodeCreateTicketRequest,
		encodeResponse,
		opts...,
	)

	generateNumberHandler := kithttp.NewServer(
		makeGenerateNumberEndpoint(bs),
		decodeGenerateNumberRequest,
		encodeResponse,
		opts...,
	)

	showTicketHandler := kithttp.NewServer(
		makeShowTicketEndpoint(bs),
		decodeShowTicketRequest,
		encodeShowTicketResponse,
		opts...,
	)

	getAllNumbersHandler := kithttp.NewServer(
		makeGetAllNumbersEndpoint(bs),
		decodeGetAllNumbersRequest,
		encodeResponse,
		opts...,
	)

	getStatsHandler := kithttp.NewServer(
		makeGetStatsEndpoint(bs),
		decodeGetStatsRequest,
		encodeResponse,
		opts...,
	)

	r := mux.NewRouter()

	r.Handle("/api/game/create", createGameHandler).Methods(http.MethodPost)
	r.Handle("/api/game/{gameId}/ticket/{username}/generate", createTicketHandler).Methods(http.MethodPost)
	r.Handle("/api/game/{gameId}/number/random", generateNumberHandler).Methods(http.MethodGet)
	r.Handle("/ticket/{ticketId}", showTicketHandler).Methods(http.MethodGet)
	r.Handle("/api/game/{gameId}/numbers", getAllNumbersHandler).Methods(http.MethodGet)
	r.Handle("/api/game/{gameId}/stats", getStatsHandler).Methods(http.MethodGet)

	return r
}

func decodeCreateGameRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return createGameRequest{}, nil
}

func decodeCreateTicketRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	gameID, ok := vars["gameId"]
	if !ok {
		return nil, errBadRoute
	}
	username, ok := vars["username"]
	if !ok {
		return nil, errBadRoute
	}
	return createTicketRequest{GameID: game.ID(gameID), Username: username}, nil
}

func decodeGenerateNumberRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	gameID, ok := vars["gameId"]
	if !ok {
		return nil, errBadRoute
	}
	return generateNumberRequest{GameID: game.ID(gameID)}, nil
}

func decodeGetAllNumbersRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	gameID, ok := vars["gameId"]
	if !ok {
		return nil, errBadRoute
	}
	return getAllNumbersRequest{GameID: game.ID(gameID)}, nil
}

func decodeGetStatsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	gameID, ok := vars["gameId"]
	if !ok {
		return nil, errBadRoute
	}
	return getStatsRequest{GameID: game.ID(gameID)}, nil
}

func decodeShowTicketRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	ticketID, ok := vars["ticketId"]
	if !ok {
		return nil, errBadRoute
	}
	return showTicketRequest{TicketID: ticket.ID(ticketID)}, nil
}

func encodeShowTicketResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	e, ok := response.(showTicketResponse)
	if ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}
	ticket := e.Ticket
	cells := ticket.Cells
	table := "<table style='table-layout:fixed' border='1px solid black';>" + encodeTableRow(cells[:9]) + encodeTableRow(cells[9:18]) + encodeTableRow(cells[18:]) + "</table>"

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(table))
	return nil
}

func encodeTableRow(row []int64) string {
	tds := ""
	for _, r := range row {
		if r != -1 {
			tds = tds + "<td style='width: 11%;'>" + strconv.FormatInt(r, 10) + "</td>"
		} else {
			tds = tds + "<td style='width: 11%;'></td>"
		}
	}
	return "<tr>" + tds + "</tr>"
}

// encode errors from business-logic
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	case game.ErrInvalidID:
		w.WriteHeader(http.StatusBadRequest)
	case ticket.ErrInvalidID:
		w.WriteHeader(http.StatusBadRequest)
	case ErrMaxNumbersGenerated:
		w.WriteHeader(http.StatusBadRequest)
	case ErrInvalidArgument:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

type errorer interface {
	error() error
}
