package controllers

import (
	"errors"
	"net/http"

	"admigo/lang"
	"admigo/models/roles"
	"admigo/models/users"

	"github.com/julienschmidt/httprouter"
)

func retError(w http.ResponseWriter, r *http.Request, from_api bool, mn string) {
	if from_api {
		APIError(w, errors.New(lang.Re("Insufficient rights")))
		return
	}

	WebError(w, r, errors.New(lang.Re("Insufficient rights")), mn)
}

func LoggedUsersAdmin(h httprouter.Handle, from_api bool) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		logged := LoggedUser(r)

		if logged != nil && logged.IsUsersAdmin() {
			h(w, r, ps)
			return
		}

		retError(w, r, from_api, users.MN_USERS)
	}
}

func LoggedRolesAdmin(h httprouter.Handle, from_api bool) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		logged := LoggedUser(r)

		if logged != nil && logged.IsRolesAdmin() {
			h(w, r, ps)
			return
		}

		retError(w, r, from_api, roles.MN_ROLES)
	}
}
