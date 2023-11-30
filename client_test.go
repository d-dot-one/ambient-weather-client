package awn

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
)

func TestConvertTimeToEpoch(t *testing.T) {
	tests := []struct {
		name string
		t    string
		want int64
	}{
		{"Test01Jan2014ToEpoch", "2014-01-01", 1388534400000},
		{"Test15Nov2023ToEpoch", "2023-11-15", 1700006400000},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := ConvertTimeToEpoch(tt.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertTimeToEpoch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertTimeToEpochBadFormat(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		t    string
		want error
	}{
		{"TestWrongDateFormat", "11-15-2021", ErrMalformedDate},
		{"TestWrongDateFormat", "11152021", ErrMalformedDate},
		{"TestWrongDateFormat", "11-15-2021:12:42", ErrMalformedDate},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ConvertTimeToEpoch(tt.t)
			if err == nil {
				t.Errorf("ConvertTimeToEpoch() = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestCreateApiConfig(t *testing.T) {
	t.Parallel()
	_, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	fd := FunctionData{"api", "app", 0, 1, ""}

	type args struct {
		api string
		app string
	}
	tests := []struct {
		name string
		args args
		want *FunctionData
	}{
		{name: "TestCreateApiConfig", args: args{"api", "app"}, want: &fd},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CreateAPIConfig(tt.args.api, tt.args.app)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateAPIConfig() = %v, want %v", got, tt.want)
			}
			if got.API != tt.want.API {
				t.Errorf("CreateAPIConfig() API = %v, want %v", got.API, tt.want.API)
			}
			if got.App != tt.want.App {
				t.Errorf("CreateAPIConfig() App = %v, want %v", got.App, tt.want.App)
			}
		})
	}
}

func TestGetLatestData(t *testing.T) {
	t.Skip("skipping test -- flaky")

	fd := FunctionData{API: "api_key_goes_here", App: "app_key_goes_here"}
	jsonData := `{"info": {}, "DeviceData": {}, "macAddress": "00:00:00:00:00:00"}`
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	s := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Write([]byte(jsonData))
		}))
	defer s.Close()

	type Tests struct {
		name     string
		baseURL  string
		ctx      context.Context
		version  string
		response *AmbientDevice
		want     error
	}
	tests := []Tests{
		{
			name:     "basic-request",
			baseURL:  s.URL,
			ctx:      ctx,
			version:  "/v1",
			response: &AmbientDevice{},
			want:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetLatestData(ctx, fd, tt.baseURL, tt.version)
			if !reflect.DeepEqual(got, tt.response) {
				t.Errorf("GetLatestData() = %v, want %v", got, tt.response)
			}
			if err != nil {
				t.Errorf("GetLatestData() Error = %v, want %v", err, tt.want)
			}
			if !errors.Is(err, tt.want) {
				t.Errorf("GetLatestData() Error = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestGetEnvVar(t *testing.T) {
	t.Parallel()
	err := os.Setenv("TEST_ENV_VAR", "test")
	if err != nil {
		t.Errorf("unable to set test environment variable")
	}

	type args struct {
		key string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"TestGetEnvVar", args{"TEST_ENV_VAR"}, "test"},
		{"TestGetEnvVarEmpty", args{"TEST_ENV_VAR_EMPTY"}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetEnvVar(tt.args.key); got != tt.want {
				t.Errorf("GetEnvVar() = %v, want %v", got, tt.want)
			}
		})
	}

	t.Cleanup(func() {
		_ = os.Unsetenv("TEST_ENV_VAR")
	})
}

func TestGetEnvVars(t *testing.T) {
	t.Parallel()

	err := os.Setenv("TEST_ENV_VAR", "test")
	if err != nil {
		t.Errorf("unable to set test environment variable")
	}
	err = os.Setenv("ANOTHER_TEST_ENV_VAR", "another_test")
	if err != nil {
		t.Errorf("unable to set test environment variable")
	}

	type args struct {
		vars []string
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{"TestGetEnvVars", args{[]string{"TEST_ENV_VAR", "ANOTHER_TEST_ENV_VAR"}}, map[string]string{"TEST_ENV_VAR": "test", "ANOTHER_TEST_ENV_VAR": "another_test"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetEnvVars(tt.args.vars); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetEnvVars() = %v, want %v", got, tt.want)
			}
		})
	}

	t.Cleanup(func() {
		err := os.Unsetenv("TEST_ENV_VAR")
		if err != nil {
			t.Errorf("unable to unset test environment variable")
		}
		err = os.Unsetenv("ANOTHER_TEST_ENV_VAR")
		if err != nil {
			t.Errorf("unable to unset another test environment variable")
		}
	})
}

//	func TestGetHistoricalData(t *testing.T) {
//		t.Parallel()
//		type args struct {
//			f FunctionData
//		}
//		tests := []struct {
//			name    string
//			args    args
//			want    []DeviceDataResponse
//			wantErr bool
//		}{
//			// TODO: Add test cases.
//		}
//		for _, tt := range tests {
//			t.Run(tt.name, func(t *testing.T) {
//				got, err := GetHistoricalData(tt.args.f)
//				if (err != nil) != tt.wantErr {
//					t.Errorf("GetHistoricalData() error = %v, wantErr %v", err, tt.wantErr)
//					return
//				}
//				if !reflect.DeepEqual(got, tt.want) {
//					t.Errorf("GetHistoricalData() got = %v, want %v", got, tt.want)
//				}
//			})
//		}
//	}
//
//	func TestGetHistoricalDataAsync(t *testing.T) {
//		t.Parallel()
//		type args struct {
//			f FunctionData
//			w *sync.WaitGroup
//		}
//		tests := []struct {
//			name    string
//			args    args
//			want    <-chan DeviceDataResponse
//			wantErr bool
//		}{
//			// TODO: Add test cases.
//		}
//		for _, tt := range tests {
//			t.Run(tt.name, func(t *testing.T) {
//				got, err := GetHistoricalDataAsync(tt.args.f, tt.args.w)
//				if (err != nil) != tt.wantErr {
//					t.Errorf("GetHistoricalDataAsync() error = %v, wantErr %v", err, tt.wantErr)
//					return
//				}
//				if !reflect.DeepEqual(got, tt.want) {
//					t.Errorf("GetHistoricalDataAsync() got = %v, want %v", got, tt.want)
//				}
//			})
//		}
//	}

func TestCreateAwnClient(t *testing.T) {
	t.Parallel()
	header := http.Header{}
	header.Add("Accept", "application/json")

	tests := []struct {
		name string
		want *resty.Client
	}{
		{name: "TestCreateAwnClient", want: &resty.Client{
			BaseURL:                "http://127.0.0.1",
			Header:                 header,
			RetryCount:             0,
			RetryWaitTime:          retryMinWaitTimeSeconds * time.Second,
			RetryMaxWaitTime:       retryMaxWaitTimeSeconds * time.Second,
			HeaderAuthorizationKey: "Authorization",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := CreateAwnClient("http://127.0.0.1", "/")
			if got.BaseURL != tt.want.BaseURL &&
				!reflect.DeepEqual(got.Header, tt.want.Header) &&
				got.RetryCount != tt.want.RetryCount &&
				got.RetryWaitTime != tt.want.RetryWaitTime &&
				got.RetryMaxWaitTime != tt.want.RetryMaxWaitTime &&
				got.HeaderAuthorizationKey != tt.want.HeaderAuthorizationKey {
				t.Errorf("createAwnClient() got = %v, want %v", got, tt.want)
			}
		})
	}
}

//func Test_getDeviceData(t *testing.T) {
//	t.Parallel()
//	type args struct {
//		f FunctionData
//	}
//	tests := []struct {
//		name    string
//		args    args
//		want    DeviceDataResponse
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, err := getDeviceData(tt.args.f)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("getDeviceData() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("getDeviceData() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
