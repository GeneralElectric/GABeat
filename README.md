# Gabeat

Welcome to Gabeat.  This is a little process that implements the [Elastic Beat](https://www.elastic.co/products/beats)
interface and gets one data point from the [Google Analytics Real Time Data API](https://developers.google.com/analytics/devguides/reporting/realtime/v3/)
and sends that data to Elastic so that a user can graph events in Elastic against
the data point from Google Analytics (GA).  For example, if you have a dashboard
in Elastic that shows metrics on the errors in your application's logs, you might
want to chart that against page views from GA for your application.

Prerequisites:

  1. Note that this version of GABeat has only been tested with Go 1.7 and Beats 5.2.
  1. Follow the instructions in [Getting Ready](https://www.elastic.co/guide/en/beats/devguide/current/newbeat-getting-ready.html) in the Beats documentation.
  1. Install the non-Go dependences mentioned in [Fetching Dependencies and Setting up the Beat](https://www.elastic.co/guide/en/beats/devguide/current/setting-up-beat.html).  Note that we skipped the "Generating Your Beat" step in the Elastic documentation.  That's because the GABeat structure has already been generated.
  1. Clone this project into the following location: `${GOPATH}/src/github.com/GeneralElectric/GABeat`
  1. Get a GA JWT token (these are the credentials to use the GA APIs) and modify `_meta/beat.yml` google_credentials_file config value to point to it.  The [GA docs](https://developers.google.com/accounts/docs/OAuth2ServiceAccount) explain how to get a token.
  1. Modify the ga_ids, ga_metrics, and ga_dimensions fields of `_meta/beat.yml` to reference your GA account view ID and the data point you want to collect.  To find your account view ID:
    1. Log into [GA](https://analytics.google.com) with your usual credentials.
    1. Click on the account name in the upper left-hand corner of the home page.
    1. Click on accounts -> Properties & Apps -> views.
    1. The view ID is displayed below each view name in the menu.
  1. A note about proxies: You will need to have the http_proxy environment variable set to a [GE-approved proxy](http://internet.ge.com/docs/integration-guide/) to download the code from GE's GitHub.  You will likely not be able to use GE's GitHub if your https_proxy environment variable is set.  However, you WILL need the https_proxy variable set to download the Go dependencies.


### Init Project
To get running with Gabeat and also install the required Go libraries, run the following command:

```
make setup
```

### Build

To build the binary for Gabeat run the command below. This will generate a binary
in the same directory with the name gabeat.

```
make
```


### Run

To run Gabeat with debugging output enabled, run:

```
./gabeat -c gabeat.yml -e -d "*"
```


### Test

To test Gabeat, run the following command:

```
make testsuite
```

alternatively:
```
make unit-tests
make system-tests
make integration-tests
make coverage-report
```

The test coverage is reported in the folder `./build/coverage/`

### Update

Each beat has a template for the mapping in elasticsearch and a documentation for the fields
which is automatically generated based on `etc/fields.yml`.
To generate etc/gabeat.template.json and etc/gabeat.asciidoc

```
make update
```


### Cleanup

To clean  Gabeat source code, run the following commands:

```
make fmt
make simplify
```

To clean up the build directory and generated artifacts, run:

```
make clean
```


## Packaging

NOTE!!!  I have not tested the packaging!!!  

The beat frameworks provides tools to crosscompile and package your beat for different platforms. This requires [docker](https://www.docker.com/) and vendoring as described above. To build packages of your beat, run the following command:

```
make package
```

This will fetch and create all images required for the build process. The hole process to finish can take several minutes.
