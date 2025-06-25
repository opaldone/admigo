package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"admigo/config"

	"github.com/julienschmidt/httprouter"
)

type route struct {
	method  string
	pattern string
	handle  httprouter.Handle
}

type routes = map[string]route

var list routes

func init() {
	list = routes{
		"home":   route{"GET", "/", Index},
		"person": route{"GET", "/person", PersonIndex},

		"login_get":   route{"GET", "/login", LoginGet},
		"login_post":  route{"POST", "/login", LoginPost},
		"logout_get":  route{"GET", "/logout", LogoutGet},
		"signup_get":  route{"GET", "/signup", SignupGet},
		"signup_post": route{"POST", "/signup", SignupPost},
		"confirm":     route{"GET", "/confirm/*filepath", ConfirmUser},

		"users":          route{"GET", "/users", UsersIndex},
		"user_delete":    route{"DELETE", "/users/:id", LoggedUsersAdmin(UserDelete, true)},
		"user_add_get":   route{"GET", "/users/add", LoggedUsersAdmin(UserAddGet, false)},
		"user_add_post":  route{"POST", "/users/add", LoggedUsersAdmin(UserAddPost, false)},
		"user_edit_get":  route{"GET", "/users/edit/:user_id", LoggedUsersAdmin(UserEditGet, false)},
		"user_edit_post": route{"POST", "/users/edit/:user_id", LoggedUsersAdmin(UserEditPost, false)},

		"roles":          route{"GET", "/roles", RolesIndex},
		"role_show":      route{"GET", "/roles/show", RolesShow},
		"roles_root":     route{"POST", "/roles/root", RolesRoot},
		"roles_node":     route{"POST", "/roles/node", RolesNode},
		"role_delete":    route{"DELETE", "/roles/:id", LoggedRolesAdmin(RoleDelete, true)},
		"role_add_get":   route{"GET", "/roles/add", LoggedRolesAdmin(RoleAddGet, false)},
		"role_add_post":  route{"POST", "/roles/add", LoggedRolesAdmin(RoleAddPost, false)},
		"role_edit_get":  route{"GET", "/roles/edit/:lid", LoggedRolesAdmin(RoleEditGet, false)},
		"role_edit_post": route{"POST", "/roles/edit/:lid", LoggedRolesAdmin(RoleEditPost, false)},
	}
}

// GetRouters returns routers
func GetRouters() (router *httprouter.Router) {
	router = httprouter.New()

	st := config.Env(false).Static
	fpt := fmt.Sprintf("/%s/*filepath", st)
	router.ServeFiles(fpt, http.Dir(st))

	for _, r := range list {
		router.Handle(r.method, r.pattern, r.handle)
	}

	return
}

func ro(alias string, pars ...string) string {
	pat := list[alias].pattern
	pata := strings.Split(pat, ":")

	if len(pata) == 1 {
		return pat
	}

	if len(pars) == 0 {
		return pat
	}

	purl := ""
	for _, par := range pars {
		if len(purl) > 0 {
			purl += "/"
		}
		purl += par
	}

	ret := pata[0] + purl

	return ret
}
