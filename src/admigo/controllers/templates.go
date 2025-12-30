package controllers

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"

	"admigo/common"
	"admigo/lang"

	"github.com/gorilla/csrf"
)

func getFm(funcs map[string]any) (fm template.FuncMap) {
	fm = template.FuncMap{
		"attr": func(s string) template.HTMLAttr {
			return template.HTMLAttr(s)
		},
		"safe": func(s string) template.HTML {
			return template.HTML(s)
		},
		"ro": ro,
		"re": lang.Re,
		"dd": func(v any) template.HTML {
			return template.HTML(
				common.ShowJSON(v, false),
			)
		},
		"rand":    common.RandUID,
		"pth":     PrintTH,
		"ppg":     PrintPages,
		"strjoin": StringsJoin,
		"dict":    TemplDict,
	}

	for ke, fu := range funcs {
		fm[ke] = fu
	}

	return
}

func getSiteTemplates(filenames []string, fm template.FuncMap) (tmpl *template.Template) {
	var files []string

	for _, file := range filenames {
		files = append(files, fmt.Sprintf("templates/%s.html", file))
	}

	tmpl = template.Must(template.New("").Funcs(fm).ParseFiles(files...))

	return
}

func GenerateHTML(writer http.ResponseWriter, r *http.Request, data any, funcs map[string]any, filenames ...string) {
	if data == nil {
		data = map[string]any{
			csrf.TemplateTag: csrf.TemplateField(r),
		}
	}

	if data != nil {
		_, ok := data.(map[string]any)[csrf.TemplateTag]

		if !ok {
			data.(map[string]any)[csrf.TemplateTag] = csrf.TemplateField(r)
		}
	}

	fm := getFm(funcs)

	getSiteTemplates(filenames, fm).ExecuteTemplate(writer, "layout", data)
}

func GetHTMLAjax(data any, funcs map[string]any, filenames ...string) (ret string, err error) {
	fm := getFm(funcs)

	t := getSiteTemplates(filenames, fm)

	var buf bytes.Buffer

	err = t.ExecuteTemplate(&buf, "layout_ajax", data)
	if err != nil {
		return
	}

	ret = buf.String()

	return
}
