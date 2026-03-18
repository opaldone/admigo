package controllers

import (
	"net/http"
	"strconv"

	"admigo/models/users"

	"github.com/julienschmidt/httprouter"
)

func getShowUser(parUser string) (showUser *users.UserModel) {
	userid, err := strconv.Atoi(parUser)
	if err != nil {
		return
	}

	showUser, err = users.UserByID(userid)
	if err != nil {
		return
	}

	return
}

func PersonIndex(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var showUser *users.UserModel

	qv := r.URL.Query()

	parUser := qv.Get("user_id")

	if len(parUser) > 0 {
		showUser = getShowUser(parUser)

		if showUser == nil {
			http.Redirect(w, r, ro("person"), http.StatusFound)
			return
		}
	}

	if showUser == nil {
		showUser = LoggedUser(r)
	}

	showUser.FillAttrs()

	info := map[string]any{}
	info["showUser"] = showUser

	setFrontContent(w, r, users.MnUsers, info, nil,
		"person/ix/index",
		"person/ix/_ri",
		"person/ix/_le",
	)
}
