package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hasura/go-graphql-client"
)

const DEFAULT_LIMIT = 1000

const DEFAULT_GRAPHQL_MACHINE_URL = "http://localhost:8081/query"

type Request struct {
	Action    string               `json:"action"`
	Machine   string               `json:"machine,omitempty"`
	LastState State                `json:"lastState,omitempty"`
	States    []*CreateStatesInput `json:"states,omitempty"`
	Log       Log                  `json:"log,omitempty"`
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
// @Success 200 {object} object
// @Failure 400 {object} string
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
	default:
		log.Println("Invalid action")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid action"})
		return
	}
}

func (app *Config) lastState(ctx *gin.Context, machine string) {
	// graphql client
	client := graphql.NewClient(DEFAULT_GRAPHQL_MACHINE_URL, nil)
	var query struct {
		States []*State `graphql:"createStates(machine: $machine, limit: $limit)"`
	}
	variables := map[string]interface{}{
		"machine": machine,
		"limit":   1,
	}
	err := client.Query(context.Background(), &query, variables)
	if err != nil {
		log.Println("error:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, query.States)
}

func (app *Config) nextPlcs(ctx *gin.Context, machine string, state State) {
	// graphql client
	client := graphql.NewClient(DEFAULT_GRAPHQL_MACHINE_URL, nil)
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
		log.Println("error:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, query.Plcs)
}

func (app *Config) insertStates(ctx *gin.Context, statesInput []*CreateStatesInput) {
	// graphql client
	client := graphql.NewClient(DEFAULT_GRAPHQL_MACHINE_URL, nil)
	var mutation struct {
		CreateStates []*CreateStatesInput `graphql:"createStates(input: $input)"`
	}
	variables := map[string]interface{}{
		"input": statesInput,
	}
	err := client.Mutate(context.Background(), &mutation, variables)
	if err != nil {
		log.Println("error:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, mutation.CreateStates)
}
