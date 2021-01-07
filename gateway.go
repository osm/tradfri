package tradfri

import (
	"github.com/dustin/go-coap"
)

// Gateway.
type Gateway struct {
	addr       string
	identifier string
	psk        string
	client     *dtlsClient
}

// AuthenticateReq contains the data model that should be sent to the IKEA
// Trådfri gateway to retrieve a new identifier/PSK pair.
type AuthenticateReq struct {
	Identity string `json:"9090"`
}

// AuthenticateResp contains the data model that is returned by the IKEA
// Trådfri gateway on a successfull authentication.
type AuthenticateResp struct {
	PSK     string `json:"9091"`
	Version string `json:"9029"`
}

// Authenticate connects to the given address with the give PSK (which can be
// found underneath the IKEA TRADFRI Gateway. On a successfull connection
// we'll return a new identity and PSK which should be used for further access
// to the gateway.
func Authenticate(addr, initPSK string) (string, string, error) {
	// The IKEA TRADFRI gateway always uses the "Client_identity" when
	// initializing a new connection to the gateway.
	dc, err := newDTLSClient(addr, "Client_identity", initPSK)
	if err != nil {
		return "", "", err
	}

	// We'll generate a new UUID to be used as the identity for the PSK we
	// are about to request from the Gateway.
	identity := newUUID()

	// Send a message to the endpoint in the Gateway which is responsible
	// for generating a new PSK for the given identity.
	resp, err := dc.send("/15011/9063", &coap.Message{
		Type:    coap.Confirmable,
		Code:    coap.POST,
		Payload: jsonEncode(AuthenticateReq{identity}),
	})

	// Make sure that we didn't receive an error, if so, return it.
	if err != nil {
		return "", "", err
	}

	// If the returned Code wasn't "Created" we should return an error.
	if resp.Code != coap.Created {
		return "", "", ErrBadRequest
	}

	var r AuthenticateResp
	if err = jsonDecode(resp.Payload, &r); err != nil {
		return "", "", err
	}

	return identity, r.PSK, nil
}

// Connect initializes a new connection to the given address with the
// identifier and psk, if the connection is successful we'll return a new
// instance of the gateway, if not an error is returned.
func Connect(addr, identifier, psk string) (*Gateway, error) {
	// Establish a new connection to the gateway.
	dc, err := newDTLSClient(addr, identifier, psk)
	if err != nil {
		return nil, err
	}

	return &Gateway{
		addr:       addr,
		identifier: identifier,
		psk:        psk,
		client:     dc,
	}, nil
}
