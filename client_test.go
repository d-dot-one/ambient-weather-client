package awn

import (
	"context"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestConvertTimeToEpoch(t *testing.T) {
	type args struct {
		t string
	}
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
			i, err := ConvertTimeToEpoch(tt.t)
			err = CheckReturn(err, "Error converting time to epoch", "warning")
			if err != nil {
				t.Errorf("CheckReturn() error = %v", err)
			}
			if got, _ := ConvertTimeToEpoch(tt.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertTimeToEpoch() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestCheckReturn(t *testing.T) {
	t.Parallel()
	type args struct {
		e     error
		msg   string
		level LogLevelForError
	}
	tests := []struct {
		name string
		args args
	}{
		{"TestCheckReturnDebug", args{nil, "Debug log message", "debug"}},
		{"TestCheckReturnInfo", args{nil, "Info log message", "info"}},
		{"TestCheckReturnWarning", args{nil, "Warning log message", "warning"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = CheckReturn(tt.args.e, tt.args.msg, tt.args.level)
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
			if got := CreateAPIConfig(tt.args.api, tt.args.app); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateAPIConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

//	func TestGetDevices(t *testing.T) {
//		t.Parallel()
//		type args struct {
//			f FunctionData
//		}
//		tests := []struct {
//			name    string
//			args    args
//			want    AmbientDevice
//			wantErr bool
//		}{
//			// TODO: Add test cases.
//		}
//		for _, tt := range tests {
//			t.Run(tt.name, func(t *testing.T) {
//				got, err := GetDevices(tt.args.f)
//				if (err != nil) != tt.wantErr {
//					t.Errorf("GetDevices() error = %v, wantErr %v", err, tt.wantErr)
//					return
//				}
//				if !reflect.DeepEqual(got, tt.want) {
//					t.Errorf("GetDevices() got = %v, want %v", got, tt.want)
//				}
//			})
//		}
//	}
func TestGetEnvVar(t *testing.T) {
	t.Parallel()
	err := os.Setenv("TEST_ENV_VAR", "test")
	if err != nil {
		t.Errorf("unable to set test environment variable")
	}
	defer os.Unsetenv("TEST_ENV_VAR")

	type args struct {
		key      string
		fallback string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"TestGetEnvVar", args{"TEST_ENV_VAR", "fallback"}, "test"},
		{"TestGetEnvVarEmpty", args{"TEST_ENV_VAR_EMPTY", "fallback"}, "fallback"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetEnvVar(tt.args.key, tt.args.fallback); got != tt.want {
				t.Errorf("GetEnvVar() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetEnvVars(t *testing.T) {
	t.Parallel()
	err := os.Setenv("TEST_ENV_VAR", "test")
	if err != nil {
		t.Errorf("unable to set test environment variable")
	}
	defer os.Unsetenv("TEST_ENV_VAR")

	type args struct {
		vars []string
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{"TestGetEnvVars", args{[]string{"TEST_ENV_VAR", "ANOTHER_TEST_ENV_VAR"}}, map[string]string{"TEST_ENV_VAR": "test", "ANOTHER_TEST_ENV_VAR": ""}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetEnvVars(tt.args.vars); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetEnvVars() = %v, want %v", got, tt.want)
			}
		})
	}
}

//func TestGetHistoricalData(t *testing.T) {
//	t.Parallel()
//	type args struct {
//		f FunctionData
//	}
//	tests := []struct {
//		name    string
//		args    args
//		want    []DeviceDataResponse
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, err := GetHistoricalData(tt.args.f)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("GetHistoricalData() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("GetHistoricalData() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestGetHistoricalDataAsync(t *testing.T) {
//	t.Parallel()
//	type args struct {
//		f FunctionData
//		w *sync.WaitGroup
//	}
//	tests := []struct {
//		name    string
//		args    args
//		want    <-chan DeviceDataResponse
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, err := GetHistoricalDataAsync(tt.args.f, tt.args.w)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("GetHistoricalDataAsync() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("GetHistoricalDataAsync() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func Test_createAwnClient(t *testing.T) {
//	t.Parallel()
//	tests := []struct {
//		name string
//		want *resty.Client
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := createAwnClient(); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("createAwnClient() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
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
