package users

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"admigo/lang"
	"admigo/models/mcom"
)

type UserListRequest struct {
	Query *mcom.DataQuery `json:"query,omitempty"`
	Users []*UserModel    `json:"users,omitempty"`
}

func getSelect() (res string) {
	res = `
		select us.id, us.ac, us.em, us.pas, us.created_at, us.confirmed, us.prot,
		coalesce(nullif(concat_ws(' ',
			(select ua.val from user_attrs ua where ua.user_id = us.id and ua.attr = 'ln'),
			(select ua.val from user_attrs ua where ua.user_id = us.id and ua.attr = 'fn'),
			(select ua.val from user_attrs ua where ua.user_id = us.id and ua.attr = 'mn')
		), ''), us.ac) nm,
		coalesce((select ua.val from user_attrs ua where ua.user_id = us.id and ua.attr = 'th'), '') th,
		coalesce((select ua.val from user_attrs ua where ua.user_id = us.id and ua.attr = 'bim'), '') bim
		from users us
	`
	return
}

func getUser(where_in string, val interface{}, hidePass bool) (mo *UserModel, err error) {
	if len(where_in) == 0 {
		return
	}

	model := UserModel{}

	que := fmt.Sprintf("%s%s",
		getSelect(),
		where_in,
	)

	err = mcom.Dbc.QueryRow(que, val).Scan(
		&model.ID, &model.Ac, &model.Em,
		&model.Pas, &model.CreatedAt, &model.Confirmed, &model.Prot,
		&model.Name, &model.Thumb, &model.Bim,
	)
	if err != nil {
		return
	}

	if hidePass {
		model.Pas = ""
	}

	mo = &model
	return
}

func List(data *UserListRequest) (resp *UserListRequest, err error) {
	que := getSelect()

	que = fmt.Sprintf("select q.*\nfrom (%s) q where q.prot = 0\n", que)

	rows, err := mcom.GetRows(data.Query, que, "q", false)
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		us := UserModel{}

		err = rows.Scan(&us.ID, &us.Ac, &us.Em, &us.Pas,
			&us.CreatedAt, &us.Confirmed, &us.Prot, &us.Name, &us.Thumb, &us.Bim,
		)
		if err != nil {
			return
		}

		data.Users = append(data.Users, &us)
	}

	resp = data

	return
}

func UserByAcem(se string) (user *UserModel, err error) {
	user, err = getUser("where $1 in (us.em, us.ac)", se, false)
	return
}

func UserByID(id int) (u *UserModel, err error) {
	u, err = getUser("where us.id = $1", id, true)
	return
}

func GetNewUser() (new_user *UserModel, err error) {
	new_user = new(UserModel)
	new_user.ID = -1
	err = new_user.FillAttrs()

	return
}

func GetEditUser(user_id int) (ed_user *UserModel, err error) {
	ed_user, err = UserByID(user_id)

	if err == sql.ErrNoRows {
		err = errors.New(lang.Re("User not found"))
		return
	}

	if err != nil {
		return
	}

	if ed_user.Prot == 1 {
		err = errors.New(lang.Re("User is protected"))
		return
	}

	err = ed_user.FillAttrs()
	if err != nil {
		return
	}

	err = ed_user.FillRoles("")
	if err != nil {
		return
	}

	ed_user.SetOus()

	return
}

func postUserRoles(r *http.Request, mo *UserModel) (err error) {
	list := r.PostForm["list_roles[]"]

	if len(list) == 0 {
		return
	}

	als := "'" + strings.Join(list, "','") + "'"

	err = mo.FillRoles(als)

	return
}

func UserEditOrAdd(r *http.Request, mo *UserModel) *UserModel {
	r.ParseForm()

	mo.Ac = r.PostFormValue("ac")
	mo.Em = r.PostFormValue("em")
	mo.Pas = r.PostFormValue("pas")

	rer := postUserRoles(r, mo)
	if rer != nil {
		mo.AddError("all", rer.Error())
		return mo
	}

	for _, attr_in := range mo.Attrs {
		attr_in.Val = r.PostFormValue(fmt.Sprintf("attr_%s", attr_in.Al))
	}

	mo.CheckFields()

	if mo.Errors != nil {
		return mo
	}

	thumb, err := mcom.SaveSomeFile(r, "user-img", "users")
	if err != nil {
		mo.AddError("all", err.Error())
		return mo
	}

	if len(thumb) > 0 {
		mo.SetThumb(thumb)
	}

	bim, err := mcom.SaveSomeFile(r, "user-bim", "users")
	if err != nil {
		mo.AddError("all", err.Error())
		return mo
	}

	if len(bim) > 0 {
		mo.SetBim(bim)
	}

	ner := mo.DoUser()
	if ner != nil {
		mo.AddError("all", ner.Error())
	}

	return mo
}
