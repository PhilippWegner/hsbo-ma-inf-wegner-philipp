package main

import (
	"fmt"

	"github.com/PhilippWegner/hsbo-ma-inf-wegner-philipp/microservice-architecture/statemachine-service/data"
)

var DEFAULT_LAST_STATE = data.State{Time: "0"}

func (c *Config) calculate(machine string) error {
	// get last states
	lastState := DEFAULT_LAST_STATE
	state, err := c.ApiRepository.GetStates(machine)
	if err != nil {
		_ = c.ApiRepository.CreateLog(data.Log{Name: "Error", Data: err.Error()})
	}
	if len(state) > 0 {
		lastState = *state[0]
	}
	_ = c.ApiRepository.CreateLog(data.Log{Name: "GetState", Data: fmt.Sprintf("%s %s %s %s %d", lastState.Time, lastState.Machine, lastState.Name, lastState.Color, lastState.Value)})
	plcs, err := c.ApiRepository.GetPlcs(machine, lastState)
	if err != nil {
		_ = c.ApiRepository.CreateLog(data.Log{Name: "error", Data: err.Error()})
		return err
	}
	_ = c.ApiRepository.CreateLog(data.Log{Name: "GetPlcs", Data: fmt.Sprintf("Plcs received: %d", len(plcs))})
	var states []*data.State
	for _, plc := range plcs {
		state := data.NewState(plc)
		states = append(states, state)
	}
	err = c.ApiRepository.CreateState(states)
	if err != nil {
		_ = c.ApiRepository.CreateLog(data.Log{Name: "Error", Data: err.Error()})
		return err
	}
	_ = c.ApiRepository.CreateLog(data.Log{Name: "CreateStates", Data: fmt.Sprintf("States created: %d", len(states))})
	return nil
}
