package tests

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/divyanshgaba/bingo/ticket"
)

func TestCreateGame(t *testing.T) {
	req, _ := http.NewRequest(http.MethodPost, srvURL+"/api/game/create", nil)
	resp, _ := http.DefaultClient.Do(req)
	if resp.StatusCode != 200 {
		t.Errorf("invalid status code: want=200 got=%v", resp.StatusCode)
	}
	game := struct {
		GameID string `json:"game_id,omitempty"`
	}{}
	jd := json.NewDecoder(resp.Body)
	err := jd.Decode(&game)
	if err != nil {
		t.Errorf("error while decoding response body err=%v", err)
	}
	if len(game.GameID) == 0 {
		t.Errorf("Empty game_id in response body")
	}
}

func TestCreateTicket(t *testing.T) {
	req, _ := http.NewRequest(http.MethodPost, srvURL+"/api/game/"+string(gameID)+"/ticket/testUser/generate", nil)
	resp, _ := http.DefaultClient.Do(req)
	if resp.StatusCode != 200 {
		t.Errorf("invalid status code: want=200 got=%v", resp.StatusCode)
	}
	ticketResp := struct {
		TicketID string `json:"ticket_id,omitempty"`
	}{}
	jd := json.NewDecoder(resp.Body)
	err := jd.Decode(&ticketResp)
	if err != nil {
		t.Errorf("error while decoding response body err=%v", err)
	}
	if len(ticketResp.TicketID) == 0 {
		t.Errorf("Empty game_id in response body")
	}
	ticketStore, err := tickets.Find(context.Background(), ticket.ID(ticketResp.TicketID))
	if err != nil {
		t.Errorf("error while retrieving ticket ID=%v err=%v", ticketResp.TicketID, err)
	}
	if ticketStore.Username != "testUser" {
		t.Errorf("error invalid username want=testUser got=%v", ticketStore.Username)
	}
}
func TestGenerateNumber(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, srvURL+"/api/game/"+string(gameID)+"/number/random", nil)
	resp, _ := http.DefaultClient.Do(req)
	if resp.StatusCode != 200 {
		t.Errorf("invalid status code: want=200 got=%v", resp.StatusCode)
	}
	number := struct {
		Number *int64 `json:"number,omitempty"`
	}{}
	jd := json.NewDecoder(resp.Body)
	err := jd.Decode(&number)
	if err != nil {
		t.Errorf("error while decoding response body err=%v", err)
	}
	if number.Number == nil {
		t.Errorf("number not set in response body")
	}
}

func TestGetAllNumbers(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, srvURL+"/api/game/"+string(gameID)+"/numbers", nil)
	resp, _ := http.DefaultClient.Do(req)
	if resp.StatusCode != 200 {
		t.Errorf("invalid status code: want=200 got=%v", resp.StatusCode)
	}
	number := struct {
		Numbers []int64 `json:"numbers,omitempty"`
	}{}
	jd := json.NewDecoder(resp.Body)
	err := jd.Decode(&number)
	if err != nil {
		t.Errorf("error while decoding response body err=%v", err)
	}
	if number.Numbers == nil {
		t.Errorf("number not set in response body")
	}
}
func TestGetStats(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, srvURL+"/api/game/"+string(gameID)+"/stats", nil)
	resp, _ := http.DefaultClient.Do(req)
	if resp.StatusCode != 200 {
		t.Errorf("invalid status code: want=200 got=%v", resp.StatusCode)
	}
	stats := struct {
		NumbersDrawn     *int64 `json:"numbers_drawn,omitempty"`
		TicketsGenerated *int64 `json:"tickets_generated,omitempty"`
	}{}
	jd := json.NewDecoder(resp.Body)
	err := jd.Decode(&stats)
	if err != nil {
		t.Errorf("error while decoding response body err=%v", err)
	}
	if stats.NumbersDrawn == nil {
		t.Errorf("number not set in response body")
	}
	if stats.TicketsGenerated == nil {
		t.Errorf("number not set in response body")
	}
}
