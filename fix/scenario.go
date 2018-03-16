package fix

type Scenario struct {
	Name   string `toml:"name"`
	Tables Tables `toml:"table"`
}

// type Scenarios []Scenario
type Scenarios struct {
	Scenarios []Scenario `toml:"scenario"`
}

type Table struct {
	Name string                   `toml:"name"`
	Data []map[string]interface{} `toml:"data"`
}

type Tables []Table
