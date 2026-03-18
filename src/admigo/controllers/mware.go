package controllers

import (
	"errors"
	"net/http"

	"admigo/lang"
	"admigo/models/roles"
	"admigo/models/users"

	"github.com/julienschmidt/httprouter"
)

func retError(w http.ResponseWriter, r *http.Request, fromAPI bool, mn string) {
	if fromAPI {
		APIError(w, errors.New(lang.Re("Insufficient rights")))
		return
	}

	WebError(w, r, errors.New(lang.Re("Insufficient rights")), mn)
}

func LoggedUsersAdmin(h httprouter.Handle, fromAPI bool) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		logged := LoggedUser(r)

		if logged != nil && logged.IsUsersAdmin() {
			h(w, r, ps)
			return
		}

		retError(w, r, fromAPI, users.MnUsers)
	}
}

func LoggedRolesAdmin(h httprouter.Handle, fromAPI bool) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		logged := LoggedUser(r)

		if logged != nil && logged.IsRolesAdmin() {
			h(w, r, ps)
			return
		}

		retError(w, r, fromAPI, roles.MnRoles)
	}
}
