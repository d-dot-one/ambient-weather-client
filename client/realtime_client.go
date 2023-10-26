// this is obviously a work in progress.

package client

const (
	// apiVersionRealtime is a string and describes the version of the real-time API.
	apiVersionRealtime = "/api=1"

	// baseUrlRealtime The base URL for the Ambient Weather real-time API as a string.
	baseURLRealtime = "wss://rt2.ambientweather.net"
)

// GetRealtimeData is a public function that will connect to the Ambient Weather real-time
// weather API via Websockets and fetch live data.
func GetRealtimeData() (string, error) {
	/*
		https://ambientweather.docs.apiary.io/#reference/ambient-realtime-api

		setup: https://rt2.ambientweather.net/?api=1&applicationKey=AppKey
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		c, _, err := websocket.Dial(ctx, "ws://localhost:8080", nil)
		if err != nil {
			// ...
		}
		defer c.CloseNow()

		err = wsjson.Write(ctx, c, "hi")
		if err != nil {
			// ...
		}

		c.Close(websocket.StatusNormalClosure, "")

	*/
	_ := baseURLRealtime + apiVersionRealtime
	return "", nil
}
