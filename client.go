// Copyright 2023 d-dot-one. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Package awn is a client that can access the Ambient Weather network API and return
// device and weather data.
package awn

import (
	"context"
	"errors"
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
	// apiVersion is a string and describes the version of the API that Ambient
	// Weather is using.
	//apiVersion = "/v1"

	// baseURL The base URL for the Ambient Weather API (Not the real-time API)
	// as a string.
	//baseURL = "https://rt.ambientweather.net"

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

	// retryMinWaitTimeSeconds An integer describing the minimum time to wait
	// to retry an API call, in seconds.
	retryMinWaitTimeSeconds = 5
)

var (
	// ErrContextTimeoutExceeded is an error message that is returned when
	// the context has timed out.
	ErrContextTimeoutExceeded = errors.New("context timeout exceeded")

	// ErrMalformedDate is a custom error and message that is returned when a
	// date is passed that does not conform to the required format.
	ErrMalformedDate = errors.New("date format is malformed. should be YYYY-MM-DD")

	// ErrRegexFailed is a custom error and message that is returned when regex
	// fails catastrophically.
	ErrRegexFailed = errors.New("regex failed")

	// ErrAPIKeyMissing is a custom error and message that is returned when no API key is
	// passed to a function that requires it.
	ErrAPIKeyMissing = errors.New("api key is missing")

	// ErrAppKeyMissing is a custom error and message that is returned when no application
	// key is passed to a function that requires it.
	ErrAppKeyMissing = errors.New("application key is missing")

	// ErrInvalidDateFormat is a custom error and message that is returned when the date
	// is not passed as an epoch time in milliseconds.
	ErrInvalidDateFormat = errors.New("date is invalid. It should be in epoch time in milliseconds")

	// ErrMacAddressMissing is a custom error and message that is returned when no MAC
	// address is passed to a function that requires it.
	ErrMacAddressMissing = errors.New("mac address missing")
)

type (
	// LogLevelForError is a type that describes the log level for an error message.
	LogLevelForError string

	// LogMessage is the message that you would like to see in the log.
	LogMessage string

	// YearMonthDay is a type that describes a date in the format YYYY-MM-DD.
	YearMonthDay string
)

// verify is a private helper function that will check that the date string passed from
// the caller is in the correct format. It will return a boolean value and an error.
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

// String is a public helper function that will return the YearMonthDay object
// as a string.
func (y YearMonthDay) String() string {
	return string(y)
}

// The ConvertTimeToEpoch public helper function that can convert a string, formatted
// as a time.DateOnly object (i.e. "2023-01-01") to a Unix epoch time in milliseconds.
// This can be helpful when you want to use the GetHistoricalData function to
// fetch data for a specific date or range of dates.
//
// Basic Usage:
//
//	epochTime, err := ConvertTimeToEpoch("2023-01-01")
func ConvertTimeToEpoch(t string) (int64, error) {
	ok, err := YearMonthDay(t).verify()
	_ = CheckReturn(err, "unable to verify date", "warning")

	if !ok {
		log.Fatalf("invalid date format, %v should be YYYY-MM-DD", t)
	}

	parsed, err := time.Parse(time.DateOnly, t)
	_ = CheckReturn(err, "unable to parse time", "warning")

	return parsed.UnixMilli(), err
}

// CreateAwnClient is a public function that is used to create a new resty-based API
// client. It takes the URL that you would like to connect to and the API version as inputs
// from the caller. This client supports retries and can be placed into debug mode when
// needed. By default, it will also set the accept content type to JSON. Finally, it
// returns a pointer to the client and an error.
//
// Basic Usage:
//
//	client, err := createAwnClient()
func CreateAwnClient(url string, version string) (*resty.Client, error) {
	client := resty.New().
		SetRetryCount(retryCount).
		SetRetryWaitTime(retryMinWaitTimeSeconds*time.Second).
		SetRetryMaxWaitTime(retryMaxWaitTimeSeconds*time.Second).
		SetBaseURL(url+version).
		SetHeader("Accept", "application/json").
		SetTimeout(defaultCtxTimeout * time.Second).
		SetDebug(debugMode).
		AddRetryCondition(
			func(r *resty.Response, e error) bool {
				return r.StatusCode() == http.StatusRequestTimeout ||
					r.StatusCode() >= http.StatusInternalServerError ||
					r.StatusCode() == http.StatusTooManyRequests
			})

	return client, nil
}

// CreateAPIConfig is a public helper function that is used to create the FunctionData
// struct, which is passed to the data gathering functions. It takes as parameters the
// API key as "api" and the Application key as "app" and returns a pointer to a
// FunctionData object.
//
// Basic Usage:
//
//	apiConfig := awn.CreateApiConfig("apiTokenHere", "appTokenHere")
func CreateAPIConfig(api string, app string) *FunctionData {
	fd := NewFunctionData()
	fd.API = api
	fd.App = app

	return fd
}

// CheckReturn is a public function to remove the usual error checking cruft while also
// logging the error message. It takes an error, a message and a log level as inputs and
// returns an error (can be nil of course). You can then use the err message for custom
// handling of the error.
//
// Basic Usage:
//
//	err = CheckReturn(err, "unable to get device data", "warning")
func CheckReturn(err error, msg string, level LogLevelForError) error {
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

// CheckResponse is a public function that will take an API response and evaluate it
// for any errors that might have occurred. The API specification does not publish all
// the possible error messages, but these are what I have found so far. It returns a
// boolean that indicates if the response has an error or not and an error message, if
// applicable.
//
// This is not currently implemented.
func CheckResponse(resp map[string]string) (bool, error) {
	message, ok := resp["error"]
	if ok {
		switch message {
		case "apiKey-missing":
			log.Panicf("API key is missing (%v). Visit https://ambientweather.net/account", message)
			return false, ErrAPIKeyMissing
		case "applicationKey-missing":
			log.Panicf("App key is missing (%v). Visit https://ambientweather.net/account", message)
			return false, ErrAppKeyMissing
		case "date-invalid":
			log.Panicf("Date is invalid (%v). It should be in epoch time in milliseconds", message)
			return false, ErrInvalidDateFormat
		case "macAddress-missing":
			log.Panicf("MAC address is missing (%v). Supply a valid MAC address for a weather station", message)
			return false, ErrMacAddressMissing
		default:
			return false, nil
		}
	}

	return true, nil
}

// GetLatestData is a public function that takes a context object, a FunctionData object, a
// URL and an API version route as inputs. It then creates an AwnClient and sets the
// appropriate query parameters for authentication, makes the request to the
// devicesEndpoint endpoint and marshals the response data into a pointer to an
// AmbientDevice object, which is returned along with any error message.
//
// This function can be used to get the latest data from the Ambient Weather Network API.
// But, it is generally used to get the MAC address of the weather station that you would
// like to get historical data from.
//
// Basic Usage:
//
//	ctx := createContext()
//	apiConfig := awn.CreateApiConfig(apiKey, appKey)
//	data, err := awn.GetLatestData(ctx, ApiConfig, baseURL, apiVersion)
func GetLatestData(ctx context.Context, funcData FunctionData, url string, version string) (*AmbientDevice, error) {
	client, err := CreateAwnClient(url, version)
	_ = CheckReturn(err, "unable to create client", "warning")

	client.R().SetQueryParams(map[string]string{
		"apiKey":         funcData.API,
		"applicationKey": funcData.App,
	})

	deviceData := &AmbientDevice{}

	_, err = client.R().SetResult(deviceData).Get(devicesEndpoint)
	_ = CheckReturn(err, "unable to handle data from devicesEndpoint", "warning")

	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return nil, errors.New("context timeout exceeded")
	}

	return deviceData, err
}

// getDeviceData is a private function takes a context object, a FunctionData object, a URL
// for the Ambient Weather Network API and the API version route as inputs. It creates the
// API client, then sets the query parameters for authentication and the maximum
// number of records to fetch in each API call to the macAddress endpoint. The response
// data is then marshaled into a pointer to a DeviceDataResponse object which is
// returned to the caller along with any errors.
//
// This function should be used if you are looking for weather data from a specific date
// or time. The "limit" parameter can be a number from 1 to 288. You should discover how
// often your weather station updates data in order to get a better understanding of how
// many records will be fetched. For example, if your weather station updates every 5
// minutes, then 288 will give you 24 hours of data. However, many people upload weather
// data less frequently, skewing this length of time.
//
// Basic Usage:
//
//	ctx := createContext()
//	apiConfig := awn.CreateApiConfig(apiKey, appKey)
//	resp, err := getDeviceData(ctx, apiConfig)
func getDeviceData(ctx context.Context, funcData FunctionData, url string, version string) (DeviceDataResponse, error) {
	client, err := CreateAwnClient(url, version)
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

	//CheckResponse(resp) // todo: check response for errors passed through resp

	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return DeviceDataResponse{}, ErrContextTimeoutExceeded
	}

	return *deviceData, err
}

// GetHistoricalData is a public function that takes a context object, a FunctionData
// object, the URL of the Ambient Weather Network API and the API version route as inputs
// and returns a list of DeviceDataResponse objects and an error.
//
// This function is useful if you would like to retrieve data from some point in the past
// until the present.
//
// Basic Usage:
//
//	ctx := createContext()
//	apiConfig := awn.CreateApiConfig(apiKey, appKey)
//	resp, err := GetHistoricalData(ctx, apiConfig)
func GetHistoricalData(ctx context.Context, funcData FunctionData, url string, version string) ([]DeviceDataResponse, error) {
	var deviceResponse []DeviceDataResponse

	for i := funcData.Epoch; i <= time.Now().UnixMilli(); i += epochIncrement24h {
		funcData.Epoch = i

		resp, err := getDeviceData(ctx, funcData, url, version)
		_ = CheckReturn(err, "unable to get device data", "warning")

		deviceResponse = append(deviceResponse, resp)
	}

	return deviceResponse, nil
}

// GetHistoricalDataAsync is a public function that takes a context object, a FunctionData
// object, the URL of the Ambient Weather Network API, the version route of the API and a
// WaitGroup object as inputs. It will return a channel of DeviceDataResponse
// objects and an error status.
//
// Basic Usage:
//
//	ctx := createContext()
//	outChannel, err := awn.GetHistoricalDataAsync(ctx, functionData, *sync.WaitGroup)
func GetHistoricalDataAsync(
	ctx context.Context,
	funcData FunctionData,
	url string,
	version string,
	w *sync.WaitGroup) (<-chan DeviceDataResponse, error) {
	defer w.Done()

	out := make(chan DeviceDataResponse)

	go func() {
		defer close(out)

		for i := funcData.Epoch; i <= time.Now().UnixMilli(); i += epochIncrement24h {
			funcData.Epoch = i

			resp, err := getDeviceData(ctx, funcData, url, version)
			_ = CheckReturn(err, "unable to get device data", "warning")

			out <- resp
		}
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
