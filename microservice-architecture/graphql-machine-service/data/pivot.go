package data

import (
	"encoding/json"
	"sort"

	"github.com/PhilippWegner/hsbo-ma-inf-wegner-philipp/microservice-architecture/graphql-machine-service/graph/model"
)

func Pivotize(influxPlcs [][]interface{}) []*model.Plc {
	var uniquetimes []string

	for _, influxPlc := range influxPlcs {
		time := influxPlc[0].(json.Number).String()
		if !contains(uniquetimes, time) {
			uniquetimes = append(uniquetimes, time)
		}
	}
	// order uniquetimes ascending
	sort.Strings(uniquetimes)

	var plcs []*model.Plc
	for _, uniquetime := range uniquetimes {
		var plc model.Plc
		plc.Time = uniquetime
		var identifiers []*model.Identifier
		for _, influxPlc := range influxPlcs {
			time := influxPlc[0].(json.Number).String()
			if time == uniquetime {
				machine, _ := influxPlc[1].(string)
				bezeichner, _ := influxPlc[2].(string)
				value, _ := influxPlc[3].(json.Number).Float64()
				identifier := model.Identifier{Name: bezeichner, Value: value}
				identifiers = append(identifiers, &identifier)
				plc.Machine = machine
				plc.Identifier = identifiers
			}
		}
		plcs = append(plcs, &plc)
	}
	// cut of the last plc, because it might not be a full plc
	plcs = plcs[:len(plcs)-1]
	return plcs
}

func contains(uniquetimes []string, time string) bool {
	for _, uniquetime := range uniquetimes {
		if uniquetime == time {
			return true
		}
	}
	return false
}
