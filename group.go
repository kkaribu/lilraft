package lilraft

import "fmt"

// Group ...
type Group struct {
	Logs   map[string]*Log
	Config Config

	inbox chan message
}

// Run ...
func (g *Group) Run() {
	g.Logs = map[string]*Log{
		"abc123": {
			id: "abc123",
		},
	}

	go func() {}()

	for id := range g.Logs {
		fmt.Printf("State of %s: %d\n", g.Logs[id].id, g.Logs[id].state)
	}
}

func (g *Group) send(to string, msg message) error {
	return nil
}
