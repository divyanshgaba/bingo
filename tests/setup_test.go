package tests

import (
	"context"
	"fmt"
	"os"
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/divyanshgaba/bingo/bingo"
	"github.com/divyanshgaba/bingo/config"
	"github.com/divyanshgaba/bingo/game"
	"github.com/divyanshgaba/bingo/mongo"
	"github.com/divyanshgaba/bingo/ticket"
	"github.com/go-kit/kit/log"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	srvURL   string
	gameID   game.ID
	ticketID ticket.ID
	games    game.Repository
	tickets  ticket.Repository
)

func TestMain(m *testing.M) {
	if config.Env() != "test" {
		fmt.Println("pass -env=test to enable testing")
		os.Exit(1)
	}
	// setup
	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	mongoClient, err := mongo.NewClient(logger)
	if err != nil {
		panic("could not create mongo client")
	}
	games = mongo.NewGameRepository(mongoClient)
	tickets = mongo.NewTicketRepository(mongoClient)
	var bs bingo.Service
	bs = bingo.NewService(games, tickets)
	httpLogger := log.With(logger, "component", "http-test")

	mux := http.NewServeMux()
	mux.Handle("/", bingo.MakeHandler(bs, httpLogger))
	// run server
	srv := httptest.NewServer(mux)
	defer srv.Close()
	// setup database
	Database()
	srvURL = srv.URL
	os.Exit(m.Run())
}

func Database() {
	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	mongoclient, err := mongo.NewClient(logger)
	if err != nil {
		panic(err)
	}
	db := mongoclient.Database(config.Mongo().Database)
	// clean db
	db.Collection("games").DeleteMany(context.Background(), bson.D{})
	db.Collection("tickets").DeleteMany(context.Background(), bson.D{})
	gameID, _ = games.Insert(context.Background(), game.Game{})
}
