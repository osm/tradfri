package tradfri

import (
	"fmt"
)

type ControlOutlet struct {
	JSON []byte

	Info Info `json:"3"`
	Plug struct {
		Power int `json:"5850"`
		Dim   int `json:"5851"`
		ID    int `json:"9003"`
	}
	ApplicationType int    `json:"5750"`
	Name            string `json:"9001"`
	CreatedAt       int    `json:"9002"`
	ID              int    `json:"9003"`
	ReachableState  int    `json:"9019"`
	LastSeen        int    `json:"9020"`
	OTAUpdateState  int    `json:"9054"`
}

// GetControlOutlet gets the control outlet with the given id.
func (gw *Gateway) GetControlOutlet(id int) (*ControlOutlet, error) {
	d, err := gw.GetDevice(id)
	if err != nil {
		return nil, err
	}

	if !d.IsControlOutlet() {
		return nil, fmt.Errorf("device %d is a '%s', not a control outlet", id, d.Info.Model)
	}

	var ret ControlOutlet
	if err = jsonDecode(d.JSON, &ret); err != nil {
		return nil, err
	}

	ret.JSON = d.JSON

	return &ret, nil
}
