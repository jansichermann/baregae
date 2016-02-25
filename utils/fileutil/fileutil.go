package fileutil

import (
	"net/http"
	"os"
)

func WriteFile(filename, mimeType string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if f, err := os.Open(filename); err != nil {
			http.Error(w, err.Error(), 500)
		} else {
			defer f.Close()
			w.Header().Set("Content-Type", mimeType)
			b := make([]byte, 1024)
			for n, err := f.Read(b); err == nil && n > 0; n, err = f.Read(b) {
				w.Write(b[:n])
			}
		}
	}
}
