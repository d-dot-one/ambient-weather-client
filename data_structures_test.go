package awn

import (
	"context"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
)

func TestAmbientDeviceToString(t *testing.T) {
	tests := []struct {
		name string
		a    AmbientDevice
		want string
	}{
		{"TestAmbientDeviceMacString", AmbientDevice{MacAddress: "00:11:22:33:44:55"}, `{"info":{"coords":{"address":"","coords":{"lat":0,"lon":0},"elevation":0,"geo":{"coordinates":null,"type":""},"location":""},"name":""},"DeviceData":{"baromabsin":0,"baromrelin":0,"batt_lightning":0,"dailyrainin":0,"date":"0001-01-01T00:00:00Z","dateutc":0,"dewPoint":0,"dewPointin":0,"eventrainin":0,"feelsLike":0,"feelsLikein":0,"hourlyrainin":0,"humidity":0,"humidityin":0,"lastRain":"0001-01-01T00:00:00Z","lightning_day":0,"lightning_distance":0,"lightning_hour":0,"lightning_time":0,"maxdailygust":0,"monthlyrainin":0,"solarradiation":0,"tempf":0,"tempinf":0,"tz":"","uv":0,"weeklyrainin":0,"winddir":0,"winddir_avg10m":0,"windgustmph":0,"windspdmph_avg10m":0,"windspeedmph":0,"yearlyrainin":0},"macAddress":"00:11:22:33:44:55"}`},
		{name: "TestAmbientDeviceInfoString", a: AmbientDevice{Info: info{Coords: coords{Address: "123 Main", Location: "Anywhere, USA"}}}, want: `{"info":{"coords":{"address":"123 Main","coords":{"lat":0,"lon":0},"elevation":0,"geo":{"coordinates":null,"type":""},"location":"Anywhere, USA"},"name":""},"DeviceData":{"baromabsin":0,"baromrelin":0,"batt_lightning":0,"dailyrainin":0,"date":"0001-01-01T00:00:00Z","dateutc":0,"dewPoint":0,"dewPointin":0,"eventrainin":0,"feelsLike":0,"feelsLikein":0,"hourlyrainin":0,"humidity":0,"humidityin":0,"lastRain":"0001-01-01T00:00:00Z","lightning_day":0,"lightning_distance":0,"lightning_hour":0,"lightning_time":0,"maxdailygust":0,"monthlyrainin":0,"solarradiation":0,"tempf":0,"tempinf":0,"tz":"","uv":0,"weeklyrainin":0,"winddir":0,"winddir_avg10m":0,"windgustmph":0,"windspdmph_avg10m":0,"windspeedmph":0,"yearlyrainin":0},"macAddress":""}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeviceDataResponseToString(t *testing.T) {
	dateVar, _ := time.Parse(time.RFC3339, "2023-07-01T12:00:30Z")
	lastRainVar, _ := time.Parse(time.RFC3339, "2023-07-04T00:30:45Z")

	tests := []struct {
		name string
		d    DeviceDataResponse
		want string
	}{
		{name: "FullSuite", d: DeviceDataResponse{
			Baromabsin:        1004.4,
			Baromrelin:        999.9,
			BattLightning:     1,
			Dailyrainin:       1.23,
			Date:              dateVar,
			Dateutc:           22220101,
			DewPoint:          174.3,
			DewPointin:        74.3,
			Eventrainin:       0.0,
			FeelsLike:         70.0,
			Hourlyrainin:      1.11,
			Humidity:          99,
			Humidityin:        88,
			LastRain:          lastRainVar,
			LightningDay:      1,
			LightningDistance: 5.254,
			LightningHour:     53,
			LightningTime:     170000000,
			Maxdailygust:      5.254,
			Monthlyrainin:     5.254,
			Solarradiation:    5.254,
			Tempf:             5.254,
			Tempinf:           5.254,
			Tz:                "GMT",
			Uv:                33,
			Weeklyrainin:      5.254,
			Winddir:           353,
			WinddirAvg10M:     53,
			Windgustmph:       5.254,
			WindspdmphAvg10M:  5.254,
			Windspeedmph:      5.254,
			Yearlyrainin:      5.254,
		}, want: `{"baromabsin":1004.4,"baromrelin":999.9,"batt_lightning":1,"dailyrainin":1.23,"date":"2023-07-01T12:00:30Z","dateutc":22220101,"dewPoint":174.3,"dewPointin":74.3,"eventrainin":0.0,"feelsLike":0.0,"feelsLikein":0.0,"hourlyrainin":1.11,"humidity":99,"humidityin":88,"lastRain":"2023-07-04T00:30:45Z","lightning_day":1,"lightning_distance":5.24,"lightning_hour":53,"lightning_time":170000000,"maxdailygust":5.254,"monthlyrainin":5.254,"solarradiation":5.254,"tempf":5.254,"tempinf":5.254,"tz":"GMT","uv":33,"weeklyrainin":5.254,"winddir":353,"winddir_avg10m":53,"windgustmph":5.254,"windspdmph_avg10m":5.254,"windspeedmph":5.254,"yearlyrainin":5.254}`},
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

	_, err := CreateAwnClient("https://rt.ambientweather.net", "/v1")
	_ = CheckReturn(err, "unable to create client", "warning")

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
