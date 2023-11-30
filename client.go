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
func ConvertTimeToEpoch(tte string) (int64, error) {
	ok, err := YearMonthDay(tte).verify() //nolint:varnamelen
	if err != nil {
		log.Printf("unable to verify date")
		err = fmt.Errorf("unable to verify date: %w", err)
		return 0, err
	}

	if !ok {
		log.Fatalf("invalid date format, %v should be YYYY-MM-DD", tte)
	}

	parsed, err := time.Parse(time.DateOnly, tte)
	if err != nil {
		log.Printf("unable to parse time")
		err = fmt.Errorf("unable to parse time: %w", err)
		return 0, err
	}

	return parsed.UnixMilli(), nil
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
	// todo: check for a valid client before returning

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
	if err != nil {
		log.Printf("unable to create client")
		wrappedErr := fmt.Errorf("unable to create client: %w", err)
		return nil, wrappedErr
	}

	client.R().SetQueryParams(map[string]string{
		"apiKey":         funcData.API,
		"applicationKey": funcData.App,
	})

	deviceData := new(AmbientDevice)

	_, err = client.R().SetResult(deviceData).Get(devicesEndpoint)
	if err != nil {
		log.Printf("unable to get data from devicesEndpoint")
		wrappedErr := fmt.Errorf("unable to get data from devicesEndpoint: %w", err)
		return nil, wrappedErr
	}

	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return nil, errors.New("context timeout exceeded")
	}

	return deviceData, nil
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
	if err != nil {
		log.Printf("unable to create client")
		return DeviceDataResponse{}, err
	}

	client.R().SetQueryParams(map[string]string{
		"apiKey":         funcData.API,
		"applicationKey": funcData.App,
		"endDate":        strconv.FormatInt(funcData.Epoch, 10),
		"limit":          strconv.Itoa(funcData.Limit),
	})

	deviceData := new(DeviceDataResponse)

	_, err = client.R().
		SetPathParams(map[string]string{
			"devicesEndpoint": devicesEndpoint,
			"macAddress":      funcData.Mac,
		}).
		SetResult(deviceData).
		Get("{devicesEndpoint}/{macAddress}")
	if err != nil {
		log.Printf("unable to get data from devicesEndpoint")
		wrappedErr := fmt.Errorf("unable to get data from devicesEndpoint: %w", err)
		return DeviceDataResponse{}, wrappedErr
	}

	// todo: check response for errors passed through resp

	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return DeviceDataResponse{}, ErrContextTimeoutExceeded //nolint:exhaustruct
	}

	return *deviceData, nil
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
func GetHistoricalData(
	ctx context.Context,
	funcData FunctionData,
	url string,
	version string) ([]DeviceDataResponse, error) {
	var deviceResponse []DeviceDataResponse

	for i := funcData.Epoch; i <= time.Now().UnixMilli(); i += epochIncrement24h {
		funcData.Epoch = i

		resp, err := getDeviceData(ctx, funcData, url, version)
		if err != nil {
			log.Printf("unable to get device data")
			wrappedErr := fmt.Errorf("unable to get device data: %w", err)
			return nil, wrappedErr
		}

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
			if err != nil {
				log.Printf("unable to get device data: %v", err)
				break
			}

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
		value := GetEnvVar(vars[v])
		if value == "" {
			log.Printf("environment variable %v is empty or not set", vars[v])
		}
		envVars[vars[v]] = value
	}

	return envVars
}

// GetEnvVar is a public function attempts to fetch an environment variable. If that
// environment variable is not found, it will return an empty string.
//
// Basic Usage:
//
//	environmentVariable := GetEnvVar("ENV_VAR_1", "fallback")
func GetEnvVar(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = ""
	}

	return value
}
