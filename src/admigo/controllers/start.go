package controllers

import (
	"net/http"

	"admigo/common"

	"github.com/julienschmidt/httprouter"
)

// Index home page
func Index(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	qv := r.URL.Query()

	some := qv.Get("some")

	var info map[string]string
	if len(some) > 0 {
		info = make(map[string]string)
		info["some"] = common.Hapas(some)
	}

	setFrontContent(w, r, "home", info, nil, "home/home")
}
