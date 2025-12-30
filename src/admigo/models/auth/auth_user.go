package auth

import (
	"admigo/lang"
	"admigo/models/mcom"
)

type AuthUser struct {
	Acem   string             `json:"acem"`
	Ac     string             `json:"ac"`
	Em     string             `json:"em"`
	Pas    string             `json:"pas"`
	Cpas   string             `json:"cpas"`
	Errors *map[string]string `json:"errors"`
}

var (
	labels map[string]string
	tras   *map[string]string
)

func re(key string) string {
	if tras == nil {
		return key
	}

	ret, ok := (*tras)[key]
	if !ok {
		return key
	}

	return ret
}

func setLabels() {
	if len(labels) > 0 {
		return
	}

	if tras == nil {
		tras = lang.LoadLabels("auth_user")
	}

	labels = map[string]string{
		"acem":           re("Email or Login"),
		"ac":             re("Login"),
		"em":             re("Email"),
		"pas":            re("Password"),
		"cpas":           re("Confirm Password"),
		"login_title":    re("Login to your account"),
		"register_title": re("Registering a new account"),
	}
}

func Label(fld string) string {
	setLabels()

	ret := labels[fld]

	if len(ret) == 0 {
		return fld
	}

	return ret
}

func (au *AuthUser) CheckLoginFields() {
	if au.Errors != nil {
		au.Errors = nil
	}

	err := make(map[string]string)

	mcom.Required(&err, map[string][]string{
		"acem": {au.Acem, Label("acem")},
		"pas":  {au.Pas, Label("pas")},
	})

	if len(err) > 0 {
		au.Errors = &err
	}
}

func (au *AuthUser) CheckSignupFields() {
	if au.Errors != nil {
		au.Errors = nil
	}

	err := make(map[string]string)

	mcom.Required(&err, map[string][]string{
		"ac":   {au.Ac, Label("ac")},
		"em":   {au.Em, Label("em")},
		"pas":  {au.Pas, Label("pas")},
		"cpas": {au.Cpas, Label("cpas")},
	})

	mcom.Confirmed(&err, map[string][]string{
		"cpas": {au.Pas, au.Cpas, Label("cpas")},
	})

	if len(err) > 0 {
		au.Errors = &err
	}
}

func (au *AuthUser) GetError(fld string) string {
	if au.Errors == nil {
		return ""
	}

	return (*au.Errors)[fld]
}

func (au *AuthUser) AddError(key string, err_msg string) {
	if au.Errors == nil {
		err := make(map[string]string)
		au.Errors = &err
	}

	(*au.Errors)[key] = err_msg
}
