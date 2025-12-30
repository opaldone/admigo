package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"admigo/config"

	"github.com/julienschmidt/httprouter"
)

type Route struct {
	method  string
	pattern string
	handle  httprouter.Handle
}

type Routes = map[string]Route

var routsList Routes

func init() {
	routsList = Routes{
		"home":   Route{"GET", "/", Index},
		"person": Route{"GET", "/person", PersonIndex},

		"login_get":   Route{"GET", "/login", LoginGet},
		"login_post":  Route{"POST", "/login", LoginPost},
		"logout_get":  Route{"GET", "/logout", LogoutGet},
		"signup_get":  Route{"GET", "/signup", SignupGet},
		"signup_post": Route{"POST", "/signup", SignupPost},
		"confirm":     Route{"GET", "/confirm/*filepath", ConfirmUser},

		"users":          Route{"GET", "/users", UsersIndex},
		"user_delete":    Route{"DELETE", "/users/:id", LoggedUsersAdmin(UserDelete, true)},
		"user_add_get":   Route{"GET", "/users/add", LoggedUsersAdmin(UserAddGet, false)},
		"user_add_post":  Route{"POST", "/users/add", LoggedUsersAdmin(UserAddPost, false)},
		"user_edit_get":  Route{"GET", "/users/edit/:user_id", LoggedUsersAdmin(UserEditGet, false)},
		"user_edit_post": Route{"POST", "/users/edit/:user_id", LoggedUsersAdmin(UserEditPost, false)},

		"roles":          Route{"GET", "/roles", RolesIndex},
		"role_show":      Route{"GET", "/roles/show", RolesShow},
		"roles_root":     Route{"POST", "/roles/root", RolesRoot},
		"roles_node":     Route{"POST", "/roles/node", RolesNode},
		"role_delete":    Route{"DELETE", "/roles/:id", LoggedRolesAdmin(RoleDelete, true)},
		"role_add_get":   Route{"GET", "/roles/add", LoggedRolesAdmin(RoleAddGet, false)},
		"role_add_post":  Route{"POST", "/roles/add", LoggedRolesAdmin(RoleAddPost, false)},
		"role_edit_get":  Route{"GET", "/roles/edit/:lid", LoggedRolesAdmin(RoleEditGet, false)},
		"role_edit_post": Route{"POST", "/roles/edit/:lid", LoggedRolesAdmin(RoleEditPost, false)},

		"loca":      Route{"GET", "/map/loca", MapLoca},
		"loca_ws":   Route{"POST", "/map/locaws", MapWs},
		"loca_show": Route{"GET", "/map/locash/:ci", MapShow},

		"apki": Route{"GET", "/apps", ApkiIndex},
	}
}

// GetRouters returns routers
func GetRouters() (router *httprouter.Router) {
	router = httprouter.New()

	st := config.Env(false).Static
	fpt := fmt.Sprintf("/%s/*filepath", st)
	router.ServeFiles(fpt, http.Dir(st))

	for _, r := range routsList {
		router.Handle(r.method, r.pattern, r.handle)
	}

	return
}

func ro(alias string, pars ...string) string {
	pat := routsList[alias].pattern
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
