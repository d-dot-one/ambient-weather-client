package awn

import (
	"encoding/json"
	"fmt"
	"time"
)

// FunctionData is a struct that is used to pass data basic API call parameters more
// easily. It contains API (API key), App (Application key), Epoch (Unix epoch time in
// milliseconds), Limit (maximum number of records to return in a single API call), and
// Mac (MAC address of the weather station.
type FunctionData struct {
	API   string `json:"api"`
	App   string `json:"app"`
	Epoch int64  `json:"epoch"`
	Limit int    `json:"limit"`
	Mac   string `json:"mac"`
}

// String is a helper function to print the FunctionData struct as a string.
func (f FunctionData) String() string {
	r, _ := json.Marshal(f) //nolint:all

	return fmt.Sprint(string(r))
}

// NewFunctionData creates a new FunctionData object with some default values and return
// it to the caller as a pointer.
func NewFunctionData() *FunctionData {
	return &FunctionData{
		API:   "",
		App:   "",
		Epoch: 0,
		Limit: 1,
		Mac:   "",
	}
}

// DeviceDataResponse is used to marshal/unmarshal the response from the
// devices/macAddress endpoint.
type DeviceDataResponse []struct {
	Baromabsin        float64   `json:"baromabsin"`
	Baromrelin        float64   `json:"baromrelin"`
	BattLightning     int       `json:"batt_lightning"`
	Dailyrainin       float64   `json:"dailyrainin"`
	Date              time.Time `json:"date"`
	Dateutc           int64     `json:"dateutc"`
	DewPoint          float64   `json:"dewPoint"`
	DewPointin        float64   `json:"dewPointin"`
	Eventrainin       float64   `json:"eventrainin"`
	FeelsLike         float64   `json:"feelsLike"`
	FeelsLikein       float64   `json:"feelsLikein"`
	Hourlyrainin      float64   `json:"hourlyrainin"`
	Humidity          int       `json:"humidity"`
	Humidityin        int       `json:"humidityin"`
	LastRain          time.Time `json:"lastRain"`
	LightningDay      int       `json:"lightning_day"`
	LightningDistance float64   `json:"lightning_distance"`
	LightningHour     int       `json:"lightning_hour"`
	LightningTime     int64     `json:"lightning_time"`
	Maxdailygust      float64   `json:"maxdailygust"`
	Monthlyrainin     float64   `json:"monthlyrainin"`
	Solarradiation    float64   `json:"solarradiation"`
	Tempf             float64   `json:"tempf"`
	Tempinf           float64   `json:"tempinf"`
	Tz                string    `json:"tz"`
	Uv                int       `json:"uv"`
	Weeklyrainin      float64   `json:"weeklyrainin"`
	Winddir           int       `json:"winddir"`
	WinddirAvg10M     int       `json:"winddir_avg10m"`
	Windgustmph       float64   `json:"windgustmph"`
	WindspdmphAvg10M  float64   `json:"windspdmph_avg10m"`
	Windspeedmph      float64   `json:"windspeedmph"`
	Yearlyrainin      float64   `json:"yearlyrainin"`
}

// String is a helper function to print the DeviceDataResponse struct as a string.
func (d DeviceDataResponse) String() string {
	r, _ := json.Marshal(d)

	return fmt.Sprint(string(r))
}

// DeviceData is used to marshal/unmarshal the response from the
// 'devices' API endpoint. This should be removed, since this data is
// not captured. It's only possible use is for a quasi-real-time data pull.
type DeviceData struct {
	Baromabsin        float64   `json:"baromabsin"`
	Baromrelin        float64   `json:"baromrelin"`
	BattLightning     int       `json:"batt_lightning"`
	Dailyrainin       int       `json:"dailyrainin"`
	Date              time.Time `json:"date"`
	Dateutc           int64     `json:"dateutc"`
	DewPoint          float64   `json:"dewPoint"`
	DewPointin        float64   `json:"dewPointin"`
	Eventrainin       int       `json:"eventrainin"`
	FeelsLike         float64   `json:"feelsLike"`
	FeelsLikein       float64   `json:"feelsLikein"`
	Hourlyrainin      int       `json:"hourlyrainin"`
	Humidity          int       `json:"humidity"`
	Humidityin        int       `json:"humidityin"`
	LastRain          time.Time `json:"lastRain"`
	LightningDay      int       `json:"lightning_day"`
	LightningDistance float64   `json:"lightning_distance"`
	LightningHour     int       `json:"lightning_hour"`
	LightningTime     int64     `json:"lightning_time"`
	Maxdailygust      float64   `json:"maxdailygust"`
	Monthlyrainin     float64   `json:"monthlyrainin"`
	Solarradiation    float64   `json:"solarradiation"`
	Tempf             float64   `json:"tempf"`
	Tempinf           float64   `json:"tempinf"`
	Tz                string    `json:"tz"`
	Uv                int       `json:"uv"`
	Weeklyrainin      float64   `json:"weeklyrainin"`
	Winddir           int       `json:"winddir"`
	WinddirAvg10M     int       `json:"winddir_avg10m"`
	Windgustmph       float64   `json:"windgustmph"`
	WindspdmphAvg10M  float64   `json:"windspdmph_avg10m"`
	Windspeedmph      float64   `json:"windspeedmph"`
	Yearlyrainin      float64   `json:"yearlyrainin"`
}

// This info struct will likely be deleted at some point in the near future since it is
// never used. It is part of the AmbientDevice and coords structs.
type geo struct {
	Coordinates []float64 `json:"coordinates"`
	Type        string    `json:"type"`
}

// This info struct will likely be deleted at some point in the near future since it is
// never used. It is part of the AmbientDevice and coords structs.
type specificCoords struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

// This info struct will likely be deleted at some point in the near future since it is
// never used. It is part of the AmbientDevice and info structs.
type coords struct {
	Address   string         `json:"address"`
	Coords    specificCoords `json:"coords"`
	Elevation float64        `json:"elevation"`
	Geo       geo            `json:"geo"`
	Location  string         `json:"location"`
}

// This info struct will likely be deleted at some point in the near future since it is
// never used. It is part of the AmbientDevice struct.
type info struct {
	Coords coords `json:"coords"`
	Name   string `json:"name"`
}

// AmbientDevice is a struct that is used in the marshal/unmarshal JSON. This structure
// is not fully required, since all we use is the MacAddress field. The rest of the data
// is thrown away.
type AmbientDevice []struct {
	Info       info       `json:"info"`
	LastData   DeviceData `json:"DeviceData"`
	MacAddress string     `json:"macAddress"`
}

// String is a helper function to print the AmbientDevice struct as a string.
func (a AmbientDevice) String() string {
	r, err := json.Marshal(a)
	CheckReturn(err, "unable to marshall json from AmbientDevice", "warning")

	return fmt.Sprint(string(r))
}
