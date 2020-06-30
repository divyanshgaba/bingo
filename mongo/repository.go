package mongo

import (
	"context"
	"errors"

	"github.com/divyanshgaba/bingo/game"
	"github.com/divyanshgaba/bingo/ticket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type gameRepository struct {
	client *mongo.Client
}

// NewGameRepository returns implementation for game.Repository with mongo store.
func NewGameRepository(client *mongo.Client) game.Repository {
	return &gameRepository{
		client: client,
	}
}

// Game is model for storing games.
type Game struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	Tickets []ticket.ID        `bson:"tickets,omitempty"`
	Numbers []int64            `bson:"numbers,omitempty"`
}

func (r gameRepository) collection() string {
	return "games"
}

func (r gameRepository) Insert(ctx context.Context, g game.Game) (game.ID, error) {
	gs := Game{Tickets: g.Tickets, Numbers: g.Numbers} // game store
	c := r.client.Database(database).Collection(r.collection())
	ior, err := c.InsertOne(ctx, gs)
	if err != nil {
		return "", err
	}
	id, ok := ior.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", errors.New("mongo: error while capturing inserted id from inserted one result")
	}
	return game.ID(id.Hex()), nil
}

func (r gameRepository) Find(ctx context.Context, gameID game.ID) (game.Game, error) {
	var g Game
	c := r.client.Database(database).Collection(r.collection())
	hexID, err := primitive.ObjectIDFromHex(string(gameID))
	if err != nil {
		return game.Game{}, game.ErrInvalidID
	}
	query := bson.M{"_id": hexID}
	sr := c.FindOne(ctx, query)
	err = sr.Decode(&g)
	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return game.Game{}, game.ErrInvalidID
		default:
			return game.Game{}, err
		}
	}
	return game.Game{
		ID:      gameID,
		Tickets: g.Tickets,
		Numbers: g.Numbers,
	}, nil
}
func (r gameRepository) AddTicket(ctx context.Context, gameID game.ID, ticketID ticket.ID) error {
	c := r.client.Database(database).Collection(r.collection())
	gameObjID, err := primitive.ObjectIDFromHex(string(gameID))
	if err != nil {
		return game.ErrInvalidID
	}
	query := bson.M{"_id": gameObjID}
	update := bson.M{"$push": bson.M{"tickets": ticketID}}
	ur, err := c.UpdateOne(ctx, query, update)
	if err != nil {
		return err
	}
	if ur.MatchedCount < 1 {
		return game.ErrInvalidID
	}
	return err
}
func (r gameRepository) AddNumber(ctx context.Context, gameID game.ID, number int64) error {
	c := r.client.Database(database).Collection(r.collection())
	gameObjID, err := primitive.ObjectIDFromHex(string(gameID))
	if err != nil {
		return game.ErrInvalidID
	}
	query := bson.M{"_id": gameObjID}
	update := bson.M{"$push": bson.M{"numbers": number}}
	ur, err := c.UpdateOne(ctx, query, update)
	if err != nil {
		return err
	}
	if ur.MatchedCount < 1 {
		return game.ErrInvalidID
	}
	return err
}

func ticketIDSlice(ids []primitive.ObjectID) []ticket.ID {
	hexIDs := make([]ticket.ID, len(ids))
	for i, id := range ids {
		hexIDs[i] = ticket.ID(id.Hex())
	}
	return hexIDs
}

type ticketRepository struct {
	client *mongo.Client
}

// NewTicketRepository returns implementation for ticket.Repository with mongo store.
func NewTicketRepository(client *mongo.Client) ticket.Repository {
	return &ticketRepository{
		client: client,
	}
}

// Ticket is model for storing tickets.
type Ticket struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Username string             `bson:"username"`
	Cells    string             `bson:"cell"`
}

func (r ticketRepository) collection() string {
	return "tickets"
}

func (r ticketRepository) Insert(ctx context.Context, t ticket.Ticket) (ticket.ID, error) {
	ts := Ticket{Username: t.Username, Cells: t.Cells} // ticket store
	c := r.client.Database(database).Collection(r.collection())
	ior, err := c.InsertOne(ctx, ts)
	if err != nil {
		return "", err
	}
	id, ok := ior.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", errors.New("mongodb: error while capturing inserted id from inserted one result")
	}
	return ticket.ID(id.Hex()), nil
}

func (r ticketRepository) Find(ctx context.Context, ticketID ticket.ID) (ticket.Ticket, error) {
	c := r.client.Database(database).Collection(r.collection())
	hexID, err := primitive.ObjectIDFromHex(string(ticketID))
	if err != nil {
		return ticket.Ticket{}, ticket.ErrInvalidID
	}
	query := bson.M{"_id": hexID}
	sr := c.FindOne(ctx, query)
	var t Ticket
	err = sr.Decode(&t)
	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return ticket.Ticket{}, ticket.ErrInvalidID
		default:
			return ticket.Ticket{}, err
		}
	}
	return ticket.Ticket{
		ID:       ticket.ID(t.ID.Hex()),
		Username: t.Username,
		Cells:    t.Cells,
	}, nil
}
