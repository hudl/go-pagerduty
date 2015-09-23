package main

import (
	"fmt"

	"github.com/hudl/go-pagerduty/pagerduty"
)

func main() {
	client := pagerduty.NewClient(nil, "subdomain", "super-secret-api-key")

	opts := &pagerduty.IncidentListOptions{
		Status: pagerduty.StatusTriggered,
	}

	// get a list of triggered incidents
	incidents, _, err := client.Incidents.List(opts)
	if err != nil {
		panic(err)
	}

	fmt.Println(incidents)
}
