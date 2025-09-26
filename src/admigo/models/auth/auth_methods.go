package auth

import (
	"errors"
	"net/http"

	"admigo/common"
	"admigo/lang"
	"admigo/mail"
	"admigo/models/users"
)

func Login(r *http.Request) (aus *AuthUser, uid string) {
	r.ParseForm()

	aus = new(AuthUser)
	aus.Acem = r.PostFormValue("acem")
	aus.Pas = r.PostFormValue("pas")

	aus.CheckLoginFields()

	if aus.Errors != nil {
		return
	}

	user, err := users.UserByAcem(aus.Acem)
	if err != nil {
		aus.AddError("acem", lang.Re("The Email or Login not found"))
		return
	}

	if user.Confirmed == 0 {
		aus.AddError("acem", lang.Re("User was not confirmed"))
		return
	}

	if !common.CheckPas(user.Pas, aus.Pas) {
		aus.AddError("pas", lang.Re("Incorrect Password"))
		return
	}

	session, err := user.CreateSession()
	if err != nil {
		aus.AddError("all", err.Error())
		return
	}

	uid = session.UID
	return
}

func CreateSignup(r *http.Request) (aus *AuthUser) {
	r.ParseForm()

	aus = new(AuthUser)
	aus.Ac = r.PostFormValue("ac")
	aus.Em = r.PostFormValue("em")
	aus.Pas = r.PostFormValue("pas")
	aus.Cpas = r.PostFormValue("cpas")

	aus.CheckSignupFields()

	if aus.Errors != nil {
		return
	}

	newUser := users.UserModel{
		ID:  -1,
		Ac:  aus.Ac,
		Em:  aus.Em,
		Pas: aus.Pas,
	}

	err := newUser.DoUser()
	if err != nil {
		aus.AddError("all", err.Error())
		return
	}

	mail.SendEmailUserRegistration(aus.Em)

	return
}

func ForceLogin(email string) (uid string, err error) {
	user, err := users.UserByAcem(email)
	if err != nil {
		err = errors.New(lang.Re("The user with this Email not found, please Signin again"))
		return
	}

	if err = user.SetConfirmed(); err != nil {
		return
	}

	sess, err := user.CreateSession()
	if err != nil {
		return
	}

	uid = sess.UID

	return
}
