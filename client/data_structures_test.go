package client

import (
	"context"
	"testing"
	"time"

	"gopkg.in/resty.v1"
)

func TestAmbientDevice_String(t *testing.T) {
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
	tests := []struct {
		name string
		d    DeviceDataResponse
		want string
	}{
		// TODO: Add test cases.
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	epoch := time.Now().UnixMilli()

	client := createAwnClient()

	type fields struct {
		Api   string
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
		{name: "FunctionDataString()", fields: {Api: "api", App: "app", Ct: createAwnClient(), Cx: ctx, Epoch: epoch, Limit: 100, Mac: "00:11:22:33:44:55"}, want: {}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := FunctionData{
				Api:   tt.fields.Api,
				App:   tt.fields.App,
				Ct:    tt.fields.Ct,
				Cx:    tt.fields.Cx,
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
