package main

import (
	"log"

	"github.com/PhilippWegner/hsbo-ma-inf-wegner-philipp/microservice-architecture/statemachine-service/data"
)

var DEFAULT_LAST_STATE = data.State{Time: "0"}

func (c *Config) calculate(machine string) error {
	// get last states
	lastState := DEFAULT_LAST_STATE
	state, err := c.ApiRepository.GetStates(machine)
	// log.Println("calculate:", state)
	if err != nil {
		log.Println("calculate err:", err)
	}
	if len(state) > 0 {
		lastState = *state[0]
	}
	// log.Println("lastState:", lastState)
	plcs, err := c.ApiRepository.GetPlcs(machine, lastState)
	if err != nil {
		return err
	}
	var states []*data.State
	for _, plc := range plcs {
		state := data.NewState(plc)
		states = append(states, state)
	}
	log.Printf("State created: %d\n", len(states))
	err = c.ApiRepository.CreateState(states)
	if err != nil {
		return err
	}
	return nil
}
