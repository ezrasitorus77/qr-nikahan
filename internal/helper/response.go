package helper

import (
	"encoding/json"
	"net/http"
	"qr-nikahan/domain"
)

func Response(w http.ResponseWriter, resp domain.Response, status int) {
	var encoder *json.Encoder

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	encoder = json.NewEncoder(w)
	encoder.Encode(resp)
}
