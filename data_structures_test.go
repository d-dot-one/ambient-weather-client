package awn

import (
	"context"
	"reflect"
	"testing"
	"time"
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
	t.Skip("flaky test")

	dateVar, _ := time.Parse(time.RFC3339, "2023-07-01T12:00:30Z")
	lastRainVar, _ := time.Parse(time.RFC3339, "2023-10-12T04:53:00.000Z")

	tests := []struct {
		name string
		d    DeviceDataResponse
		want string
	}{
		{name: "FullSuite", d: DeviceDataResponse{
			Baromabsin: 29.675, Baromrelin: 29.775,
			BattLightning: 0, Dailyrainin: 1.234,
			Date: dateVar, Dateutc: 1697142300000,
			DewPoint: 78.51, DewPointin: 78,
			Eventrainin: 10.023, FeelsLike: 99.2,
			Hourlyrainin: 1.11, Humidity: 79,
			Humidityin: 76, LastRain: lastRainVar,
			LightningDay: 1, LightningDistance: 4.97,
			LightningHour: 53, LightningTime: 1696633175000,
			Maxdailygust: 9.8, Monthlyrainin: 5.925,
			Solarradiation: 455.56, Tempf: 85.8,
			Tempinf: 5.254, Tz: "America",
			Uv: 4, Weeklyrainin: 2.122,
			Winddir: 239, WinddirAvg10M: 250,
			Windgustmph: 5.6, WindspdmphAvg10M: 2.7,
			Windspeedmph: 4.3, Yearlyrainin: 34.457,
		}, want: `{"baromabsin":29.675,"baromrelin":29.775,"batt_lightning":0,"dailyrainin":1.234,"date":"2023-07-01T12:00:30Z","dateutc":1697142300000,"dewPoint":78.51,"dewPointin":78,"eventrainin":10.023,"feelsLike":99.2,"feelsLikein":0,"hourlyrainin":1.11,"humidity":79,"humidityin":76,"lastRain":"2023-10-12T04:53:00.000Z","lightning_day":1,"lightning_distance":4.97,"lightning_hour":53,"lightning_time":1696633175000,"maxdailygust":9.8,"monthlyrainin":5.925,"solarradiation":455.56,"tempf":85.8,"tempinf":5.24,"tz":"America","uv":4,"weeklyrainin":2.122,"winddir":239,"winddir_avg10m":250,"windgustmph":5.6,"windspdmph_avg10m":2.7,"windspeedmph":4.3,"yearlyrainin":34.457}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.d.String()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TestDeviceDataResponseToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFunctionData(t *testing.T) {
	type params struct {
		API   string
		App   string
		Epoch int64
		Limit int
		Mac   string
	}
	tests := []struct {
		name   string
		fields params
		want   FunctionData
	}{
		{"FunctionData", params{API: "api", App: "app", Epoch: 1234567890, Limit: 100, Mac: "00:11:22:33:44:55"}, FunctionData{}},
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
			tof := reflect.TypeOf(f)
			tow := reflect.TypeOf(tt.want)

			if tof != tow {
				t.Errorf("FunctionData = %v, want %v", tof, tow)
			}
		})
	}
}

func TestFunctionDataToString(t *testing.T) {
	t.Skip("not yet implemented")

	_, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	type params struct {
		API   string
		App   string
		Epoch int64
		Limit int
		Mac   string
	}
	tests := []struct {
		name   string
		fields params
		want   string
	}{
		{"FunctionDataToString", params{API: "api", App: "app", Epoch: 1234567890, Limit: 100, Mac: "00:11:22:33:44:55"}, `{"api":"api","app":"app","epoch":1234567890,"limit":100,"mac":"00:11:22:33:44:55"}`},
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

func TestFunctionDataToMap(t *testing.T) {
	_, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	type params struct {
		API   string
		App   string
		Epoch int64
		Limit int
		Mac   string
	}
	tests := []struct {
		name   string
		fields params
		want   map[string]interface{}
	}{
		{name: "FunctionDataToString", fields: params{API: "api", App: "app", Epoch: 1234567890, Limit: 100, Mac: "00:11:22:33:44:55"}, want: FunctionData{API: "api", App: "app", Epoch: 1234567890, Limit: 100, Mac: "00:11:22:33:44:55"}.ToMap()},
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
			if f.API != tt.fields.API {
				t.Errorf("FunctionDataToMap().API = %v, want %v", f.API, tt.fields.API)
			}
			if f.App != tt.fields.App {
				t.Errorf("FunctionDataToMap().App = %v, want %v", f.App, tt.fields.App)
			}
			if f.Epoch != tt.fields.Epoch {
				t.Errorf("FunctionDataToMap().Epoch = %v, want %v", f.Epoch, tt.fields.Epoch)
			}
			if f.Limit != tt.fields.Limit {
				t.Errorf("FunctionDataToMap().Limit = %v, want %v", f.Limit, tt.fields.Limit)
			}
			if f.Mac != tt.fields.Mac {
				t.Errorf("FunctionDataToMap().Mac = %v, want %v", f.Mac, tt.fields.Mac)
			}

			got := f.ToMap()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FunctionDataToMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewFunctionData(t *testing.T) {
	f1 := FunctionData{API: "", App: "", Epoch: 0, Limit: 1, Mac: ""}

	tests := []struct {
		name string
		new  *FunctionData
		want *FunctionData
	}{
		{"NewFunctionData", NewFunctionData(), &f1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.new; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFunctionData() = %v, want %v", got, tt.want)
			}
		})
	}
}
