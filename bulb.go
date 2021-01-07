package tradfri

import (
	"fmt"
	"strings"

	"github.com/dustin/go-coap"
)

// Bulb contains the data structure that the Gateway returns when a bulb
// device is requested.
type Bulb struct {
	RawJSON string
	On      bool
	Info    struct {
		Manufacturer          string `json:"0"`
		Model                 string `json:"1"`
		Serial                string `json:"2"`
		FirmwareVersion       string `json:"3"`
		AvailablePowerSources int    `json:"6"`
		BatteryLevel          int    `json:"9"`
	} `json:"3"`
	LightControl []struct {
		Color    string `json:"5706"`
		ColorHue int    `json:"5707"`
		ColorSat int    `json:"5708"`
		ColorX   int    `json:"5709"`
		ColorY   int    `json:"5710"`
		Power    int    `json:"5850"`
		Dim      int    `json:"5851"`
		Mireds   int    `json:"5711"`
		Duration int    `json:"5712"`
	} `json:"3311"`
	ApplicationType int    `json:"5750"`
	Name            string `json:"9001"`
	CreatedAt       int    `json:"9002"`
	ID              int    `json:"9003"`
	ReachableState  int    `json:"9019"`
	LastSeen        int    `json:"9020"`
	OTAUpdateState  int    `json:"9054"`
}

// GetBulb gets the bulb with the given ID from the gateway, if the requested
// device isn't a bulb or if it doesn't exist an error is returned. Otherwise
// a filled Bulb object is returned.
func (gw *Gateway) GetBulb(id int) (*Bulb, error) {
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

	// Decode the payload into a Bulb object.
	var b Bulb
	if err = jsonDecode(resp.Payload, &b); err != nil {
		return nil, err
	}

	// All bulbs should begin with the "TRADFRI bulb"-string, if not we
	// have gotten a device that isn't a bulb, so return an error.
	if !strings.HasPrefix(b.Info.Model, "TRADFRI bulb") {
		return nil, fmt.Errorf("device %d is a '%s', not a bulb", id, b.Info.Model)
	}

	// Store the raw JSON payload in the bulb.
	b.RawJSON = string(resp.Payload)

	// Set the convenient on property of the Bulb.
	for _, lc := range b.LightControl {
		if lc.Power == 1 {
			b.On = true
			break
		}
	}

	return &b, nil
}
