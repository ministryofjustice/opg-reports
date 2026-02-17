package respond

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

// AsJSON writes data as pure json, not html / template stack involved
func AsJSON(ctx context.Context, request *http.Request, writer http.ResponseWriter, data any) {
	var (
		err        error
		dataBytes  []byte
		buffer     = new(bytes.Buffer)
		buffWriter = bufio.NewWriter(buffer)
	)

	// convert the data struct into jsonified bytes
	dataBytes, err = json.MarshalIndent(data, "", "  ")

	if err != nil {
		buffWriter.WriteString(err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Header().Set("Content-Type", "text/html")
	} else {
		buffWriter.Write(dataBytes)
		writer.WriteHeader(http.StatusOK)
		writer.Header().Set("Content-Type", "application/json")
	}
	buffWriter.Flush()
	writer.Write(buffer.Bytes())
}
