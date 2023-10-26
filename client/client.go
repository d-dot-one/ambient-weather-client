// Package client is a client that can access the Ambient Weather network API and
// return device and weather data.
package client

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"gopkg.in/resty.v1"
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

// The ConvertTimeToEpoch help function can convert any Go time.Time object to a Unix epoch time in milliseconds.
// func ConvertTimeToEpoch(t time.Time) int64 {
//	return t.UnixMilli()
//}

// The createAwnClient function is used to create a new resty-based API client. This client
// supports retries and can be placed into debug mode when needed. By default, it will
// also set the accept content type to JSON. Finally, it returns a pointer to the client.
func createAwnClient() *resty.Client {
	client := resty.New().
		SetRetryCount(retryCount).
		SetRetryWaitTime(retryMinWaitTimeSeconds * time.Second).
		SetRetryMaxWaitTime(retryMaxWaitTimeSeconds * time.Second).
		SetHostURL(baseURL + apiVersion).
		SetTimeout(defaultCtxTimeout * time.Second).
		SetDebug(debugMode).
		AddRetryCondition(
			func(r *resty.Response) (bool, error) {
				return r.StatusCode() == http.StatusRequestTimeout ||
					r.StatusCode() >= http.StatusInternalServerError ||
					r.StatusCode() == http.StatusTooManyRequests, error(nil)
			})

	client.SetHeader("Accept", "application/json")

	return client
}

// CreateAPIConfig is a helper function that is used to create the FunctionData struct,
// which is passed to the data gathering functions.
func CreateAPIConfig(api string, app string) FunctionData {
	functionData := FunctionData{
		Api: api,
		App: app,
		Ct:  createAwnClient(),
	}

	return functionData
}

// The GetDevices function takes a client, sets the appropriate query parameters for
// authentication, makes the request to the devicesEndpoint endpoint and marshals the
// response data into a pointer to an AmbientDevice object, which is returned along with
// any error messages.
func GetDevices(ctx context.Context, funcData FunctionData) (AmbientDevice, error) {
	funcData.Ct.SetQueryParams(map[string]string{
		"apiKey":         funcData.Api,
		"applicationKey": funcData.App,
	})

	deviceData := &AmbientDevice{}

	_, err := funcData.Ct.R().SetResult(deviceData).Get(devicesEndpoint)
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
	funcData.Ct.SetQueryParams(map[string]string{
		"apiKey":         funcData.Api,
		"applicationKey": funcData.App,
		"endDate":        strconv.FormatInt(funcData.Epoch, 10),
		"limit":          strconv.Itoa(funcData.Limit),
	})

	deviceData := &DeviceDataResponse{}

	_, err := funcData.Ct.R().
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

// The GetHistoricalData function takes a FunctionData object as input and returns a and
// will return a list of client.DeviceDataResponse object.
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
