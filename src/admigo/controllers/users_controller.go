package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"admigo/models/mcom"
	"admigo/models/users"

	"github.com/julienschmidt/httprouter"
)

func doUsersWebError(w http.ResponseWriter, r *http.Request, err error) {
	WebError(w, r, err, users.MN_USERS)
}

func getUsersSort(qv url.Values) *map[string]string {
	so := qv.Get("sort")

	if len(so) == 0 {
		return &map[string]string{
			"id": "",
		}
	}

	de := ""
	if strings.Contains(so, "-") {
		de = "desc"
		so = strings.ReplaceAll(so, "-", "")
	}

	return &map[string]string{
		so: de,
	}
}

func getPage(qv url.Values) int {
	pg := qv.Get("pg")

	ret, _ := strconv.Atoi(pg)

	if ret <= 0 {
		ret = 1
	}

	return ret
}

func getUsersFilter(qv url.Values) *map[string]string {
	fil := map[string]string{}

	nm := qv.Get("nm")
	if len(nm) > 0 {
		fil["nm"] = nm
	}

	em := qv.Get("em")
	if len(em) > 0 {
		fil["em"] = em
	}

	pk := qv.Get("pk")
	if len(pk) > 0 {
		intPk, err := strconv.Atoi(pk)
		if err != nil {
			intPk = -1
		}
		fil["pk"] = strconv.Itoa(intPk)
	}

	return &fil
}

func UsersIndex(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	qv := r.URL.Query()

	list, err := users.List(&users.UserListRequest{
		Query: &mcom.DataQuery{
			Page:    getPage(qv),
			PerPg:   7,
			Link:    r.URL,
			OrderBy: getUsersSort(qv),
			Filter:  getUsersFilter(qv),
		},
	})
	if err != nil {
		doUsersWebError(w, r, err)
		return
	}

	setFrontContent(w, r, users.MN_USERS, list, map[string]any{"label": users.Label},
		"users/ix/index", "users/ix/_filter",
		"users/ix/_table", "users/ix/_row", "users/ix/_row_empty",
		"stru/dlg", "stru/hidden_sort",
	)
}

func userForm(w http.ResponseWriter, r *http.Request, mo *users.UserModel) {
	forFo := map[string]any{}
	forFo["mo"] = mo
	setFrontContent(w, r, users.MN_USERS, forFo, map[string]any{"label": users.Label},
		"users/ed/user_form",
		"users/ed/_sign",
		"users/ed/_thumb",
		"users/ed/_bimg",
		"users/ed/_attrs",
		"users/ed/_user_roles",
	)
}

func UserAddGet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	mo, err := users.GetNewUser()
	if err != nil {
		doUsersWebError(w, r, err)
		return
	}

	userForm(w, r, mo)
}

func UserAddPost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	mo, err := users.GetNewUser()
	if err != nil {
		doUsersWebError(w, r, err)
		return
	}

	users.UserEditOrAdd(r, mo)

	if mo.Errors != nil {
		userForm(w, r, mo)
		return
	}

	redir := fmt.Sprintf("%s?sort=-id", ro(users.MN_USERS))
	http.Redirect(w, r, redir, http.StatusFound)
}

func UserEditGet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID, _ := strconv.Atoi(ps.ByName("user_id"))

	mo, err := users.GetEditUser(userID)
	if err != nil {
		doUsersWebError(w, r, err)
		return
	}

	userForm(w, r, mo)
}

func UserEditPost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID, _ := strconv.Atoi(ps.ByName("user_id"))

	mo, err := users.GetEditUser(userID)
	if err != nil {
		doUsersWebError(w, r, err)
		return
	}

	users.UserEditOrAdd(r, mo)

	if mo.Errors != nil {
		userForm(w, r, mo)
		return
	}

	redir := fmt.Sprintf("%s?pk=%d", ro(users.MN_USERS), mo.ID)
	http.Redirect(w, r, redir, http.StatusFound)
}

func UserDelete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(ps.ByName("id"))

	user, err := users.UserByID(id)
	if err != nil {
		APIError(w, err)
		return
	}

	if user.Prot == 1 {
		APIError(w, errors.New("user is protected"))
		return
	}

	err = user.DeleteUser()
	if err != nil {
		APIError(w, err)
		return
	}

	ok := mcom.GetOk("User was deleted")
	output, _ := json.MarshalIndent(ok, "", "\t")

	w.Write(output)
}
