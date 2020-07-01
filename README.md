# Bingo
![bingo](docs/assets/bingo.jpeg?raw=true)

## Deploy
```
docker-compose up --build  
```
## Test cases
```
TODO: Add automated test cases
```
## Code structure
This project uses [go-kit](https://github.com/go-kit/kit) framework for microservices and [MongoDB](https://www.mongodb.com/) as data store.
1. Package `config` and `config.yaml`, holds data store credentials. 
2. Package `bingo`,
    1. `service.go` holds all the buisiness logic.
    2. `endpoint.go` holds request and response schema for each endpoint. 
    3. `transport.go` holds details how each request is decoded and corresponding response is encoded.
3. Package `game` holds domain model and repository interface for games.
4. Package `ticket` holds domain model and repository interface for tickets.
5. Package `mongo` holds implementation of `game` and `ticket` repository with mongo store.

## Models
1. `Game` has,
    1. Unique string ID.
    2. List of IDs of tickets generated.
    3. List of Numbers drawn
2. `Ticket` has,
    1. Unique string ID.
    2. Username for which this ticket was generated.
    3. Cell values represented as semicolon(;) seperated integer values. `-1` signifies empty cell.


## APIs
For manual testing: [Postman collection](docs/assets/Bingo.postman_collection.json?raw=true)


```
1. POST /api/game/create
Success HTTP 200
{
    "game_id": "5efbbb5c0880edb67dd9fd33"
}
```
```
2. POST /api/game/{game_id}/ticket/{username}/generate
Success HTTP 200
{
    "ticket_id": "5efbbf88abaf8fdbd4aae78d"
}
```
```
3. GET /api/game/{game_id}/number/random
Success HTTP 200
{
    "number": 79
}
```
```
4. GET /api/game/{game_id}/numbers
Success HTTP 200
{
    "numbers": [
        36,
        79
    ]
}
```
```
5. GET /api/game/{game_id}/stats
Success HTTP 200
{
    "numbers_drawn": 2,
    "tickets_generated": 1
}
```
```
6. GET /ticket/{ticket_id}
Success HTTP 200
```
![ticket](docs/assets/ticket.png?raw=true)
