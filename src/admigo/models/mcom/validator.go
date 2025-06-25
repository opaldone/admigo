package mcom

import (
	"fmt"

	"admigo/lang"
)

const (
	F_ERROR = "error"
)

type Result struct {
	Message string            `json:"message"`
	Name    string            `json:"name,omitempty"`
	Errors  map[string]string `json:"errors,omitempty"`
}

func GetErrorResult(errors map[string]string) (res *Result) {
	res = &Result{Message: F_ERROR, Errors: errors}
	return
}

func GetOk(messageIn string) (res *Result) {
	res = &Result{Message: messageIn}
	return
}

func Required(errs *map[string]string, fields map[string][]string) {
	for field, val := range fields {
		if len(val[0]) == 0 {
			(*errs)[field] = fmt.Sprintf(lang.Re("The %s field is required"), val[1])
		}
	}
}

func Confirmed(errs *map[string]string, values map[string][]string) {
	for field, val := range values {
		if val[0] != val[1] {
			(*errs)[field] = fmt.Sprintf(lang.Re("The %s confirmation does not match"), val[2])
		}
	}
}
