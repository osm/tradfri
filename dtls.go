package tradfri

import (
	"net"

	"github.com/dustin/go-coap"
	"github.com/pion/dtls"
)

// dtlsClient holds the dtls connection and a message ID variable that gets
// incremented for each request that is made.
type dtlsClient struct {
	conn      *dtls.Conn
	messageID uint16
}

// newDTLSClient initializes a new connection to the given addr with the
// identity and psk and returns a pointer to the dtls.Conn object.
func newDTLSClient(addr, identity, psk string) (*dtlsClient, error) {
	conn, err := dtls.Dial(
		"udp",
		&net.UDPAddr{IP: net.ParseIP(addr), Port: 5684},
		&dtls.Config{
			PSK: func(_ []byte) ([]byte, error) {
				return []byte(psk), nil
			},
			PSKIdentityHint: []byte(identity),
			CipherSuites:    []dtls.CipherSuiteID{dtls.TLS_PSK_WITH_AES_128_CCM_8},
		},
	)
	if err != nil {
		if err.Error() == "alert: Alert LevelFatal: BadRecordMac" {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	return &dtlsClient{
		conn:      conn,
		messageID: 1,
	}, nil
}

// send sends the given message to the url.
func (d *dtlsClient) send(url string, msg *coap.Message) (*coap.Message, error) {
	// Set the path on the message.
	msg.SetPathString(url)

	// Set the MessageID in the message.
	msg.MessageID = d.messageID

	// Increment the message ID. More info on the message ID can be found
	// in the RFC: https://tools.ietf.org/html/rfc7252#section-2.1
	d.messageID++

	// Marshal the message.
	data, err := msg.MarshalBinary()
	if err != nil {
		return nil, err
	}

	// Write the message on the already established connection. If there's
	// an error we'll return it.
	_, err = d.conn.Write(data)
	if err != nil {
		return nil, err
	}

	// Allocate a buffer, 1152 bytes should be large enough accoring to
	// the RFC: https://tools.ietf.org/html/rfc7252#section-4.6
	// If there's an error we'll return it.
	buff := make([]byte, 1152)
	bytes, err := d.conn.Read(buff)
	if err != nil {
		return nil, err
	}

	// Parse the buffer, we'll only use the returned number of read bytes
	// for the message, otherwise we'll get a large zero padded Payload in
	// the parsed message.
	ret, err := coap.ParseMessage(buff[0:bytes])
	if err != nil {
		return nil, err
	}

	// Return the parsed message.
	return &ret, nil
}
