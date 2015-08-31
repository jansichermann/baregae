package utils

import (
	"encoding/json"
	"net/http"
)

func WriteJson(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if err := enc.Encode(v); err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func ReadJson(w http.ResponseWriter, r *http.Request, o interface{}) {
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(o); err != nil {
		http.Error(w, err.Error(), 500)
	}
}
