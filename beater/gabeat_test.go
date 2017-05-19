package beater

import (
	"fmt"
	"testing"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/publisher"

	"GABeat/config"
	"GABeat/ga"

	"github.com/stretchr/testify/assert"
)

func TestConstructor(t *testing.T) {

	parentBeat := &beat.Beat{}
	parentConfig := common.NewConfig()

	constructed, err := New(parentBeat, parentConfig)
	assert.Nil(t, err, "Error constructing gabeat: %v", err)
	assert.NotNil(t, constructed)
}

func TestRunAndStop(t *testing.T) {

	runResult := make(chan bool)

	parentBeat := &beat.Beat{}
	stubPub := &stubPublisher{}
	parentBeat.Publisher = stubPub
	parentConfig := common.NewConfig()

	constructed, err := New(parentBeat, parentConfig)
	assert.Nil(t, err, "Error constructing gabeat: %v", err)
	assert.NotNil(t, constructed)

	gabeater, typeAssertOK := constructed.(*Gabeat)
	assert.True(t, typeAssertOK, "Error converting Beater to Gabeat")

	go func() {
		result := runFunctionally(gabeater, parentBeat, mockDataRetriever)
		if result == nil {
			runResult <- true
		} else {
			runResult <- false
		}
	}()
	time.Sleep(2 * time.Second)
	constructed.Stop()
	assert.True(t, stubPub.isConnected, "The stub publisher should have connected")
	assert.True(t, <-runResult, "The Run method should not have returned an error")
	assert.NotNil(t, stubPub.client.events, "Slice of published events should not be nil")
	eventsLength := len(stubPub.client.events)
	assert.True(t, (eventsLength > 0 && eventsLength < 4),
		fmt.Sprintf("Slice of published events should be between 1 and 4 but is %d", eventsLength))
}

type stubPublisher struct {
	isConnected bool
	client      *stubClient
}

func (publisher *stubPublisher) Connect() publisher.Client {
	publisher.isConnected = true
	publisher.client = &stubClient{events: make([]common.MapStr, 0, 10)}
	return publisher.client
}

type stubClient struct {
	events []common.MapStr
}

func (c *stubClient) Close() error {
	return nil
}

func (c *stubClient) PublishEvent(event common.MapStr, opts ...publisher.ClientOption) bool {
	c.events = append(c.events, event)
	return true
}

func (c *stubClient) PublishEvents(events []common.MapStr, opts ...publisher.ClientOption) bool {
	copy(c.events, events[:])
	return true
}

func mockDataRetriever(gaConfig config.GoogleAnalyticsConfig) (gaDataPoint []ga.GABeatDataPoint, err error) {
	rows := []ga.GABeatDataPoint{}
	row1 := ga.GABeatDataPoint{5, "dimension_name", "metric_name"}
	rows = append(rows, row1)
	return rows, nil
}
