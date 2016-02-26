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
			p := make([]byte, 1024)
			for {
				n, _ := f.Read(p)
				if n == 0 {
					break
				}
				w.Write(p[:n])
			}
		}
	}
}
