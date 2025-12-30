package users

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"admigo/applog"
	"admigo/common"
	"admigo/config"
	"admigo/lang"
	"admigo/models/attrs"
	"admigo/models/mcom"
	"admigo/models/roles"
)

// UserModel model for a user
type UserModel struct {
	ID        int                `json:"id,omitempty"`
	Ac        string             `json:"ac,omitempty"`
	Em        string             `json:"em,omitempty"`
	Pas       string             `json:"pas,omitempty"`
	CreatedAt *time.Time         `json:"created_at,omitempty"`
	Confirmed int                `json:"confirmed,omitempty"`
	Prot      uint               `json:"prot,omitempty"`
	Name      string             `json:"name,omitempty"`
	Thumb     string             `json:"thumb,omitempty"`
	Bim       string             `json:"bim,omitempty"`
	Errors    *map[string]string `json:"errors,omitempty"`
	Attrs     []*attrs.Attr      `json:"attrs,omitempty"`
	Roles     []*roles.RoleModel `json:"roles,omitempty"`
	Ous       *UserModel         `json:"ous,omitempty"`
	canRoles  map[string]bool
}

const (
	IMGPATH           = "/images/users/"
	ADMIN_ROLE string = "admu"
	MN_USERS   string = "users"
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
		tras = lang.LoadLabels("user")
	}

	labels = map[string]string{
		"ac":         lang.Re("Login"),
		"em":         lang.Re("Email"),
		"pas":        lang.Re("Password"),
		"th":         re("Thumb"),
		"bim":        re("Big image"),
		"name":       re("Name"),
		"sign":       re("Sign"),
		"attributes": re("Attributes"),
		"roles":      re("Roles"),
		"add_title":  re("Add new user"),
		"edit_title": re("Edit user"),
		"ch_role":    re("Choose role"),
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

func (user *UserModel) CheckFields() {
	if user.Errors != nil {
		user.Errors = nil
	}

	err := make(map[string]string)

	mcom.Required(&err, map[string][]string{
		"ac": {user.Ac, Label("ac")},
		"em": {user.Em, Label("em")},
	})

	if user.ID < 0 {
		mcom.Required(&err, map[string][]string{
			"pas": {user.Pas, Label("pas")},
		})
	}

	if len(err) > 0 {
		user.Errors = &err
	}
}

func (user *UserModel) GetError(fld string) string {
	if user.Errors == nil {
		return ""
	}

	return (*user.Errors)[fld]
}

func (user *UserModel) AddError(key string, err_msg string) {
	if user.Errors == nil {
		err := make(map[string]string)
		user.Errors = &err
	}

	(*user.Errors)[key] = err_msg
}

func (user *UserModel) imgPath() (ret string) {
	ret = fmt.Sprintf("/%s%s", config.Env(false).Static, IMGPATH)

	return
}

func (user *UserModel) GetThumb() (pth string) {
	if len(user.Thumb) == 0 {
		return ""
	}

	pth = user.imgPath() + user.Thumb

	return
}

func (user *UserModel) GetBim() (pth string) {
	if len(user.Bim) == 0 {
		return ""
	}

	pth = user.imgPath() + user.Bim

	return
}

func (user *UserModel) userSession() (session *SessionModel, err error) {
	que := `
		select uid, user_id, created_at
		from sessions
		where user_id = $1
	`
	se := SessionModel{}

	err = mcom.Dbc.QueryRow(que, user.ID).Scan(&se.UID, &se.UserID, &se.CreatedAt)
	session = &se
	return
}

func (user *UserModel) CreateSession() (session *SessionModel, err error) {
	if session, err = user.userSession(); err == nil {
		return
	}

	statement := `
		insert into sessions (uid, user_id)
		values ($1, $2)
		returning uid, user_id, created_at
	`

	stmt, err := mcom.Dbc.Prepare(statement)
	if err != nil {
		return
	}

	defer stmt.Close()

	se := SessionModel{}
	err = stmt.QueryRow(common.CreateUID(), user.ID).Scan(&se.UID, &se.UserID, &se.CreatedAt)

	session = &se
	return
}

func (user *UserModel) DeleteSessions() (err error) {
	que := "delete from sessions where user_id = $1"
	_, err = mcom.Dbc.Exec(que, user.ID)
	return
}

func (user *UserModel) SetConfirmed() (err error) {
	que := "update users set confirmed = 1 where id = $1"
	_, err = mcom.Dbc.Exec(que, user.ID)
	return
}

func (user *UserModel) removeImg(img string) (err error) {
	if len(img) == 0 {
		return
	}

	pt := fmt.Sprintf(".%s/%s", user.imgPath(), img)

	err = os.Remove(pt)
	if err != nil {
		return
	}

	return
}

func (user *UserModel) removeThumb() (err error) {
	return user.removeImg(user.Thumb)
}

func (user *UserModel) removeBim() (err error) {
	return user.removeImg(user.Bim)
}

func (user *UserModel) removeFiles() (err error) {
	err = user.removeThumb()
	if err != nil {
		return
	}

	err = user.removeBim()
	if err != nil {
		return
	}

	return
}

func (user *UserModel) setAttr(al_in string, val_in string) {
	if user.Attrs == nil {
		return
	}

	for _, attr_in := range user.Attrs {
		if attr_in.Al != al_in {
			continue
		}

		attr_in.Val = val_in
		return
	}
}

func (user *UserModel) SetThumb(val_in string) {
	user.setAttr("th", val_in)
	user.Thumb = val_in
}

func (user *UserModel) SetBim(val_in string) {
	user.setAttr("bim", val_in)
	user.Bim = val_in
}

func (user *UserModel) FillAttrs() (err error) {
	que := `
		select a.al, a.nm, a.tp, a.so, coalesce(ua.val, '') val
		from attrs a
		left join user_attrs ua on ua.attr = a.al and ua.user_id = $1
		where a.tp = 'users'
		order by so
	`

	rows, err := mcom.Dbc.Query(que, user.ID)
	if err != nil {
		return
	}

	defer rows.Close()

	list := []*attrs.Attr{}

	for rows.Next() {
		ar := attrs.Attr{}

		err = rows.Scan(&ar.Al, &ar.Nm, &ar.Tp, &ar.So, &ar.Val)
		if err != nil {
			return
		}

		list = append(list, &ar)
	}

	user.Attrs = list

	return
}

func (user *UserModel) getRolesByUser() (rows *sql.Rows, err error) {
	role_sel := roles.GetSelect()

	que := `
		select qu.al, qu.alnm
		from usro_links ur
		join (` + role_sel + `) qu on qu.al = ur.role_al
		where ur.user_id = $1
		order by al
	`

	rows, err = mcom.Dbc.Query(que, user.ID)
	if err != nil {
		return
	}

	return
}

func (user *UserModel) getRolesByAliases(aliases string) (rows *sql.Rows, err error) {
	role_sel := roles.GetSelect()
	que := `
		select qu.al, qu.alnm
		from (` + role_sel + `) qu
		where qu.al in (` + aliases + `)
		and length(qu.al) > 0
		order by al
	`

	rows, err = mcom.Dbc.Query(que)
	if err != nil {
		return
	}

	return
}

func (user *UserModel) FillRoles(aliases string) (err error) {
	var rows *sql.Rows

	if len(aliases) == 0 {
		rows, err = user.getRolesByUser()
	} else {
		rows, err = user.getRolesByAliases(aliases)
	}

	if err != nil {
		return
	}

	defer rows.Close()

	list := []*roles.RoleModel{}

	for rows.Next() {
		rl := roles.RoleModel{}

		err = rows.Scan(&rl.Al, &rl.AlNm)
		if err != nil {
			return
		}

		list = append(list, &rl)
	}

	user.Roles = list

	return
}

func (user *UserModel) fillCanRoles() (err error) {
	que := `
		with recursive tf as (
			select o.*, 0 lev
			from role_links o
			where o.child in (select ur.role_al from usro_links ur where ur.user_id = $1)
			union all
			select f.*, r.lev+1 lev
			from tf r
			join role_links f on f.parent = r.child
		)
		select distinct child al
		from tf q
		order by al
	`

	rows, err := mcom.Dbc.Query(que, user.ID)
	if err != nil {
		return
	}

	defer rows.Close()

	mapa := make(map[string]bool)
	var map_key string

	for rows.Next() {

		err = rows.Scan(&map_key)
		if err != nil {
			return
		}

		mapa[map_key] = true
	}

	user.canRoles = mapa

	return
}

func (user *UserModel) canRole(role_al string) bool {
	if user.canRoles == nil {
		err := user.fillCanRoles()
		if err != nil {
			applog.Danger("logged.FillCanRoles", err)
		}
	}

	return user.canRoles[role_al]
}

func (user *UserModel) IsUsersAdmin() bool {
	return user.canRole(ADMIN_ROLE)
}

func (user *UserModel) IsRolesAdmin() bool {
	return user.canRole(roles.ADMIN_ROLE)
}

func (user *UserModel) SetOus() {
	if user.Ous == nil {
		user.Ous = new(UserModel)
	}

	user.Ous.ID = user.ID
	user.Ous.Ac = user.Ac
	user.Ous.Em = user.Em
	user.Ous.Thumb = user.Thumb
	user.Ous.Bim = user.Bim
}

func (user *UserModel) dbAttrs() (err error) {
	if len(user.Attrs) == 0 {
		return
	}

	un := ""
	for _, attr_in := range user.Attrs {
		if len(attr_in.Val) == 0 {
			continue
		}

		if len(un) > 0 {
			un = un + " union\n"
		}

		un = fmt.Sprintf("%sselect %d user_id, '%s' attr, '%s' val", un,
			user.ID, attr_in.Al, attr_in.Val,
		)
	}

	if len(un) == 0 {
		return
	}

	scr := `
		insert into user_attrs(user_id, attr, val)
		select q.user_id, q.attr, q.val
		from (
			%[1]s
		) q
		where not exists(
			select 1
			from user_attrs ua
			where ua.user_id = q.user_id
			and ua.attr = q.attr
		);
		delete from user_attrs gg
		where gg.user_id = %d
		and gg.attr not in (
			select q.attr
			from (
				%[1]s
			) q
		);
		update user_attrs ua set
		val = q.val
		from (
			%[1]s
		) q
		where ua.user_id = q.user_id
		and ua.attr = q.attr
		and ua.val != q.val;
	`
	scr = fmt.Sprintf(scr, un, user.ID)

	_, err = mcom.Dbc.Exec(scr)

	return
}

func (user *UserModel) dbRoles() (err error) {
	if len(user.Roles) == 0 {
		return
	}

	un := ""
	for _, rl_in := range user.Roles {
		if len(rl_in.Al) == 0 {
			continue
		}

		if len(un) > 0 {
			un = un + " union\n"
		}

		un = fmt.Sprintf("%sselect %d user_id, '%s' role_al", un,
			user.ID, rl_in.Al,
		)
	}

	if len(un) == 0 {
		return
	}

	scr := `
		insert into usro_links(user_id, role_al)
		select q.user_id, q.role_al
		from (
			%[1]s
		) q
		where not exists(
			select 1
			from usro_links ua
			where ua.user_id = q.user_id
			and ua.role_al = q.role_al
		);
		delete from usro_links gg
		where gg.user_id = %d
		and gg.role_al not in (
			select q.role_al
			from (
				%[1]s
			) q
		);
	`

	scr = fmt.Sprintf(scr, un, user.ID)

	_, err = mcom.Dbc.Exec(scr)

	return
}

func (user *UserModel) insertUser() (err error) {
	ins := `
		insert into users (ac, em, pas)
		values ($1, $2, $3)
		returning id
	`

	new_pas := common.Hapas(user.Pas)

	row, err := mcom.Dbc.Query(ins, user.Ac, user.Em, new_pas)
	if err != nil {
		return
	}

	defer row.Close()

	if !row.Next() {
		return
	}

	err = row.Scan(&user.ID)
	if err != nil {
		return
	}

	return
}

func (user *UserModel) updatePas() (err error) {
	if len(user.Pas) == 0 {
		return
	}

	upd := `
		update users set
		pas = $1,
		created_at = now()
		where id = $2
	`

	new_pas := common.Hapas(user.Pas)

	_, err = mcom.Dbc.Exec(upd, new_pas, user.ID)

	return
}

func (user *UserModel) updateThumbFile() (err error) {
	if user.Thumb == user.Ous.Thumb {
		return
	}

	err = user.Ous.removeThumb()

	return
}

func (user *UserModel) updateBimFile() (err error) {
	if user.Bim == user.Ous.Bim {
		return
	}

	err = user.Ous.removeBim()

	return
}

func (user *UserModel) updateUser() (err error) {
	upd := `
		update users set
		ac = $1,
		em = $2,
		created_at = now()
		where id = $3
		and (
			ac != $1 or
			em != $2
		)
	`

	_, err = mcom.Dbc.Exec(upd, user.Ac, user.Em, user.ID)
	if err != nil {
		return
	}

	err = user.updatePas()
	if err != nil {
		return
	}

	err = user.updateThumbFile()
	if err != nil {
		return
	}

	err = user.updateBimFile()

	return
}

func (user *UserModel) DoUser() (err error) {
	if user.ID < 0 {
		err = user.insertUser()
		if err != nil {
			return
		}
	} else {
		err = user.updateUser()
		if err != nil {
			return
		}
	}

	err = user.dbAttrs()
	if err != nil {
		return
	}

	err = user.dbRoles()

	return
}

func (user *UserModel) DeleteUser() (err error) {
	if err = user.removeFiles(); err != nil {
		return
	}

	del := "delete from users where id = $1"
	_, err = mcom.Dbc.Exec(del, user.ID)

	return
}
