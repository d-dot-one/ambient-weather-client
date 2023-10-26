# Ambient Weather Network API Client

## WIP
1662409800000 = 2022-09-06 14:30:00
1662496200000 = 2022-09-06 14:30:00

## Overview

This is a feature-complete Go version of a client that can connect to the Ambient Weather Network API in order to pull information about your weather station and the data that it has collected. It supports the normal API as well as the Websockets-based realtime API.

## Installation

```bash
go get github.com/d-dot-one/ambient-weather-network-client
```
... or you can simply import it in your project and use it.

```go
import "github.com/d-dot-one/ambient-weather-network-client"
```
You'll need to do a `go get` in the terminal to actually fetch the package.

## Environment Variables

In order for all of this to work, you will need the following environment variables:

| Variable        | Required | Description                                  |
|-----------------|----------|----------------------------------------------|
| `AWN_API_KEY`   | Yes      | Your Ambient Weather Network API key         |
| `AWN_APP_KEY`   | Yes      | Your Ambient Weather Network application key |
| `AWN_LOG_LEVEL` | No       | The log level to use. Defaults to `info`     |


## Usage

### Get Weather Station Data
To fetch the current weather and the weather station device data, you can use the following code:

```go
package main

import (
	"context"
	"fmt"

	client "github.com/d-dot-one/ambient-weather-network-client"
)

func main() {
	// create a context
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()
	
	// fetch required environment variables
	requiredVars := []string{"AWN_API_KEY", "AWN_APP_KEY", "AWN_LOG_LEVEL"}
	environmentVariables := client.GetEnvVars(requiredVars)

	// set the API key
	apiKey := fmt.Sprintf("%v", environmentVariables["AWN_API_KEY"])

	// set the application key
	appKey := fmt.Sprintf("%v", environmentVariables["AWN_APP_KEY"])

	// create an object to hold the API configuration
	ApiConfig := client.CreateApiConfig(apiKey, appKey, ctx)

	// fetch the device data and return it as an AmbientDevice
	data, err := client.GetDevices(ApiConfig)
	client.CheckReturn(err, "failed to get devices", "critical")
	
	// see the MAC address of the weather station
	fmt.Println(data.MacAddress)
}
```

### Get Historical Weather Station Data
```go
package main

import (
	"context"
	"fmt"

	client "github.com/d-dot-one/ambient-weather-network-client"
)

func main() {
	// create a context
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()
	
	// fetch required environment variables
	requiredVars := []string{"AWN_API_KEY", "AWN_APP_KEY", "AWN_LOG_LEVEL"}
	environmentVariables := client.GetEnvVars(requiredVars)

	// set the API key
	apiKey := fmt.Sprintf("%v", environmentVariables["AWN_API_KEY"])

	// set the application key
	appKey := fmt.Sprintf("%v", environmentVariables["AWN_APP_KEY"])

	// create an object to hold the API configuration
	ApiConfig := client.CreateApiConfig(apiKey, appKey, ctx)

	// fetch the device data and return it as an AmbientDevice
	data, err := client.GetDevices(ApiConfig)
	client.CheckReturn(err, "failed to get devices", "critical")
	
	// see the MAC address of the weather station
	fmt.Println(data.MacAddress)
}
```

## Dependencies

I purposefully chose to use as few dependencies as possible for this project. I wanted to keep it as simple and close to the standard library. The only exception is the `resty` library, which is used to make the API calls. It was too helpful with retries to not use it.

## Constrictions

The Ambient Weather API has a cap on the number of API calls that one can make in a given second. This is set to 1 call per second. This means that if you have more than one weather station, you will need to make sure that you are not making more than 1 call per second. This is done by using a `time.Sleep(1 * time.Second)` after each API call. Generally speaking, this is all handled in the background with a retry mechanism, but it is something to be aware of.

## Contributing

You are more than welcome to submit a PR if you would like to contribute to this project. I am always open to suggestions and improvements. Reference [CONTRIBUTING.md](./CONTRIBUTING.md) for more information.

## License

This package is made available under an MIT license. See [LICENSE.md](./LICENSE.md) for more information.
