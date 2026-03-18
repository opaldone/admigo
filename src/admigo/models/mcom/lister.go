// Package mcom
package mcom

import (
	"database/sql"
	"fmt"
	"math"
	"net/url"
	"strings"
)

type DataQuery struct {
	Page    int                `json:"page,omitempty"`
	PerPg   int                `json:"perpage,omitempty"`
	CntPage int                `json:"cntpage,omitempty"`
	Link    *url.URL           `json:"link,omitempty"`
	OrderBy *map[string]string `json:"order_by,omitempty"`
	Filter  *map[string]string `json:"filter,omitempty"`
}

func GetRows(query *DataQuery, quein string, alias string, needwhere bool) (*sql.Rows, error) {
	que := fmt.Sprintf("%s%s", quein, _where(query, alias, needwhere))
	pages := _page(query, que)
	que = fmt.Sprintf("%s%s%s", que, _order(query), pages)

	return Dbc.Query(que)
}

func _where(query *DataQuery, alias string, needwhere bool) (wh string) {
	if query == nil {
		return
	}

	if query.Filter == nil {
		return
	}

	whlist := []string{}

	for fld, val := range *query.Filter {
		if fld == "pk" {
			whlist = append(whlist, fmt.Sprintf("%s.id = %s", alias, val))
			continue
		}

		whlist = append(whlist, fmt.Sprintf("lower(%s.%s) like lower('%%%s%%')", alias, fld, val))
	}

	if len(whlist) == 0 {
		return
	}

	empstr := "and "
	if needwhere {
		empstr = "where "
	}

	wh = fmt.Sprintf(empstr+"%s\n", strings.Join(whlist, " and "))

	return
}

func _order(query *DataQuery) (by string) {
	if query == nil {
		return
	}

	ord := ""

	for fld, so := range *query.OrderBy {
		if len(ord) > 0 {
			ord = ord + ","
		}
		ord = ord + fld
		if len(so) > 0 {
			ord = ord + " " + so
		}
	}

	if len(ord) > 0 {
		by = fmt.Sprintf("order by %s\n", ord)

		return
	}

	return
}

func _page(query *DataQuery, sql string) string {
	if query == nil {
		return ""
	}

	if query.Page == 0 && query.PerPg == 0 {
		return ""
	}

	count := 0
	cntSQL := fmt.Sprintf("select count(1) from (%s) q_cnt", sql)
	Dbc.QueryRow(cntSQL).Scan(&count)

	query.CntPage = int(math.Ceil(float64(count) / float64(query.PerPg)))

	if query.Page > query.CntPage {
		query.Page = 1
	}

	return fmt.Sprintf("limit %d offset %d\n", query.PerPg, query.PerPg*(query.Page-1))
}
