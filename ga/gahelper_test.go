package ga

// This file is mandatory as otherwise the gabeat.test binary is not generated correctly.

import (
	cfg "github.com/GeneralElectric/GABeat/config"
	"errors"
	"github.com/stretchr/testify/assert"
	analytics "google.golang.org/api/analytics/v3"
	googleapi "google.golang.org/api/googleapi"
	"os"
	"testing"
)

var gaTestConfig cfg.GoogleAnalyticsConfig = cfg.GoogleAnalyticsConfig{"/", "ids", "metric", "dimension"}
const metricNameFormatted = "rt_pageviews"
const metricName = "rt:pageViews"
const dimensionNameFormatted = "pagename"
const dimensionName = "PageName"

func TestValidateCorrectConfig(t *testing.T) {
	gadata, err := validateConfig("foo", "fooConfig")
	assert.Nil(t, err, "validation of corerct config failed %v", err)
	assert.True(t, gadata[0].Value == 0, "validation of corerct config failed, gaData=%d", gadata)
}

func TestValidateIncorrectConfig(t *testing.T) {
	gadata, err := validateConfig("", "emptyConfig")
	assert.NotNil(t, err, "validation of empty config failed %v", err)
	assert.True(t, gadata[0].Value == -1, "validation of empty config failed, gaData=%d", gadata)
}

func TestFormat(t *testing.T) {
	result := format(metricName)
	assert.EqualValues(t, metricNameFormatted, result, "formatting failed")
}

func TestFormatEventMetric(t *testing.T) {
	result := format(" ab*$&%  (c) d:")
	assert.EqualValues(t, "_ab__c_d_", result, "formatting failed")
}

func TestParseGAResponse(t *testing.T) {
	gaDataPoints, err := getGARealtimeData("3")
	assert.Nil(t, err, "parseGAResponse falied with error %v", err)
	testDataPoint(3, gaDataPoints, 0, t)
	testMetric(metricNameFormatted, gaDataPoints, 0, t)
	testDimension(dimensionNameFormatted, gaDataPoints, 0, t)
}

func TestParseGAMulitRowResponse(t *testing.T) {
	gaDataPoints, err := getGAMultiRowRealtimeData()
	assert.Nil(t, err, "parseGAResponse falied with error %v", err)
	testDataPoint(63, gaDataPoints, 0, t)
	testDataPoint(2, gaDataPoints, 1, t)
	testDataPoint(35, gaDataPoints, 2, t)
	testMetric("rt_totalevents", gaDataPoints, 0, t)
	testMetric("rt_totalevents", gaDataPoints, 1, t)
	testMetric("rt_totalevents", gaDataPoints, 2, t)
	testDimension("not_set_not_set", gaDataPoints, 0, t)
	testDimension("about_ge_our_company", gaDataPoints, 1, t)
	testDimension("actions_open", gaDataPoints, 2, t)
}

func testDataPoint(expected int, dataPoints []GABeatDataPoint, index int, t *testing.T){
	assert.EqualValues(t, expected, dataPoints[index].Value,
		"data point should be %d but was: %d", expected, dataPoints[index].Value)
}

func testMetric(expected string, dataPoints []GABeatDataPoint, index int, t *testing.T){
	assert.EqualValues(t, expected, dataPoints[index].MetricName,
		"metric should be %s but was: %s", expected, dataPoints[index].MetricName)
}

func testDimension(expected string, dataPoints []GABeatDataPoint, index int, t *testing.T){
	assert.EqualValues(t, expected, dataPoints[index].DimensionName,
		"dimension should be %s but was: %s", expected, dataPoints[index].DimensionName)
}

func TestParseGAResponse_NaNShouldFail(t *testing.T) {
	gaDataPoints, err := getGARealtimeData("NaN")
	assert.NotNil(t, err, "parseGAResponse should have failed")
	assert.EqualValues(t, -1, gaDataPoints[0].Value,
		"data point should be -1 but was: %d", gaDataPoints[0].Value)
}

func createGASingleRowRealtimeData(value string) *analytics.RealtimeData {
	rows := [][]string{}
	row1 := []string{dimensionName, value}
	rows = append(rows, row1)
	columnHeaders := []*analytics.RealtimeDataColumnHeaders{
		createDimensionColumnHeader("rt:pageTitle"),
		createMetricColumnHeader("rt:pageViews")}
	return createGARealtimeData(columnHeaders, "3", rows)
}

func createGAMulitRowRealtimeData() *analytics.RealtimeData {
	rows := [][]string{}
	row1 := []string{"(not set)", "(not set)", "63"}
	row2 := []string{"About GE", "Our Company", "2"}
	row3 := []string{"Actions", "Open", "35"}
	rows = append(rows, row1, row2, row3)
	columnHeaders := []*analytics.RealtimeDataColumnHeaders{
		createDimensionColumnHeader("rt:eventAction"),
		createDimensionColumnHeader("rt:eventLabel"),
		createMetricColumnHeader("rt:totalEvents")}
	return createGARealtimeData(columnHeaders, "100", rows)
}

func createGAEmptyRealtimeData() *analytics.RealtimeData {
	rows := [][]string{}
	columnHeaders := []*analytics.RealtimeDataColumnHeaders{}
	return createGARealtimeData(columnHeaders, "0", rows)
}

func createGARealtimeData(columnHeaders []*analytics.RealtimeDataColumnHeaders,
	totalsForAllResultsValue string,
	rows [][]string) *analytics.RealtimeData {
	var totalsForAllResults = map[string]string{metricName: totalsForAllResultsValue}
	var header = map[string][]string{"foo": {"bar", "baz"}}
	var serverResponse = googleapi.ServerResponse{200, header}
	realtimeData := &analytics.RealtimeData{
		columnHeaders,       //ColumnHeaders []*RealtimeDataColumnHeaders `json:"columnHeaders,omitempty"`
		"",                  //Id string `json:"id,omitempty"`
		"",                  //Kind string `json:"kind,omitempty"`
		nil,                 //ProfileInfo *RealtimeDataProfileInfo `json:"profileInfo,omitempty"`
		nil,                 //Query *RealtimeDataQuery `json:"query,omitempty"`
		rows,                //Rows [][]string `json:"rows,omitempty"`
		"",                  //SelfLink string `json:"selfLink,omitempty"`
		100,                 //TotalResults int64 `json:"totalResults,omitempty"`
		totalsForAllResults, //TotalsForAllResults map[string]string `json:"totalsForAllResults,omitempty"`
		serverResponse,      //googleapi.ServerResponse `json:"-"`
		nil,                 //ForceSendFields []string `json:"-"`
		nil,                 //NullFields is a list of field names (e.g. "ColumnHeaders") to include in API requests with the JSON null value.
	}
	return realtimeData
}

func createMetricColumnHeader(name string) *analytics.RealtimeDataColumnHeaders {
	return createColumnHeader(name, "METRIC", "INTEGER")
}

func createDimensionColumnHeader(name string) *analytics.RealtimeDataColumnHeaders {
	return createColumnHeader(name, "DIMENSION", "STRING")
}

func createColumnHeader(name string, columnType string, dataType string) *analytics.RealtimeDataColumnHeaders {
	return &analytics.RealtimeDataColumnHeaders{
		columnType, //ColumnType string `json:"columnType,omitempty"`
		dataType, //DataType string `json:"dataType,omitempty"`
		name, //Name string `json:"name,omitempty"`
		nil, //ForceSendFields []string `json:"-"`
		nil, //NullFields []string `json:"-"`
	}
}

func getGARealtimeData(value string) ([]GABeatDataPoint, error) {
	realtimeData := createGASingleRowRealtimeData(value)
	gaDataPoint, err := parseGAResponse(realtimeData)
	return gaDataPoint, err
}

func getGAMultiRowRealtimeData() ([]GABeatDataPoint, error) {
	realtimeData := createGAMulitRowRealtimeData()
	gaDataPoint, err := parseGAResponse(realtimeData)
	return gaDataPoint, err
}

func TestInitCredentialsPath(t *testing.T) {

	_, err := initCredentialsPath(gaTestConfig)
	assert.Nil(t, err, "init credentials test failed")
	newEnv := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	assert.EqualValues(t, gaTestConfig.GoogleCredentialsFilePath, newEnv, "env var not set")
}

func TestInitCredentialsEmptyPathShouldFail(t *testing.T) {
	testInitCredentialsPathShouldFail(t, "", "empty path should fail validation")
}

func TestInitCredentialsMissingPathShouldFail(t *testing.T) {
	testInitCredentialsPathShouldFail(t, "foo", "incorrect path should fail validation")
}

func testInitCredentialsPathShouldFail(t *testing.T, pathValue string, message string) {

	gaTestConfigFail := cfg.GoogleAnalyticsConfig{pathValue, "ids", "metric", "dimension"}

	result, err := initCredentialsPath(gaTestConfigFail)
	assert.NotNil(t, err, message)
	assert.EqualValues(t, -1, result[0].Value, message)
}

func TestGAReturnsSuccess(t *testing.T) {
	dataPoints, err := getGAReportData(gaTestConfig, getGADataSuccess)
	assert.Nil(t, err, "Error getting GA report data %v", err)
	assert.EqualValues(t, 3, dataPoints[0].Value,
		"data point should be 3 but was: %d", dataPoints[0].Value)
}

func TestGAReturnsEmpty(t *testing.T) {
	dataPoints, err := getGAReportData(gaTestConfig, getGADataEmpty)
	assert.Nil(t, err, "Error getting GA report data %v", err)
	assert.EqualValues(t, len(emptyResults), len(dataPoints),
		"length of data points should be %d but was %d", len(emptyResults), len(dataPoints))
	testDataPoint(emptyResult.Value, dataPoints, 0, t)
	testMetric(emptyResult.MetricName, dataPoints, 0, t)
	testDimension(emptyResult.DimensionName, dataPoints, 0, t)
}

func TestGAReturnsFail(t *testing.T) {
	_, err := getGAReportData(gaTestConfig, getGADataFail)
	assert.NotNil(t, err, "Should have failed to get GA data")
}

func TestBadCredsPath(t *testing.T) {
	var gaTestBadCredsConfig cfg.GoogleAnalyticsConfig = cfg.GoogleAnalyticsConfig{"foo", "ids", "metric", "dimension"}
	testBadConfig(t, gaTestBadCredsConfig, "Should have failed credentials path validation")
}

func TestEmptyIDsConfig(t *testing.T) {
	var gaTestEmptyIDsConfig cfg.GoogleAnalyticsConfig = cfg.GoogleAnalyticsConfig{"/", "", "metric", "dimension"}
	testBadConfig(t, gaTestEmptyIDsConfig, "Should have failed empty id string in config")
}

func TestEmptyMetricConfig(t *testing.T) {
	var badConfig cfg.GoogleAnalyticsConfig = cfg.GoogleAnalyticsConfig{"/", "id", "", "dimension"}
	testBadConfig(t, badConfig, "Should have failed empty metrics string in config")
}

func TestEmptyDimensionConfig(t *testing.T) {
	var badConfig cfg.GoogleAnalyticsConfig = cfg.GoogleAnalyticsConfig{"/", "id", "metrics", ""}
	testBadConfig(t, badConfig, "Should have failed empty dimensions string in config")
}

func testBadConfig(t *testing.T, config cfg.GoogleAnalyticsConfig, message string) {
	_, err := getGAReportData(config, getGADataSuccess)
	assert.NotNil(t, err, message)
}

func getGADataSuccess(gaIds string, gaMetrics string, gaDimensions string) (gaData *analytics.RealtimeData, err error) {
	return createGASingleRowRealtimeData("3"), nil
}

func getGADataEmpty(gaIds string, gaMetrics string, gaDimensions string) (gaData *analytics.RealtimeData, err error) {
	return createGAEmptyRealtimeData(), nil
}

func getGADataFail(gaIds string, gaMetrics string, gaDimensions string) (gaData *analytics.RealtimeData, err error) {
	return nil, errors.New("Testing failure condition")
}
