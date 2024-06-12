package data

const (
	stateValue              = 1
	stateRuestenName        = "ruesten"
	stateRuestenColor       = "orange"
	stateLaufzeitName       = "laufzeit"
	stateLaufzeitColor      = "lightgreen"
	stateProduktivzeitName  = "produktivzeit"
	stateProduktivzeitColor = "darkgreen"
	stateStillstandName     = "stillstand"
	stateStillstandColor    = "grey"
)

func NewState(plc *Plc) *State {
	// plcIdentifier struct to map
	plcIdentifierMap := make(map[string]float64)
	for _, identifier := range plc.Identifier {
		plcIdentifierMap[identifier.Name] = identifier.Value
	}
	var state State
	state.Time = plc.Time
	state.Machine = plc.Machine
	state.Value = stateValue

	if plcIdentifierMap["position_presse_geoeffnet"] == 1.0 {
		// if position_presse_geoeffnet == 1 -> state = "ruesten" -> color = "orange"
		state.Name = stateRuestenName
		state.Color = stateRuestenColor
	} else if plcIdentifierMap["einspritzzeit_ist"] == 0.0 && plcIdentifierMap["heizzeit_ist"] != plcIdentifierMap["heizzeit_soll"] && plcIdentifierMap["heizzeit_ist"] > 0.0 {
		// else if einspritzzeit_ist == 0.0 AND heizzeit_ist != heizzeit_soll AND heizzeit_ist > 0.0 -> state = "laufzeit" -> color = "lightgreen"
		state.Name = stateLaufzeitName
		state.Color = stateLaufzeitColor
	} else if plcIdentifierMap["einspritzzeit_ist"] != plcIdentifierMap["einspritzzeit_soll"] && plcIdentifierMap["heizzeit_ist"] > 0.0 {
		// else if einspritzzeit_ist != einspritzzeit_soll AND heizzeit_ist > 0.0 -> state = "produktivzeit" -> color = "darkgreen"
		state.Name = stateProduktivzeitName
		state.Color = stateProduktivzeitColor
	} else {
		// else -> state = "stillstand" -> color = "grey"
		state.Name = stateStillstandName
		state.Color = stateStillstandColor
	}

	return &state
}
