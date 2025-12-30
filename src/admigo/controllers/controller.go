// Package controllers
package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"admigo/config"
	"admigo/lang"
	"admigo/models/mcom"
)

const (
	BACKA string = "backa"
)

func setFrontContent(w http.ResponseWriter, r *http.Request,
	menuitem string, info any, funcs map[string]any,
	pages ...string,
) {
	qv := r.URL.Query()
	hsb := qv.Get("hsb")

	logged := LoggedUser(r)

	data := map[string]any{"logged": logged}

	if info != nil {
		data["info"] = info
	}

	sb, tagItem := CreateSidebarWeb(menuitem)

	auth := false
	if tagItem != nil {
		auth = tagItem.Auth
	}

	data["sb"] = sb

	data["lang"] = config.Env(false).Lang
	data["stat"] = config.Env(false).Static
	data["needtra"] = 0
	data["hidesb"] = 0

	if lang.NeedTra() {
		data["needtra"] = 1
	}

	if len(hsb) > 0 {
		data["hidesb"] = 1
	}

	if logged == nil && auth {
		cookie := http.Cookie{
			Name:     BACKA,
			Path:     "/",
			Value:    ro(menuitem),
			HttpOnly: true,
			Expires:  time.Now().Add(60 * time.Minute),
		}
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, ro("login_get"), http.StatusFound)

		return
	}

	listPages := []string{
		"lays/layout",
		"stru/sidebar",
		"stru/unava",
		"stru/uava",
	}

	listPages = append(listPages, pages...)

	GenerateHTML(w, r, data, funcs, listPages...)
}

func getFrontContentAjax(r *http.Request, info any, funcs map[string]any, pages ...string) (ret string, err error) {
	logged := LoggedUser(r)

	listPages := []string{
		"lays/layout_ajax",
	}

	data := map[string]any{
		"logged": logged,
		"info":   info,
	}

	listPages = append(listPages, pages...)

	return GetHTMLAjax(data, funcs, listPages...)
}

func APIError(w http.ResponseWriter, err error) {
	res := mcom.GetErrorResult(map[string]string{"api": err.Error()})
	w.WriteHeader(500)
	output, _ := json.MarshalIndent(res, "", "\t")
	w.Write(output)
}

func WebError(w http.ResponseWriter, r *http.Request, err error, menuitem string) {
	dber := map[string]string{
		"error": err.Error(),
	}

	setFrontContent(w, r, menuitem, dber, nil, "stru/web_error")
}
