package controllers

import (
	"net/http"
	"strconv"

	"admigo/models/users"

	"github.com/julienschmidt/httprouter"
)

func getShowUser(par_user string) (show_user *users.UserModel) {
	userid, err := strconv.Atoi(par_user)
	if err != nil {
		return
	}

	show_user, err = users.UserByID(userid)
	if err != nil {
		return
	}

	return
}

func PersonIndex(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var show_user *users.UserModel

	qv := r.URL.Query()

	par_user := qv.Get("user_id")

	if len(par_user) > 0 {
		show_user = getShowUser(par_user)

		if show_user == nil {
			http.Redirect(w, r, ro("person"), http.StatusFound)
			return
		}
	}

	if show_user == nil {
		show_user = LoggedUser(r)
	}

	show_user.FillAttrs()

	info := map[string]interface{}{}
	info["show_user"] = show_user

	setFrontContent(w, r, users.MN_USERS, info, nil,
		"person/ix/index",
		"person/ix/_ri",
		"person/ix/_le",
	)
}
