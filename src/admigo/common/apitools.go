package common

import (
	"encoding/json"
	"io"
	"net/http"

	"admigo/applog"
)

type AjaxAnswer struct {
	Cont string `json:"cont,omitempty"`
}

type TreeRequest struct {
	Lev   string    `json:"lev,omitempty"`
	Par   string    `json:"par,omitempty"`
	Prid  string    `json:"prid,omitempty"`
	Fid   string    `json:"fid,omitempty"`
	Pk    string    `json:"pk,omitempty"`
	Root  string    `json:"root,omitempty"`
	Paths *[]string `json:"paths,omitempty"`
}

func doJsonFromBody(body io.Reader, v interface{}) {
	decoder := json.NewDecoder(body)
	if err := decoder.Decode(v); err != nil {
		applog.Danger("Cannot parse form", err)
	}
}

func ReqTreePars(r *http.Request) TreeRequest {
	tp := TreeRequest{}
	doJsonFromBody(r.Body, &tp)

	return tp
}
