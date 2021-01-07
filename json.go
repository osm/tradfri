package tradfri

import (
	"encoding/json"
)

// jsonEncode encodes the given input and returns a JSON encoded slice of
// bytes. We don't do any error checking in this function since it is assumed
// that the user of the function is 100 % sure that the given data is possible
// to be encoded.
func jsonEncode(p interface{}) []byte {
	r, _ := json.Marshal(p)
	return r
}

// jsonDecode decodes the given input into the empty interface.
func jsonDecode(b []byte, r interface{}) error {
	return json.Unmarshal(b, &r)
}
