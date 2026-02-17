package cnv

import "encoding/json"

// Convert takes original struct of T and by marshaling and then unmarshaling applied its
// content to destination R
func Convert[T any, R any](source T, destination R) (err error) {
	var bytes []byte
	if bytes, err = json.MarshalIndent(source, "", "  "); err == nil {
		err = json.Unmarshal(bytes, destination)
	}
	return
}
