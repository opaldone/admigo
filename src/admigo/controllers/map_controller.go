package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"admigo/common"
	"admigo/config"

	"github.com/julienschmidt/httprouter"
)

func MapLoca(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	setFrontContent(w, r, "loca", nil, nil,
		"map/ix/loca",
		"map/ix/_mus",
		"map/ix/_logs",
		"map/ix/_ros",
	)
}

func MapWs(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	e := config.Env(false)

	cid := fmt.Sprintf("re-%s", common.RandUID())

	link := fmt.Sprintf("%s/ws/%s/0", e.Map.Ws, cid)

	js := map[string]string{
		"cid":        cid,
		"link":       link,
		"startpoint": e.Map.StartPoint,
		"routeurl":   e.Map.RouteURL,
		"routekey":   e.Map.RouteKey,
	}

	output, _ := json.Marshal(js)

	w.Write(output)
}
