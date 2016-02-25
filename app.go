package web

import (
	"appengine"
	"appengine/delay"
	"github.com/gorilla/mux"
	"net/http"
	"utils"
	"utils/datastoreutil"
	"utils/jsonutil"
	"utils/fileutil"
)

func init() {
	r := mux.NewRouter()

	r.HandleFunc("/robots.txt", fileutil.WriteFile("robots.txt", "text/plain")).Methods("GET")

	r.HandleFunc("/get", handleGet).Methods("GET")
	r.HandleFunc("/post", handlePost).Methods("POST")
	r.HandleFunc("/put", handle("put")).Methods("PUT")
	r.HandleFunc("/delete", handleDelete).Methods("DELETE")

	r.HandleFunc("/delay", handleDelay).Methods("GET")

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

func handlePost(w http.ResponseWriter, r *http.Request) {
	var o PostData
	jsonutil.ReadJson(w, r, &o)
	ctx := appengine.NewContext(r)
	if err := datastoreutil.PutObject(ctx, &o); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	jsonutil.WriteJson(w, o)
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	var o PostData
	ctx := appengine.NewContext(r)
	if err := datastoreutil.GetObject(ctx, &o); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	jsonutil.WriteJson(w, o)
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	var o PostData
	ctx := appengine.NewContext(r)
	if err := datastoreutil.DeleteObject(ctx, &o); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	jsonutil.WriteJson(w, o)
}

func handle(s string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		v := struct{ Method string }{s}
		jsonutil.WriteJson(w, v)
	}
}

func delayedFunction(ctx appengine.Context, data string) {
	ctx.Infof("delayedFunction called with data: %s", data)
}

func handleDelay(w http.ResponseWriter, r *http.Request) {
	var later = delay.Func("key", delayedFunction)
	later.Call(appengine.NewContext(r), "the data")
}
