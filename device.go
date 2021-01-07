package tradfri

import (
	"github.com/dustin/go-coap"
)

// ListDevices retrieves a list of all the device ids that has been connected
// to the gateway.
func (gw *Gateway) ListDeviceIDs() ([]int, error) {
	// Get the /15001-route on the gateway, which corresponds to a route
	// that lists device ids for all the connected entities.
	resp, err := gw.client.send("/15001", &coap.Message{
		Type: coap.Confirmable,
		Code: coap.GET,
	})
	if err != nil {
		return nil, err
	}

	// The response is an array of integers, so we do not need to specifiy
	// a special struct for this route.
	var r []int
	if err = jsonDecode(resp.Payload, &r); err != nil {
		return nil, err
	}

	return r, nil
}
