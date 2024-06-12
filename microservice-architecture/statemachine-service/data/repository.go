package data

type Repository interface {
	GetStates(machine string) ([]*State, error)
	GetPlcs(machine string, state State) ([]*Plc, error)
	CreateState(state []*State) error
}
