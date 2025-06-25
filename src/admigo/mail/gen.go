package mail

import (
	"bytes"
	"fmt"
	"text/template"
)

func getMailTemplates(filenames []string) (tmpl *template.Template) {
	var files []string

	for _, file := range filenames {
		files = append(files, fmt.Sprintf("templates/%s.html", file))
	}

	tmpl = template.Must(template.ParseFiles(files...))

	return
}

// GenerateMail genereates email
func GenerateMail(bufer *bytes.Buffer, data interface{}, filenames ...string) {
	getMailTemplates(filenames).ExecuteTemplate(bufer, "layout_mail", data)
}
