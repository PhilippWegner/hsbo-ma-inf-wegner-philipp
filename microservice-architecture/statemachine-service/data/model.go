package data

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

type Plc struct {
	Time       string `json:"time"`
	Machine    string `json:"machine"`
	Identifier []*Identifier
}

type Identifier struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

type State struct {
	Time    string `json:"time"`
	Machine string `json:"machine"`
	Name    string `json:"name"`
	Color   string `json:"color"`
	Value   int    `json:"value"`
}

type RequestPayload struct {
	Action    string   `json:"action"`
	Machine   string   `json:"machine,omitempty"`
	LastState State    `json:"lastState,omitempty"`
	Plcs      []*Plc   `json:"plcs,omitempty"`
	States    []*State `json:"states,omitempty"`
	// Log       Log      `json:"log,omitempty"`
}

type ResponsePayload struct {
	Plcs   []*Plc   `json:"plcs,omitempty"`
	States []*State `json:"states,omitempty"`
	Error  string   `json:"error,omitempty"`
}

type ApiRepository struct {
	restful_api string
}

func NewApiRepository(restful_api string) *ApiRepository {
	return &ApiRepository{restful_api: restful_api}
}

func (r *ApiRepository) GetPlcs(machine string, laststate State) ([]*Plc, error) {
	log.Println("GetPlcs")
	RequestPayload := RequestPayload{
		Action:    "next-plcs",
		Machine:   machine,
		LastState: laststate,
	}
	requestPayloadJson, err := json.Marshal(RequestPayload)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest("POST", r.restful_api, bytes.NewBuffer(requestPayloadJson))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	var responsePayload *ResponsePayload
	err = json.NewDecoder(response.Body).Decode(&responsePayload)
	if err != nil {
		return nil, err
	}
	// log len of plcs
	log.Printf("len(plcs): %v\n", len(responsePayload.Plcs))
	return responsePayload.Plcs, nil

}

func (r *ApiRepository) GetStates(machine string) ([]*State, error) {
	log.Println("GetStates")
	RequestPayload := RequestPayload{
		Action:  "last-state",
		Machine: machine,
	}
	requestPayloadJson, err := json.Marshal(RequestPayload)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest("POST", r.restful_api, bytes.NewBuffer(requestPayloadJson))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	var responsePayload *ResponsePayload
	err = json.NewDecoder(response.Body).Decode(&responsePayload)
	if err != nil {
		return nil, err
	}
	// log len of states
	log.Printf("len(states): %v\n", len(responsePayload.States))
	return responsePayload.States, nil
}

func (r *ApiRepository) CreateState(states []*State) error {
	log.Println("CreateState")
	RequestPayload := RequestPayload{
		Action: "insert-states",
		States: states,
	}
	requestPayloadJson, err := json.Marshal(RequestPayload)
	if err != nil {
		return err
	}
	request, err := http.NewRequest("POST", r.restful_api, bytes.NewBuffer(requestPayloadJson))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	return nil
}
