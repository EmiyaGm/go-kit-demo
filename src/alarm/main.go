package main

import (
	"time"
	"alarm/client"
)

func main() {
	client.Init()

	go client.Run()
	// Set up a connection to the server.
	client.Alarm() <- &client.Message{
		ID:       "5h9gE3VVIu",
		FlowID:   123,
		Source:   "jtt808_position",
		Type:     "telematicsbox_powerfail",
		Time:     time.Now(),
		Strategy: "create",
		Target:   "EnYtWhHl2g",
		SourceID: "jtt808_position",
		Location: []float64{
			106.354564138953,
			26.6664378969494,
		},
	}

	client.Alarm() <- &client.Message{
		ID:       "5h9gE3VVIu",
		Source:   "jtt808_position",
		Type:     "telematicsbox_powerfail",
		Time:     time.Now(),
		Strategy: "add",

		Target: "O6zwIIRi70",
		Location: []float64{
			106.354564138953,
			26.6664378969494,
		},
	}

	client.Alarm() <- &client.Message{
		ID:       "5h9gE3VVIu",
		Source:   "jtt808_position",
		Type:     "telematicsbox_powerfail",
		Time:     time.Now(),
		Strategy: "end",

		Target: "f87scL4oLu",
		Location: []float64{
			106.354564138953,
			26.6664378969494,
		},
	}
	time.Sleep(1 * time.Second)
}
