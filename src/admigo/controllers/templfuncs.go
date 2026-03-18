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

	newLink, _ := url.ParseRequestURI(rp)

	pars := newLink.Query()
	pars.Set(par, val)
	newLink.RawQuery = pars.Encode()

	return newLink.RequestURI()
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

	headHref := getNewLink(qu.Link, "sort", sig+field, lastpath)

	str := fmt.Sprintf(
		`<th%s style="width:%s%%">
		<a href="%s"><span>%s</span></a>
		</th>`,
		cls, width, headHref, label,
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
	var pgHref string

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

	listr := ""
	for i := range imax {
		pg = i + 1
		ac = fmt.Sprintf("<span>%d</span>", pg)

		if pg != qu.Page {
			pgHref = getNewLink(qu.Link, "pg", strconv.Itoa(pg), lastpath)
			ac = fmt.Sprintf("<a href=\"%s\">%d</a>", pgHref, pg)
		}

		listr = fmt.Sprintf("%s<li>%s</li>", listr, ac)
	}

	ret = fmt.Sprintf(ret, colspan, listr)

	return template.HTML(ret)
}

func StringsJoin(dlm string, strIn ...string) string {
	ret := strings.Join(strIn, dlm)

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
