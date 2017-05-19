package beater

import (
	"fmt"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"

	"GABeat/config"
	"GABeat/ga"
)

var debugf = logp.MakeDebug("gabeat")

type gaDataRetriever func(gaConfig config.GoogleAnalyticsConfig) (gaDataPoints []ga.GABeatDataPoint, err error)

type Gabeat struct {
	done   chan struct{}
	config config.Config
	client publisher.Client
}

// Creates beater
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	config := config.DefaultConfig
	if err := cfg.Unpack(&config); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &Gabeat{
		done:   make(chan struct{}),
		config: config,
	}
	logp.Info("Config: %s", bt.config)
	return bt, nil
}

func (bt *Gabeat) Run(b *beat.Beat) error {
	return runFunctionally(bt, b, ga.GetGAReportData)
}

func runFunctionally(bt *Gabeat, b *beat.Beat, dataFunc gaDataRetriever) error {
	logp.Info("gabeat is running! Hit CTRL-C to stop it.")
	bt.client = b.Publisher.Connect()
	ticker := time.NewTicker(bt.config.Period)
	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
		} //end select
		beatOnce(bt.client, b.Name, bt.config.Googleanalytics, dataFunc)
	} //end for
} //end func

func beatOnce(client publisher.Client, beatName string, gaConfig config.GoogleAnalyticsConfig, dataFunc gaDataRetriever) {
	GAData, err := dataFunc(gaConfig)
	if err == nil {
		publishToElastic(client, beatName, GAData)
	} else {
		logp.Err("gadata was null, not publishing: %v", err)
	}

}

func makeEvent(beatType string, GAData []ga.GABeatDataPoint) map[string]interface{} {
	event := common.MapStr{
		"@timestamp": common.Time(time.Now()),
		"type":       beatType,
		"count":      1,  //The number of transactions that this event represents
	}
	for _, gaDataPoint := range GAData {
		gaDataName := gaDataPoint.DimensionName + "_" + gaDataPoint.MetricName
		event.Put(gaDataName, gaDataPoint.Value)
	}
	return event
}

func publishToElastic(client publisher.Client, beatType string, GAData []ga.GABeatDataPoint) {
	event := makeEvent(beatType, GAData)
	succeeded := client.PublishEvent(event)
	if !succeeded {
		logp.Err("Publisher couldn't publish event to Elastic")
	}
}

func (bt *Gabeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
