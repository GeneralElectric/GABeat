package main

import (
	"os"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/logp"

	"github.com/GeneralElectric/GABeat/beater"
)

func main() {
	logp.Info("Starting GA Beat...")
	err := beat.Run("gabeat", "", beater.New)
	os.Exit(getExitStatus(err))
}

func getExitStatus(err error) (status int) {
	if err != nil {
		logp.Info("Stopping GA Beat with error status...")
		return 1
	}
	logp.Info("Stopping GA Beat normally")
	return 0
}
