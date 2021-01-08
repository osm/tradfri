package tradfri

import (
	"fmt"
	"strings"

	"github.com/dustin/go-coap"
)

// Info contains the common structure that all IKEA Trådfri devices share.
type Info struct {
	Manufacturer          string `json:"0"`
	Model                 string `json:"1"`
	Serial                string `json:"2"`
	FirmwareVersion       string `json:"3"`
	AvailablePowerSources int    `json:"6"`
	BatteryLevel          int    `json:"9"`
}

// Device contains the structure of an unknown IKEA Trådfri device.
type Device struct {
	JSON []byte

	Info           Info   `json:"3"`
	Name           string `json:"9001"`
	CreatedAt      int    `json:"9002"`
	ID             int    `json:"9003"`
	ReachableState int    `json:"9019"`
	LastSeen       int    `json:"9020"`
}

// IsBulb checks whether or not the device is a bulb or not.
func (d *Device) IsBulb() bool {
	return strings.HasPrefix(d.Info.Model, "TRADFRI bulb")
}

// IsControlOutlet checks whether or not the device is a control outlet or
// not.
func (d *Device) IsControlOutlet() bool {
	return strings.HasPrefix(d.Info.Model, "TRADFRI control outlet")
}

// GetDevice fetches the device with the given ID.
func (gw *Gateway) GetDevice(id int) (*Device, error) {
	// Request the given device.
	resp, err := gw.client.send(fmt.Sprintf("/15001/%d", id), &coap.Message{
		Type: coap.Confirmable,
		Code: coap.GET,
	})
	if err != nil {
		return nil, err
	}
	if resp.Code == coap.NotFound {
		return nil, ErrNotFound
	}

	// Decode the payload.
	var d Device
	if err = jsonDecode(resp.Payload, &d); err != nil {
		return nil, err
	}

	// Store the raw JSON payload in the bulb.
	d.JSON = resp.Payload

	return &d, nil
}

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
