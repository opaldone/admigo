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

func doJSONFromBody(body io.Reader, v any) {
	decoder := json.NewDecoder(body)
	if err := decoder.Decode(v); err != nil {
		applog.Danger("Cannot parse form", err)
	}
}

func ReqTreePars(r *http.Request) TreeRequest {
	tp := TreeRequest{}
	doJSONFromBody(r.Body, &tp)

	return tp
}
