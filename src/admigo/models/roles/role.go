package roles

import (
	"strings"

	"admigo/lang"
	"admigo/models/mcom"
)

type RoleModel struct {
	Lid     int                `json:"lid"`
	Al      string             `json:"al"`
	Nm      string             `json:"nm"`
	AlNm    string             `json:"alnm"`
	Paral   string             `json:"paral"`
	IsPar   uint               `json:"ispar"`
	Prot    uint               `json:"prot"`
	Parents *[]*RoleModel      `json:"parents"`
	Paths   *[]string          `json:"paths"`
	Errors  *map[string]string `json:"errors"`
}

const (
	ADMIN_ROLE string = "admr"
	MN_ROLES   string = "roles"
)

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
		tras = lang.LoadLabels("role")
	}

	labels = map[string]string{
		"al":         re("Alias"),
		"nm":         re("Description"),
		"paral":      re("Parent"),
		"frole":      re("Find role"),
		"add_title":  re("Add new role"),
		"edit_title": re("Edit role"),
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

func (role *RoleModel) checkLoop() (res int, err error) {
	if len(role.Paral) == 0 {
		return
	}

	que := `
		with recursive tp as (
			select op.*, 0 lev
			from role_links op
			where op.child = $1
			union all
			select fp.*, rp.lev+1 lev
			from role_links fp
			join tp rp on rp.child = fp.parent
		),
		tc as (
			select oc.*, 0 lev
			from role_links oc
			where oc.child = $2
			union all
			select fc.*, rc.lev+1 lev
			from role_links fc
			join tc rc on rc.child = fc.parent
		)
		select case when exists(
			select 1
			from tp
			join tc on tc.parent = tp.parent
			where tp.parent != ''
			and tc.parent != ''
			and tp.child != $1 and tp.parent != $2
		) then 1 else 0 end exists_loop
	`

	err = mcom.Dbc.QueryRow(que, role.Al, role.Paral).Scan(&res)
	if err != nil {
		return
	}

	return
}

func (role *RoleModel) CheckFields() {
	if role.Errors != nil {
		role.Errors = nil
	}

	err := make(map[string]string)

	mcom.Required(&err, map[string][]string{
		"al": {role.Al, Label("al")},
		"nm": {role.Nm, Label("nm")},
	})

	if len(err) > 0 {
		role.Errors = &err
		return
	}

	loo, lo_err := role.checkLoop()
	if lo_err != nil {
		role.AddError("all", lo_err.Error())
		return
	}

	if loo == 1 {
		role.AddError("all", "Infinite recursive loop")
	}
}

func (role *RoleModel) GetError(fld string) string {
	if role.Errors == nil {
		return ""
	}

	return (*role.Errors)[fld]
}

func (role *RoleModel) AddError(key string, err_msg string) {
	if role.Errors == nil {
		err := make(map[string]string)
		role.Errors = &err
	}

	(*role.Errors)[key] = err_msg
}

func (role *RoleModel) FillPaths() (err error) {
	que := `
		with recursive tf as (
			select f.id, f.child, f.parent, f.id lid, 0 lev
			from role_links f
			where f.child = $1
			union all
			select f.id, f.child, f.parent, r.lid, r.lev+1 lev
			from role_links f
			join tf r on r.parent = f.child
		)
		select string_agg(cast(f.id as varchar), ',' order by f.lev desc) pt
		from tf f
		group by f.lid
	`

	rows, err := mcom.Dbc.Query(que, role.Al)
	if err != nil {
		return
	}

	defer rows.Close()

	list := []string{}

	for rows.Next() {
		var row string

		err = rows.Scan(&row)
		if err != nil {
			return
		}

		list = append(list, row)
	}

	if len(list) == 0 {
		return
	}

	role.Paths = &list

	return
}

func (mo *RoleModel) FillAvaiParents() (err error) {
	sel := GetSelect()
	que := sel + `where r.al != $1
		order by al
	`

	mo_al := "-"
	if len(mo.Al) > 0 {
		mo_al = mo.Al
	}

	rows, err := mcom.Dbc.Query(que, mo_al)
	if err != nil {
		return
	}

	defer rows.Close()

	list := []*RoleModel{}

	for rows.Next() {
		item := new(RoleModel)

		err = rows.Scan(&item.Al, &item.Nm, &item.AlNm)
		if err != nil {
			return
		}

		list = append(list, item)
	}

	if len(list) == 0 {
		return
	}

	mo.Parents = &list

	return
}

func (link *RoleModel) DeleteLinkRole() (err error) {
	up := `
		update role_links up set
		parent = default
		where up.id in (
			with recursive tf as (
				select o.id, o.child, o.parent, o.id lid, 0 lev
				from role_links o
				where o.id = $1
				union all
				select f.id, f.child, f.parent, r.lid, r.lev+1 lev
				from role_links f
				join tf r on r.child = f.parent
			)
			select q.id
			from tf q
			where not exists(
				select 1
				from role_links ch
				where ch.child = q.child
				and ch.parent = ''
			)
		)
	`

	_, err = mcom.Dbc.Exec(up, link.Lid)
	if err != nil {
		return
	}

	del_link := "delete from role_links where id = $1"
	_, err = mcom.Dbc.Exec(del_link, link.Lid)
	if err != nil {
		return
	}

	del_ro := `
		delete
		from roles r
		where r.al = $1
		and r.prot = 0
		and not exists(
			select 1
			from role_links l
			where r.al in (l.child, l.parent)
		)
	`
	_, err = mcom.Dbc.Exec(del_ro, link.Al)

	return
}

func (mo *RoleModel) insertLink() (err error) {
	ins := `
		insert into role_links(child, parent)
		select q.child, q.parent
		from (
			select $1 child, $2 parent
		) q
		where not exists(
			select 1
			from role_links e
			where e.child = q.child
			and e.parent = q.parent
		)
	`

	_, err = mcom.Dbc.Exec(ins, mo.Al, mo.Paral)

	return
}

func (mo *RoleModel) updateLink() (err error) {
	upd := `
		update role_links set
		parent = $2
		where id = $1
		and parent != $2
	`

	_, err = mcom.Dbc.Exec(upd, mo.Lid, mo.Paral)

	return
}

func (mo *RoleModel) insertRole() (err error) {
	ins := `
		insert into roles(al, nm)
		select q.al, q.nm
		from (
			select $1 al, $2 nm
		) q
		where not exists(
			select 1
			from roles e
			where e.al = q.al
		)
		returning al
	`

	mo.Al = strings.Trim(strings.ToLower(mo.Al), " ")
	mo.Nm = strings.Trim(mo.Nm, " ")

	_, err = mcom.Dbc.Exec(ins, mo.Al, mo.Nm)
	if err != nil {
		return
	}

	err = mo.insertLink()

	return
}

func (mo *RoleModel) updateRole() (err error) {
	upd := `
		update roles set
		nm = $1
		where al = $2
		and nm != $1
	`

	mo.Nm = strings.Trim(mo.Nm, " ")

	_, err = mcom.Dbc.Exec(upd, mo.Nm, mo.Al)
	if err != nil {
		return
	}

	err = mo.updateLink()

	return
}

func (mo *RoleModel) DoRole() error {
	if mo.Lid < 0 {
		return mo.insertRole()
	}

	return mo.updateRole()
}
