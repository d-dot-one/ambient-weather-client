// Copyright 2023 d-dot-one. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Package ambient_weather_client is a client that can access the Ambient Weather
// network API and return device and weather data.
package awn

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	// apiVersion is a string and describes the version of the API that Ambient Weather
	// is using.
	apiVersion = "/v1"

	// baseURL The base URL for the Ambient Weather API (Not the real-time API) as a string.
	baseURL = "https://rt.ambientweather.net"

	// debugMode Enable verbose logging by setting this boolean value to true.
	debugMode = false

	// defaultCtxTimeout Set the context timeout, in seconds, as an int.
	defaultCtxTimeout = 30

	// devicesEndpoint The 'devices' endpoint as a string.
	devicesEndpoint = "devices"

	// epochIncrement24h is the number of milliseconds in a 24-hour period.
	epochIncrement24h int64 = 86400000

	// retryCount An integer describing the number of times to retry in case of
	// failure or rate limiting.
	retryCount = 3

	// retryMaxWaitTimeSeconds An integer describing the maximum time to wait to
	// retry an API call, in seconds.
	retryMaxWaitTimeSeconds = 15

	// retryMinWaitTimeSeconds An integer describing the minimum time to wait to retry
	// an API call, in seconds.
	retryMinWaitTimeSeconds = 5

	// Version # of client
	Version = "0.0.1"
)

// The ConvertTimeToEpoch help function can convert any Go time.Time object to a Unix epoch time in milliseconds.
// func ConvertTimeToEpoch(t time.Time) int64 {
//	return t.UnixMilli()
//}

// The createAwnClient function is used to create a new resty-based API client. This client
// supports retries and can be placed into debug mode when needed. By default, it will
// also set the accept content type to JSON. Finally, it returns a pointer to the client.
func createAwnClient() (*resty.Client, error) {
	client := resty.New().
		SetRetryCount(retryCount).
		SetRetryWaitTime(retryMinWaitTimeSeconds * time.Second).
		SetRetryMaxWaitTime(retryMaxWaitTimeSeconds * time.Second).
		SetBaseURL(baseURL + apiVersion).
		SetTimeout(defaultCtxTimeout * time.Second).
		SetDebug(debugMode).
		AddRetryCondition(
			func(r *resty.Response, e error) bool {
				return r.StatusCode() == http.StatusRequestTimeout ||
					r.StatusCode() >= http.StatusInternalServerError ||
					r.StatusCode() == http.StatusTooManyRequests
			})

	client.SetHeader("Accept", "application/json")

	return client, nil
}

// CreateAPIConfig is a helper function that is used to create the FunctionData struct,
// which is passed to the data gathering functions. It takes as parameters the API key
// as api and the Application key as app and returns a pointer to a FunctionData object.
func CreateAPIConfig(api string, app string) *FunctionData {
	fd := NewFunctionData()
	fd.API = api
	fd.App = app

	return fd
}

// GetDevices is a public function takes a client, sets the appropriate query parameters
// for authentication, makes the request to the devicesEndpoint endpoint and marshals the
// response data into a pointer to an AmbientDevice object, which is returned along with
// any error messages.
func GetDevices(ctx context.Context, funcData FunctionData) (AmbientDevice, error) {
	client, err := createAwnClient()
	CheckReturn(err, "unable to create client", "warning")

	client.R().SetQueryParams(map[string]string{
		"apiKey":         funcData.API,
		"applicationKey": funcData.App,
	})

	deviceData := &AmbientDevice{}

	_, err = client.R().SetResult(deviceData).Get(devicesEndpoint)
	CheckReturn(err, "unable to handle data from devicesEndpoint", "warning")

	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return nil, errors.New("context timeout exceeded")
	}

	return *deviceData, err
}

// The getDeviceData function takes a client and the Ambient Weather device MAC address
// as inputs. It then sets the query parameters for authentication and the maximum
// number of records to fetch in this API call to the macAddress endpoint. The response
// data is then marshaled into a pointer to a DeviceDataResponse object which is
// returned to the caller along with any errors.
func getDeviceData(ctx context.Context, funcData FunctionData) (DeviceDataResponse, error) {
	client, err := createAwnClient()
	CheckReturn(err, "unable to create client", "warning")

	client.R().SetQueryParams(map[string]string{
		"apiKey":         funcData.API,
		"applicationKey": funcData.App,
		"endDate":        strconv.FormatInt(funcData.Epoch, 10),
		"limit":          strconv.Itoa(funcData.Limit),
	})

	deviceData := &DeviceDataResponse{}

	_, err = client.R().
		SetPathParams(map[string]string{
			"devicesEndpoint": devicesEndpoint,
			"macAddress":      funcData.Mac,
		}).
		SetResult(deviceData).
		Get("{devicesEndpoint}/{macAddress}")
	CheckReturn(err, "unable to handle data from the devices endpoint", "warning")
	// todo: check call for errors passed through resp
	// if mac is missing, you get the devices endpoint response, so test for mac address
	// if apiKey is missing, you get {"error": "apiKey-missing"}
	// if appKey is missing, you get {"error": "applicationKey-missing"}
	// if date is wrong, you get {"error":"date-invalid","message":"Please refer
	//		to: http://momentjs.com/docs/#/parsing/string/"}

	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return nil, errors.New("ctx timeout exceeded")
	}

	return *deviceData, err
}

// GetHistoricalData is a public function takes a FunctionData object as input and
// returns a and will return a list of client.DeviceDataResponse object.
func GetHistoricalData(ctx context.Context, funcData FunctionData) ([]DeviceDataResponse, error) {
	var deviceResponse []DeviceDataResponse

	for i := funcData.Epoch; i <= time.Now().UnixMilli(); i += epochIncrement24h {
		funcData.Epoch = i

		resp, err := getDeviceData(ctx, funcData)
		CheckReturn(err, "unable to get device data", "warning")

		deviceResponse = append(deviceResponse, resp)
	}

	return deviceResponse, nil
}

type LogLevelForError string

// CheckReturn is a helper function to remove the usual error checking cruft.
func CheckReturn(err error, msg string, level LogLevelForError) {
	if err != nil {
		switch level {
		case "panic":
			log.Panicf("%v: %v", msg, err)
		case "fatal":
			log.Fatalf("%v: %v", msg, err)
		case "warning":
			log.Printf("%v: %v\n", msg, err)
		case "info":
			log.Printf("%v: %v\n", msg, err)
		case "debug":
			log.Printf("%v: %x\n", msg, err)
		}
	}
}

// CheckResponse is a helper function that will take an API response and evaluate it for
// for any errors that might have occurred. The API specification does not publish all of
// the possible error messages, but these are what I have found so far.
func CheckResponse(resp map[string]string) bool {
	message, ok := resp["error"]
	if ok {
		switch message {
		case "apiKey-missing":
			log.Panicf("API key is missing (%v). Visit https://ambientweather.net/account", message)
		case "applicationKey-missing":
			log.Panicf("App key is missing (%v). Visit https://ambientweather.net/account", message)
		case "date-invalid":
			log.Panicf("Date is invalid (%v). It should be in epoch time in milliseconds", message)
		case "macAddress-missing":
			log.Panicf("MAC address is missing (%v). Supply a valid MAC address for a weather station", message)
		default:
			return true
		}
	}

	return true
}

// GetHistoricalDataAsync is a public function that takes a context object, a FunctionData
// object and a WaitGroup object as inputs. It will return a channel of DeviceDataResponse
// objects and an error status.
func GetHistoricalDataAsync(
	ctx context.Context,
	funcData FunctionData,
	w *sync.WaitGroup) (<-chan DeviceDataResponse, error) {
	defer w.Done()

	out := make(chan DeviceDataResponse)

	go func() {
		for i := funcData.Epoch; i <= time.Now().UnixMilli(); i += epochIncrement24h {
			funcData.Epoch = i

			resp, err := getDeviceData(ctx, funcData)
			CheckReturn(err, "unable to get device data", "warning")

			out <- resp
		}
		close(out)
	}()

	return out, nil
}

// GetEnvVars is a public function that will attempt to read the environment variables that
// are passed in as a list of strings. It will return a map of the environment variables.
func GetEnvVars(vars []string) map[string]string {
	envVars := make(map[string]string)

	for v := range vars {
		value := GetEnvVar(vars[v], "")
		envVars[vars[v]] = value
	}

	return envVars
}

// GetEnvVar is a public function attempts to fetch an environment variable. If that
// environment variable is not found, it will return 'fallback'.
func GetEnvVar(key string, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}

	return value
}
