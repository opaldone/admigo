package controllers

import (
	"errors"
	"fmt"
	"html/template"
	"net/url"
	"strconv"
	"strings"

	"admigo/models/mcom"
)

func getNewLink(link *url.URL, par, val string, lastpath bool) string {
	rp := link.RequestURI()

	if !lastpath {
		spl := strings.Split(rp, "/")
		rp = strings.Join(spl[:len(spl)-1], "/")
	}

	new_link, _ := url.ParseRequestURI(rp)

	pars := new_link.Query()
	pars.Set(par, val)
	new_link.RawQuery = pars.Encode()

	return new_link.RequestURI()
}

func PrintTH(field string, label string, width string, qu *mcom.DataQuery, lastpath bool) template.HTML {
	kind, so := (*qu.OrderBy)[field]

	cls := ""
	sig := ""

	if so {
		cls = "sort"
		sig = "-"
	}

	if len(kind) > 0 {
		cls = cls + " " + kind
		sig = ""
	}

	if len(cls) > 0 {
		cls = " class=\"" + cls + "\""
	}

	head_href := getNewLink(qu.Link, "sort", sig+field, lastpath)

	str := fmt.Sprintf(
		`<th%s style="width:%s%%">
		<a href="%s"><span>%s</span></a>
		</th>`,
		cls, width, head_href, label,
	)

	return template.HTML(str)
}

func PrintPages(colspan int, qu *mcom.DataQuery, lastpath bool) template.HTML {
	imax := qu.CntPage

	if imax <= 1 {
		return ""
	}

	var pg int
	var ac string
	var pg_href string

	ret := `
      <tfoot>
        <tr>
          <td colspan="%d">
            <ul class="page-list clearfix">
			%s
            </ul>
          </td>
        </tr>
      </tfoot>
	`

	li_str := ""
	for i := 0; i < imax; i++ {
		pg = i + 1
		ac = fmt.Sprintf("<span>%d</span>", pg)

		if pg != qu.Page {
			pg_href = getNewLink(qu.Link, "pg", strconv.Itoa(pg), lastpath)
			ac = fmt.Sprintf("<a href=\"%s\">%d</a>", pg_href, pg)
		}

		li_str = fmt.Sprintf("%s<li>%s</li>", li_str, ac)
	}

	ret = fmt.Sprintf(ret, colspan, li_str)

	return template.HTML(ret)
}

func StringsJoin(dlm string, str_in ...string) string {
	ret := strings.Join(str_in, dlm)

	return strings.Trim(ret, dlm)
}

func TemplDict(values ...any) (map[string]any, error) {
	if len(values)%2 != 0 {
		return nil, errors.New("invalid dict call")
	}

	dict := make(map[string]any, len(values)/2)

	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)

		if !ok {
			return nil, errors.New("dict keys must be strings")
		}

		dict[key] = values[i+1]
	}

	return dict, nil
}
