package main

import (
	"fmt"
	"github.com/brettniven/embedqueries"
	"github.com/sirupsen/logrus"
)

func main() {

	// Construct a client. This will load the embedded queries and parse them into Go Templates, read for variable substitution
	c, err := embedqueries.NewClient()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to construct client")
	}

	// invoke some API calls
	res, err := c.MissionByID("6C42550")
	handleResponse("Mission 6C42550", res, err)

	res, err = c.MissionsByManufacturer("Orbital ATK", 10)
	handleResponse("Missions of Orbital ATK", res, err)

	res, err = c.PastLaunches(3)
	handleResponse("Last 3 Launches", res, err)

	res, err = c.Rockets(10)
	handleResponse("Rockets", res, err)
}

func handleResponse(queryDesc string, b []byte, err error) {
	if err != nil {
		logrus.WithError(err).Fatal(fmt.Sprintf("Failed to obtain %s", queryDesc))
	}
	logrus.WithField("response", string(b)).Info(queryDesc)
}
