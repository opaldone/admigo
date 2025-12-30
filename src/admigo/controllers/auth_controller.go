package controllers

import (
	"errors"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"admigo/applog"
	"admigo/common"
	"admigo/config"
	"admigo/lang"
	"admigo/models/auth"
	"admigo/models/users"

	"github.com/gorilla/csrf"
	"github.com/julienschmidt/httprouter"
)

// UUIDCOOKI constant for the name of cookie
const (
	UUIDCOOKI string = "uuidcookie"
)

// LoggedUser returns logged user
func LoggedUser(r *http.Request) (user *users.UserModel) {
	if cook, err := r.Cookie(UUIDCOOKI); err == nil {
		user = users.SessionUser(cook.Value)
	}

	return
}

func setUUIDCookie(w http.ResponseWriter, uuid string) {
	cookie := http.Cookie{
		Name:     UUIDCOOKI,
		Path:     "/",
		Value:    uuid,
		HttpOnly: true,
		Expires:  time.Now().Add(30 * 24 * time.Hour),
	}
	http.SetCookie(w, &cookie)
}

func removeUUIDCookie(w http.ResponseWriter) {
	rc := http.Cookie{
		Name:    UUIDCOOKI,
		Path:    "/",
		MaxAge:  -1,
		Expires: time.Unix(1, 0),
	}
	http.SetCookie(w, &rc)
}

func getLoginFm(funcs map[string]any) (fm template.FuncMap) {
	fm = template.FuncMap{
		"ro": ro,
		"re": lang.Re,
		"dd": func(v interface{}) template.HTML {
			return template.HTML(
				common.ShowJSON(v, false),
			)
		},
	}

	for ke, fu := range funcs {
		fm[ke] = fu
	}

	return
}

func simpleContent(w http.ResponseWriter, r *http.Request, au *auth.AuthUser, form_par map[string]string, page string) {
	new_data := map[string]interface{}{
		"au": au,
		"fp": form_par,
	}

	new_data["footer"] = map[string]string{
		"year":    strconv.Itoa(time.Now().Year()),
		"appname": config.Env(false).Appname,
	}

	new_data["cs"] = map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
	}

	filenames := []string{
		"lays/layout_login",
		page,
	}

	fm := getLoginFm(map[string]any{
		"label": auth.Label,
	})

	tmpl := getSiteTemplates(filenames, fm)

	tmpl.ExecuteTemplate(w, "layout_login", new_data)
}

func setLoginContent(w http.ResponseWriter, r *http.Request, au *auth.AuthUser) {
	form_par := map[string]string{
		"action": ro("login_post"),
		"button": lang.Re("Sign in"),
	}

	simpleContent(w, r, au, form_par, "auth/login")
}

func setSignupContent(w http.ResponseWriter, r *http.Request, au *auth.AuthUser) {
	form_par := map[string]string{
		"action": ro("signup_post"),
		"button": lang.Re("Signup"),
	}

	simpleContent(w, r, au, form_par, "auth/signup")
}

// LoginGet login to the site
func LoginGet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	aus := new(auth.AuthUser)
	setLoginContent(w, r, aus)
}

func LoginPost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	aus, uid := auth.Login(r)

	if aus.Errors != nil {
		setLoginContent(w, r, aus)
		return
	}

	if len(uid) > 0 {
		setUUIDCookie(w, uid)
	}

	cook_back := ""
	if cook, err := r.Cookie(BACKA); err == nil {
		cook_back = cook.Value
	}

	if len(cook_back) == 0 {
		cook_back = ro("home")
	}

	http.Redirect(w, r, cook_back, http.StatusFound)
}

func SignupGet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	aus := new(auth.AuthUser)
	setSignupContent(w, r, aus)
}

func SignupPost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	aus := auth.CreateSignup(r)

	if aus.Errors != nil {
		setSignupContent(w, r, aus)
		return
	}

	http.Redirect(w, r, ro("home"), http.StatusFound)
}

func ConfirmUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	vals := r.URL.Query()
	key := vals.Get("key")
	em := vals.Get("email")

	if !common.CheckPas(key, em) {
		err := errors.New(lang.Re("Unfortunately this link does not allow registration"))
		WebError(w, r, err, "")
		return
	}

	uid, err := auth.ForceLogin(vals.Get("email"))
	if err != nil {
		WebError(w, r, err, "")
		return
	}

	setUUIDCookie(w, uid)
	http.Redirect(w, r, ro("home"), http.StatusFound)
}

func LogoutGet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	user := LoggedUser(r)

	if user == nil {
		http.Redirect(w, r, ro("home"), http.StatusFound)
		return
	}

	err := user.DeleteSessions()
	if err != nil {
		applog.Danger("DeleteSessions", err)
		http.Redirect(w, r, ro("home"), http.StatusFound)
		return
	}

	removeUUIDCookie(w)
	http.Redirect(w, r, ro("home"), http.StatusFound)
}
