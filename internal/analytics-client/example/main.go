package main

import (
	"dev.khulnasoft.com/analytics-client/client"
)

func main() {
	analyticsClient := client.NewClient()

	// record a new event
	analyticsClient.RecordEvent(client.Event{
		"my-table": map[string]interface{}{
			"field": "my-field",
		},
	})

	// flush the events
	analyticsClient.Flush()
}
