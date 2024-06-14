package main

import (
	"context"
	"log"
	"net/http"

	"github.com/PhilippWegner/hsbo-ma-inf-wegner-philipp/microservice-architecture/restful-broker-service/model"
	"github.com/gin-gonic/gin"
	"github.com/hasura/go-graphql-client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const DEFAULT_LIMIT = 1000

const DEFAULT_GRAPHQL_URL = "graphql-api:80/query"
const DEFAULT_GRPC_URL = "grpc-api:80"

type Request struct {
	Action    string               `json:"action"`
	Machine   string               `json:"machine,omitempty"`
	LastState State                `json:"lastState,omitempty"`
	States    []*CreateStatesInput `json:"states,omitempty"`
	Log       Log                  `json:"log,omitempty"`
}

type Response struct {
	Plcs   []*Plc   `json:"plcs,omitempty"`
	States []*State `json:"states,omitempty"`
	Error  string   `json:"error,omitempty"`
}

type State struct {
	Time    string `json:"time"`
	Machine string `json:"machine"`
	Name    string `json:"name"`
	Color   string `json:"color"`
	Value   int64  `json:"value"`
}

type CreateStatesInput struct {
	Time    string `json:"time"`
	Machine string `json:"machine"`
	Name    string `json:"name"`
	Color   string `json:"color"`
	Value   int64  `json:"value"`
}

type Plc struct {
	Time       string        `json:"time"`
	Machine    string        `json:"machine"`
	Identifier []*Identifier `json:"identifier"`
}

type Identifier struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

type Log struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

// Handle godoc
// @Summary Handle request
// @Description Handle request
// @Tags api
// @Accept  application/json
// @Produce  application/json
// @Param request body Request true "Request"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /handle [post]
func (app *Config) Handle(ctx *gin.Context) {
	log.Println("Handle request")
	var request Request
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	switch request.Action {
	case "last-state":
		log.Println("last-state")
		app.lastState(ctx, request.Machine)
	case "next-plcs":
		log.Println("next-plcs")
		app.nextPlcs(ctx, request.Machine, request.LastState)
	case "insert-states":
		log.Println("insert-states")
		app.insertStates(ctx, request.States)
	case "log":
		log.Println("log")
		app.insertLog(ctx, request.Log)
	default:
		log.Println("Invalid action")
		response := Response{Error: "Invalid action"}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
}

func (app *Config) lastState(ctx *gin.Context, machine string) {
	// graphql client
	client := graphql.NewClient("http://"+DEFAULT_GRAPHQL_URL, nil)
	var query struct {
		States []*State `graphql:"states(machine: $machine, limit: $limit)"`
	}
	variables := map[string]interface{}{
		"machine": machine,
		"limit":   1,
	}
	err := client.Query(context.Background(), &query, variables)
	if err != nil {
		response := Response{Error: err.Error()}
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}
	response := Response{States: query.States}
	ctx.JSON(http.StatusOK, response)
}

func (app *Config) nextPlcs(ctx *gin.Context, machine string, state State) {
	// graphql client
	client := graphql.NewClient("http://"+DEFAULT_GRAPHQL_URL, nil)
	var query struct {
		Plcs []*Plc `graphql:"plcs(machine: $machine, time: $time, limit: $limit, filter: {identifier: {in: $in}})"`
	}
	variables := map[string]interface{}{
		"machine": machine,
		"time":    state.Time,
		"limit":   DEFAULT_LIMIT,
		"in":      []string{"heizzeit_ist", "heizzeit_soll", "einspritzzeit_ist", "einspritzzeit_soll", "position_presse_geoeffnet"},
	}
	// log.Println("variables:", variables)
	err := client.Query(context.Background(), &query, variables)
	if err != nil {
		// log.Println("error:", err)
		response := Response{Error: err.Error()}
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}
	response := Response{Plcs: query.Plcs}
	ctx.JSON(http.StatusOK, response)
}

func (app *Config) insertStates(ctx *gin.Context, statesInput []*CreateStatesInput) {
	// graphql client
	client := graphql.NewClient("http://"+DEFAULT_GRAPHQL_URL, nil)
	var mutation struct {
		CreateStates []*State `graphql:"createStates(input: $input)"`
	}
	variables := map[string]interface{}{
		"input": statesInput,
	}
	err := client.Mutate(context.Background(), &mutation, variables)
	if err != nil {
		response := Response{Error: err.Error()}
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}
	response := Response{States: mutation.CreateStates}
	ctx.JSON(http.StatusOK, response)
}

func (app *Config) insertLog(ctx *gin.Context, logEntry Log) {
	// grpc client
	conn, err := grpc.NewClient(DEFAULT_GRPC_URL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		response := Response{Error: err.Error()}
		log.Println("Error:", err)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}
	defer conn.Close()
	client := model.NewLogServiceClient(conn)
	client.WriteLog(context.Background(), &model.LogRequest{LogEntry: &model.Log{Name: logEntry.Name, Data: logEntry.Data}})
	log.Println("Log written")
}
