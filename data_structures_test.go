package awn

import (
	"context"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
)

func TestAmbientDevice_String(t *testing.T) {
	t.Skip("not yet implemented")

	tests := []struct {
		name string
		a    AmbientDevice
		want string
	}{
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeviceDataResponse_String(t *testing.T) {
	t.Skip("not yet implemented")

	tests := []struct {
		name string
		d    DeviceDataResponse
		want string
	}{
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFunctionData_String(t *testing.T) {
	t.Skip("not yet implemented")

	_, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	//epoch := time.Now().UnixMilli()

	_, err := createAwnClient()
	CheckReturn(err, "unable to create client", "warning")

	type fields struct {
		API   string
		App   string
		Ct    *resty.Client
		Cx    context.Context
		Epoch int64
		Limit int
		Mac   string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		//{name: "FunctionDataString()", fields: {API: "api", App: "app", Epoch: epoch, Limit: 100, Mac: "00:11:22:33:44:55"}, want: {}},
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := FunctionData{
				API:   tt.fields.API,
				App:   tt.fields.App,
				Epoch: tt.fields.Epoch,
				Limit: tt.fields.Limit,
				Mac:   tt.fields.Mac,
			}
			if got := f.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
