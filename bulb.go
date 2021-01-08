package tradfri

import (
	"fmt"
)

// Bulb contains the data structure that the Gateway returns when a bulb
// device is requested.
type Bulb struct {
	JSON []byte

	Info         Info `json:"3"`
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
	d, err := gw.GetDevice(id)
	if err != nil {
		return nil, err
	}

	if !d.IsBulb() {
		return nil, fmt.Errorf("device %d is a '%s', not a bulb", id, d.Info.Model)
	}

	var ret Bulb
	if err = jsonDecode(d.JSON, &ret); err != nil {
		return nil, err
	}

	ret.JSON = d.JSON

	return &ret, nil
}
