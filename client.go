// Copyright 2023 d-dot-one. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Package awn is a client that can access the Ambient Weather network API and return
// device and weather data.
package awn

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
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
)

// ErrContextTimeoutExceeded is an error message that is returned when the context has
// timed out.
var ErrContextTimeoutExceeded = errors.New("context timeout exceeded")
var ErrMalformedDate = errors.New("date format is malformed. should be YYYY-MM-DD")
var ErrRegexFailed = errors.New("regex failed")

// LogLevelForError is a type that describes the log level for an error message.
type LogLevelForError string

// LogMessage is the message that you would like to see in the log.
type LogMessage string

// YearMonthDay is a type that describes a date in the format YYYY-MM-DD.
type YearMonthDay string

// verify is a private helper function that will verify that the date is in the correct
// format. It will return a boolean value.
func (y YearMonthDay) verify() (bool, error) {
	match, err := regexp.MatchString(`\d{4}-\d{2}-\d{2}`, y.String())
	if err != nil {
		return false, ErrRegexFailed
	}
	if !match {
		return false, ErrMalformedDate
	}
	return true, nil
}

// String is a public helper function that will return the YearMonthDay object as a string.
func (y YearMonthDay) String() string {
	return string(y)
}

// The ConvertTimeToEpoch help function can convert a string, formatted as a time.DateOnly
// object (2023-01-01) to a Unix epoch time in milliseconds. This can be helpful when you
// want to use the GetHistoricalData function to fetch data for a specific date or range
// of dates.
//
// Basic Usage:
//
//	epochTime, err := ConvertTimeToEpoch("2023-01-01")
func ConvertTimeToEpoch(ymd YearMonthDay) (int64, error) {
	result, err := ymd.verify()
	if err != nil {
		return 0, ErrMalformedDate
	}

	if result {
		parsed, err := time.Parse(time.DateOnly, ymd.String())
		_ = CheckReturn(err, "unable to parse time", "warning")
		return parsed.UnixMilli(), err
	}
	return 0, ErrMalformedDate
}

// The CreateAwnClient function is used to create a new resty-based API client. This client
// supports retries and can be placed into debug mode when needed. By default, it will
// also set the accept content type to JSON. Finally, it returns a pointer to the client.
//
// Basic Usage:
//
//	client, err := createAwnClient()
func CreateAwnClient() (*resty.Client, error) {
	client := resty.New().
		SetRetryCount(retryCount).
		SetRetryWaitTime(retryMinWaitTimeSeconds*time.Second).
		SetRetryMaxWaitTime(retryMaxWaitTimeSeconds*time.Second).
		SetBaseURL(baseURL+apiVersion).
		SetHeader("Accept", "application/json").
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
//
// Basic Usage:
//
//	apiConfig := client.CreateApiConfig("apiTokenHere", "appTokenHere")
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
//
// Basic Usage:
//
//	ctx := createContext()
//	ApiConfig := client.CreateApiConfig(apiKey, appKey)
//	data, err := client.GetDevices(ApiConfig)
func GetDevices(ctx context.Context, funcData FunctionData) (AmbientDevice, error) {
	client, err := CreateAwnClient()
	_ = CheckReturn(err, "unable to create client", "warning")

	client.R().SetQueryParams(map[string]string{
		"apiKey":         funcData.API,
		"applicationKey": funcData.App,
	})

	deviceData := &AmbientDevice{}

	_, err = client.R().SetResult(deviceData).Get(devicesEndpoint)
	_ = CheckReturn(err, "unable to handle data from devicesEndpoint", "warning")

	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return nil, fmt.Errorf("%q: %w", deviceData, ErrContextTimeoutExceeded)
	}

	return *deviceData, fmt.Errorf("%w", err)
}

// The getDeviceData function takes a client and the Ambient Weather device MAC address
// as inputs. It then sets the query parameters for authentication and the maximum
// number of records to fetch in this API call to the macAddress endpoint. The response
// data is then marshaled into a pointer to a DeviceDataResponse object which is
// returned to the caller along with any errors.
//
// Basic Usage:
//
//	ctx := createContext()
//	apiConfig := client.CreateApiConfig(apiKey, appKey)
//	resp, err := getDeviceData(ctx, apiConfig)
func getDeviceData(ctx context.Context, funcData FunctionData) (DeviceDataResponse, error) {
	client, err := CreateAwnClient()
	_ = CheckReturn(err, "unable to create client", "warning")

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
	_ = CheckReturn(err, "unable to handle data from the devices endpoint", "warning")

	// CheckResponse(resp) // todo: check call for errors passed through resp

	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return nil, ErrContextTimeoutExceeded
	}

	return *deviceData, fmt.Errorf("%w", err)
}

// GetHistoricalData is a public function takes a FunctionData object as input and
// returns a and will return a list of client.DeviceDataResponse object.
//
// Basic Usage:
//
//	ctx := createContext()
//	apiConfig := client.CreateApiConfig(apiKey, appKey)
//	resp, err := GetHistoricalData(ctx, apiConfig)
func GetHistoricalData(ctx context.Context, funcData FunctionData) ([]DeviceDataResponse, error) {
	var deviceResponse []DeviceDataResponse

	for i := funcData.Epoch; i <= time.Now().UnixMilli(); i += epochIncrement24h {
		funcData.Epoch = i

		resp, err := getDeviceData(ctx, funcData)
		_ = CheckReturn(err, "unable to get device data", "warning")

		deviceResponse = append(deviceResponse, resp)
	}

	return deviceResponse, nil
}

// CheckReturn is a helper function to remove the usual error checking cruft while also
// logging the error message. It takes an error, a message and a log level as inputs and
// returns an error (can be nil of course).
//
// Basic Usage:
//
//	err = CheckReturn(err, "unable to get device data", "warning")
//	if err != nil {
//		log.Printf("Error: %v", err)
//	}
func CheckReturn(err error, msg LogMessage, level LogLevelForError) error {
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

	return err
}

// CheckResponse is a helper function that will take an API response and evaluate it
// for any errors that might have occurred. The API specification does not publish all
// the possible error messages, but these are what I have found so far.
//
// This is not currently implemented.
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
			_ = CheckReturn(err, "unable to get device data", "warning")

			out <- resp
		}
		close(out)
	}()

	return out, nil
}

// GetEnvVars is a public function that will attempt to read the environment variables that
// are passed in as a list of strings. It will return a map of the environment variables.
//
// Basic Usage:
//
//	listOfEnvironmentVariables := GetEnvVars([]string{"ENV_VAR_1", "ENV_VAR_2"})
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
//
// Basic Usage:
//
//	environmentVariable := GetEnvVar("ENV_VAR_1", "fallback")
func GetEnvVar(key string, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}

	return value
}
