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
	menuitem string, info interface{}, funcs map[string]any,
	pages ...string,
) {
	logged := LoggedUser(r)

	data := map[string]interface{}{"logged": logged}

	if info != nil {
		data["info"] = info
	}

	sb, tag_item := CreateSidebarWeb(menuitem)

	auth := false
	if tag_item != nil {
		auth = tag_item.Auth
	}

	data["sb"] = sb

	data["lang"] = config.Env(false).Lang
	data["stat"] = config.Env(false).Static
	data["needtra"] = 0
	if lang.NeedTra() {
		data["needtra"] = 1
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

	list_pages := []string{
		"lays/layout",
		"stru/sidebar",
		"stru/unava",
		"stru/uava",
	}

	list_pages = append(list_pages, pages...)

	GenerateHTML(w, r, data, funcs, list_pages...)
}

func getFrontContentAjax(r *http.Request, info interface{}, funcs map[string]any, pages ...string) (ret string, err error) {
	logged := LoggedUser(r)

	list_pages := []string{
		"lays/layout_ajax",
	}

	data := map[string]interface{}{
		"logged": logged,
		"info":   info,
	}

	list_pages = append(list_pages, pages...)

	return GetHTMLAjax(data, funcs, list_pages...)
}

func ApiError(w http.ResponseWriter, err error) {
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
