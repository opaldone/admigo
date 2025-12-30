package roles

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"admigo/common"
	"admigo/lang"
	"admigo/models/mcom"
)

type RoleListRequest struct {
	Query  *mcom.DataQuery `json:"query,omitempty"`
	Tree   map[string]any  `json:"tree,omitempty"`
	Finded *RoleModel      `json:"finded,omitempty"`
	List   []*RoleModel    `json:"list,omitempty"`
}

func GetSelect() (que string) {
	que = `
		select r.al, r.nm, concat(r.nm, ' ', U&'\2015', ' ', r.al) alnm
		from roles r
	`

	return
}

func getSelectLinks(fids string) (res string) {
	res = `
		select l.id lid, rc.al, rc.nm, l.parent paral, case when exists(
			select 1
			from role_links lp
			where lp.parent = rc.al`

	if len(fids) > 0 {
		res = fmt.Sprintf("%s\nand lp.id in (%s)", res, fids)
	}

	res = res + `
		) then 1 else 0 end ispar,
		concat(rc.nm, ' ', U&'\2015', ' ', rc.al) alnm,
		rc.prot
		from role_links l
		join roles rc on rc.al = l.child
		join roles rp on rp.al = l.child
	`

	return
}

func getRoleOrLink(sel string, where_in string, val interface{}, dest ...any) (err error) {
	if len(where_in) == 0 {
		return
	}

	que := fmt.Sprintf("%s%s",
		sel,
		where_in,
	)

	err = mcom.Dbc.QueryRow(que, val).Scan(
		dest...,
	)
	if err != nil {
		return
	}

	return
}

func getRolesData(data *RoleListRequest, que string) (resp *RoleListRequest, err error) {
	rows, err := mcom.GetRows(data.Query, que, "l", false)
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		role := RoleModel{}

		err = rows.Scan(&role.Lid, &role.Al, &role.Nm, &role.Paral, &role.IsPar, &role.AlNm, &role.Prot)
		if err != nil {
			return
		}

		data.List = append(data.List, &role)
	}

	resp = data

	return
}

func GetRolesList(data *RoleListRequest) (resp *RoleListRequest, err error) {
	par := data.Tree["par"]
	paths := data.Tree["paths"]
	fid := data.Tree["fid"]
	pk := data.Tree["pk"]

	wh := ""
	wh_ids := ""
	fids := ""

	if len(pk.(string)) > 0 {
		fids = pk.(string)
		wh_ids = fmt.Sprintf("l.id in (%s)", fids)
	}

	if len(fids) == 0 && paths.(*[]string) != nil && par.(string) != fid.(string) {
		fids = strings.Join(*paths.(*[]string), ",")
		wh_ids = fmt.Sprintf(" and l.id in (%s)", fids)
	}

	que := getSelectLinks(fids)

	if len(pk.(string)) == 0 {
		wh = "l.parent = ''"

		if len(par.(string)) > 0 {
			wh = "l.parent = '" + par.(string) + "'"
		}
	}

	wh = fmt.Sprintf("%s%s", wh, wh_ids)

	que = fmt.Sprintf("%s where %s\n", que, wh)

	return getRolesData(data, que)
}

func GetForShow(fi string) (resp *RoleListRequest, err error) {
	que := GetSelect() + `
		where (
			lower(r.al) like concat('%', lower($1), '%') or
			lower(r.nm) like concat('%', lower($1), '%')
		)
		and length(r.al) > 0
		order by al
	`

	rows, err := mcom.Dbc.Query(que, fi)
	if err != nil {
		return
	}

	defer rows.Close()

	list := []*RoleModel{}

	for rows.Next() {
		model := RoleModel{}

		err = rows.Scan(&model.Al, &model.Nm, &model.AlNm)
		if err != nil {
			return
		}

		list = append(list, &model)
	}

	ret := RoleListRequest{
		List: list,
	}

	resp = &ret

	return
}

func RoleByAl(al string) (mo *RoleModel, err error) {
	mo = &RoleModel{}

	err = getRoleOrLink(GetSelect(),
		"where r.al = $1", al, &mo.Al, &mo.Nm, &mo.AlNm,
	)

	return
}

func LinkByID(id int) (mo *RoleModel, err error) {
	mo = &RoleModel{}

	err = getRoleOrLink(getSelectLinks(""),
		"where l.id = $1", id, &mo.Lid, &mo.Al, &mo.Nm, &mo.Paral, &mo.IsPar, &mo.AlNm, &mo.Prot,
	)

	return
}

func RoleWithPaths(al string) (role *RoleModel, err error) {
	role, err = RoleByAl(al)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return
	}

	role.FillPaths()

	return
}

func GetNewRole(par_in string) (new_role *RoleModel, err error) {
	new_role = new(RoleModel)
	new_role.Lid = -1
	new_role.Al = fmt.Sprintf("ral_%s", common.RandUID())

	if len(par_in) > 0 {
		new_role.Paral = par_in
	}

	err = new_role.FillAvaiParents()

	return
}

func GetEditRole(lid int) (ed_role *RoleModel, err error) {
	ed_role, err = LinkByID(lid)

	if err == sql.ErrNoRows {
		err = errors.New(lang.Re("Role not found"))
		return
	}

	if err != nil {
		return
	}

	if ed_role.Prot == 1 {
		err = errors.New(lang.Re("Role is protected"))
		return
	}

	err = ed_role.FillAvaiParents()

	return
}

func RoleEditOrAdd(r *http.Request, mo *RoleModel) *RoleModel {
	r.ParseForm()

	mo.Al = r.PostFormValue("al")
	mo.Nm = r.PostFormValue("nm")
	mo.Paral = r.PostFormValue("paral")

	mo.CheckFields()

	if mo.Errors != nil {
		return mo
	}

	ner := mo.DoRole()

	if ner != nil {
		mo.AddError("all", ner.Error())
	}

	return mo
}
