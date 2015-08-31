package web

import (
	"appengine"
	"github.com/gorilla/mux"
	"net/http"
	"utils"
)

func init() {
	r := mux.NewRouter()

	r.HandleFunc("/robots.txt", utils.WriteFile("robots.txt", "text/plain")).Methods("GET")

	r.HandleFunc("/get", handleGet()).Methods("GET")
	r.HandleFunc("/post", handlePost()).Methods("POST")
	r.HandleFunc("/put", handle("put")).Methods("PUT")
	r.HandleFunc("/delete", handle("delete")).Methods("DELETE")

	http.Handle("/", r)
}

type PostData struct {
	Message string `json:"message"`
}

func (p *PostData) ObjectId() string {
	return "foobar"
}

func (p *PostData) EntityType() string {
	return "PostData"
}

func handlePost() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var o PostData
		utils.ReadJson(w, r, &o)
		ctx := appengine.NewContext(r)
		if err := utils.PutObject(ctx, &o); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		utils.WriteJson(w, o)
	}
}

func handleGet() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var o PostData
		ctx := appengine.NewContext(r)
		if err := utils.GetObject(ctx, &o); err != nil {
			http.Error(w, err.Error(), 500)
		}
		utils.WriteJson(w, o)
	}
}

func handle(s string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		v := struct{ Method string }{s}
		utils.WriteJson(w, v)
	}
}
