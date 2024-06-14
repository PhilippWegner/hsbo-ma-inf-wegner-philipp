package database

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/PhilippWegner/hsbo-ma-inf-wegner-philipp/microservice-architecture/graphql-machine-service/data"
	"github.com/PhilippWegner/hsbo-ma-inf-wegner-philipp/microservice-architecture/graphql-machine-service/graph/model"
	"github.com/influxdata/influxdb/client/v2"
)

var (
	dbHost = "influxdb"
	dbPort = "8086"
	dbUser = "root"
	dbPass = "root"
	dbName = "reich"
)

type Influxdb struct {
	client client.HTTPClient
}

func ConnectInfluxdb() *Influxdb {
	client, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://" + dbHost + ":" + dbPort,
		Username: dbUser,
		Password: dbPass,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	return &Influxdb{client: client}
}

func (db *Influxdb) CreateStates(createStatesInput []*model.CreateStatesInput) ([]*model.State, error) {
	// log.Println("CreateStates")
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  dbName,
		Precision: "ns",
	})
	if err != nil {
		return nil, err
	}
	for _, createStateInput := range createStatesInput {
		tags := map[string]string{
			"machine": createStateInput.Machine,
			"state":   createStateInput.Name,
			"color":   createStateInput.Color,
		}
		fields := map[string]interface{}{
			"value": createStateInput.Value,
		}
		time_int64, _ := strconv.ParseInt(createStateInput.Time, 10, 64)
		pt, err := client.NewPoint("statemachine", tags, fields, time.Unix(0, time_int64))
		if err != nil {
			return nil, err
		}
		bp.AddPoint(pt)
	}
	err = db.client.Write(bp)
	if err != nil {
		return nil, err
	}
	var states []*model.State
	for _, createStateInput := range createStatesInput {
		state := model.State{
			Time:    createStateInput.Time,
			Machine: createStateInput.Machine,
			Name:    createStateInput.Name,
			Color:   createStateInput.Color,
			Value:   createStateInput.Value,
		}
		states = append(states, &state)
	}
	return states, nil
}

func (db *Influxdb) GetStates(machine string, limit int) ([]*model.State, error) {
	// log.Println("GetStates")
	stmt := fmt.Sprintf("SELECT time, machine, state, color, value from statemachine WHERE machine = '%s' ORDER BY time DESC LIMIT %d", machine, limit)
	// log.Println(stmt)
	response, err := db.executeQuery(stmt)
	if err != nil {
		return nil, err
	}
	if len(response.Results[0].Series) == 0 {
		return nil, fmt.Errorf("no data found")
	}
	values := response.Results[0].Series[0].Values
	var states []*model.State
	for _, value := range values {
		time := value[0].(json.Number).String()
		machine, _ := value[1].(string)
		name, _ := value[2].(string)
		color, _ := value[3].(string)
		value_int64, _ := value[4].(json.Number).Int64()
		value := int(value_int64)
		state := model.State{
			Time:    time,
			Machine: machine,
			Name:    name,
			Color:   color,
			Value:   value,
		}
		states = append(states, &state)
	}
	return states, nil
}

func (db *Influxdb) GetPlcs(machine string, time string, limit int, filter *model.IdentifierFilterInput) ([]*model.Plc, error) {
	// log.Println("GetPlcs")
	stmt := ""
	if filter == nil || filter.Identifier == nil {
		// log.Println("filter is nil")
		stmt = fmt.Sprintf("SELECT time, maschine, bezeichner, value from data WHERE maschine = '%s' AND time >= %s ORDER BY time ASC LIMIT %d", machine, time, limit)
	} else if len(filter.Identifier.In) > 0 {
		// log.Println("filter.Identifier.In is filled")
		// create empty list
		var identifiers_in []string
		// for all elements in the filter.Identifier.In array create a query
		for _, id := range filter.Identifier.In {
			identifier_in := fmt.Sprintf("bezeichner = '%s'", *id)
			identifiers_in = append(identifiers_in, identifier_in)
		}
		identifers_in_join := strings.Join(identifiers_in, " OR ")
		stmt = fmt.Sprintf("SELECT time, maschine, bezeichner, value from data WHERE maschine = '%s' AND time >= %s AND (%s) ORDER BY time ASC LIMIT %d", machine, time, identifers_in_join, limit)
	}
	// log.Println(stmt)
	// execute query
	response, err := db.executeQuery(stmt)
	if err != nil {
		return nil, err
	}
	// get values
	values := response.Results[0].Series[0].Values
	// log.Println(values)
	// pivotize
	plcs := data.Pivotize(values)
	return plcs, nil
}

func (db *Influxdb) executeQuery(stmt string) (*client.Response, error) {
	query := client.NewQuery(stmt, dbName, "ns")
	response, err := db.client.Query(query)
	if err != nil {
		return nil, err
	}
	if response.Error() != nil {
		return nil, response.Error()
	}
	return response, nil
}
