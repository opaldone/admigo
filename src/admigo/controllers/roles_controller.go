package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"admigo/common"
	"admigo/models/mcom"
	"admigo/models/roles"

	"github.com/julienschmidt/httprouter"
)

func doRolesWebError(w http.ResponseWriter, r *http.Request, err error) {
	WebError(w, r, err, roles.MN_ROLES)
}

func getRolesSort(qv url.Values) *map[string]string {
	so := qv.Get("sort")

	if len(so) == 0 {
		return &map[string]string{
			"al": "",
		}
	}

	de := ""
	if strings.Contains(so, "-") {
		de = "desc"
		so = strings.ReplaceAll(so, "-", "")
	}

	return &map[string]string{
		so: de,
	}
}

func findRole(qv url.Values) (ret *roles.RoleModel, err error) {
	sid := qv.Get("sid")

	if len(sid) == 0 {
		return
	}

	return roles.RoleWithPaths(sid)
}

func getRolesFilter(qv url.Values) *map[string]string {
	fil := map[string]string{}

	pk := qv.Get("pk")

	if len(pk) > 0 {
		int_pk, err := strconv.Atoi(pk)
		if err != nil {
			int_pk = -1
		}
		fil["pk"] = strconv.Itoa(int_pk)
	}

	return &fil
}

func RolesIndex(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	qv := r.URL.Query()

	var paths *[]string
	fid := ""
	var finded_role *roles.RoleModel
	var err error

	fil := getRolesFilter(qv)
	pk := ""

	if len((*fil)["pk"]) > 0 {
		pk = (*fil)["pk"]
	}

	if len(pk) == 0 {
		finded_role, err = findRole(qv)
		if err != nil {
			doRolesWebError(w, r, err)
			return
		}

		if finded_role != nil {
			paths = finded_role.Paths
			fid = finded_role.Al
		}
	}

	info := &roles.RoleListRequest{
		Query: &mcom.DataQuery{
			Link:    r.URL,
			OrderBy: getRolesSort(qv),
			Filter:  fil,
		},
		Tree: map[string]interface{}{
			"el":    "roles-tr-root",
			"root":  ro("roles_root"),
			"node":  ro("roles_node"),
			"paths": paths,
			"fid":   fid,
			"pk":    pk,
		},
		Finded: finded_role,
	}

	setFrontContent(w, r, roles.MN_ROLES, info, map[string]any{"label": roles.Label},
		"roles/ix/index", "roles/ix/_filter", "stru/tree", "stru/dlg", "stru/hidden_sort",
	)
}

func roleForm(w http.ResponseWriter, r *http.Request, mo *roles.RoleModel) {
	for_fo := map[string]interface{}{}
	for_fo["mo"] = mo
	setFrontContent(w, r, roles.MN_ROLES, for_fo, map[string]any{"label": roles.Label},
		"roles/ed/role_form", "roles/ed/_sign",
	)
}

func getTreePars(r *http.Request, isroot bool) (ret map[string]interface{}, err error) {
	rt := common.ReqTreePars(r)

	qlev := rt.Lev
	qpar := rt.Par
	qprid := rt.Prid
	qpaths := rt.Paths
	qfid := rt.Fid
	qpk := rt.Pk
	qroot := rt.Root

	ex_step := 30
	em_step := 10

	lev := -1
	if len(qlev) > 0 {
		lev, err = strconv.Atoi(qlev)
		if err != nil {
			return
		}
	}

	tl := lev + 1
	pd := 0
	stp := ""

	if lev >= 0 {
		pd = tl * ex_step
	}

	stp = "style=\"padding-left:" + strconv.Itoa(pd)
	if pd > 0 {
		stp = stp + "px;\""
	} else {
		stp = stp + ";\""
	}

	show_line := "true"
	if tl == 0 {
		show_line = ""
	}

	dp := ""
	if len(qprid) > 0 {
		dp = fmt.Sprintf("%s data-prid=\"%s\"", dp, qprid)
	}

	ret = map[string]interface{}{
		"ex_step":   strconv.Itoa(ex_step),
		"em_step":   strconv.Itoa(em_step),
		"stp":       stp,
		"show_line": show_line,
		"tl":        strconv.Itoa(tl),
		"dp":        dp,
		"par":       qpar,
		"paths":     qpaths,
		"fid":       qfid,
		"pk":        qpk,
		"isroot":    isroot,
		"root":      qroot,
	}

	return
}

func RolesRoot(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	qv := r.URL.Query()

	tr, err := getTreePars(r, true)
	if err != nil {
		ApiError(w, err)
		return
	}

	list, err := roles.GetRolesList(&roles.RoleListRequest{
		Query: &mcom.DataQuery{
			Link:    r.URL,
			OrderBy: getRolesSort(qv),
		},
		Tree: tr,
	})
	if err != nil {
		ApiError(w, err)
		return
	}

	co, err := getFrontContentAjax(r, list, map[string]any{"label": roles.Label},
		"roles/ix/_root",
		"roles/ix/_head",
		"roles/ix/_data",
	)
	if err != nil {
		ApiError(w, err)
		return
	}

	ans := common.AjaxAnswer{
		Cont: co,
	}

	output, _ := json.Marshal(ans)

	w.Write(output)
}

func RolesNode(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	qv := r.URL.Query()

	tr, err := getTreePars(r, false)
	if err != nil {
		ApiError(w, err)
		return
	}

	list, err := roles.GetRolesList(&roles.RoleListRequest{
		Query: &mcom.DataQuery{
			Link:    r.URL,
			OrderBy: getRolesSort(qv),
		},
		Tree: tr,
	})
	if err != nil {
		ApiError(w, err)
		return
	}

	co, err := getFrontContentAjax(r, list, nil, "roles/ix/_node", "roles/ix/_data")
	if err != nil {
		ApiError(w, err)
		return
	}

	ans := common.AjaxAnswer{
		Cont: co,
	}

	output, _ := json.Marshal(ans)

	w.Write(output)
}

func RolesShow(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	qv := r.URL.Query()

	qfi := qv.Get("fi")

	list, err := roles.GetForShow(qfi)
	if err != nil {
		ApiError(w, err)
		return
	}

	co := ""

	if len(list.List) > 0 {
		co, err = getFrontContentAjax(r, list, nil, "roles/ix/_slist")
		if err != nil {
			ApiError(w, err)
			return
		}
	}

	ans := common.AjaxAnswer{
		Cont: co,
	}

	output, _ := json.Marshal(ans)

	w.Write(output)
}

func RoleDelete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(ps.ByName("id"))

	mo, err := roles.LinkByID(id)
	if err != nil {
		ApiError(w, err)
		return
	}

	err = mo.DeleteLinkRole()
	if err != nil {
		ApiError(w, err)
		return
	}

	ok := mcom.GetOk("Link was deleted")
	output, _ := json.MarshalIndent(ok, "", "\t")

	w.Write(output)
}

func RoleAddGet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	qv := r.URL.Query()

	par := qv.Get("par")

	mo, err := roles.GetNewRole(par)
	if err != nil {
		doRolesWebError(w, r, err)
		return
	}

	roleForm(w, r, mo)
}

func RoleAddPost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	mo, err := roles.GetNewRole("")
	if err != nil {
		doRolesWebError(w, r, err)
		return
	}

	roles.RoleEditOrAdd(r, mo)

	if mo.Errors != nil {
		roleForm(w, r, mo)
		return
	}

	http.Redirect(w, r, ro(roles.MN_ROLES), http.StatusFound)
}

func RoleEditGet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	lid, _ := strconv.Atoi(ps.ByName("lid"))

	mo, err := roles.GetEditRole(lid)
	if err != nil {
		doRolesWebError(w, r, err)
		return
	}

	roleForm(w, r, mo)
}

func RoleEditPost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	lid, _ := strconv.Atoi(ps.ByName("lid"))

	mo, err := roles.GetEditRole(lid)
	if err != nil {
		doRolesWebError(w, r, err)
		return
	}

	roles.RoleEditOrAdd(r, mo)

	if mo.Errors != nil {
		roleForm(w, r, mo)
		return
	}

	http.Redirect(w, r, ro(roles.MN_ROLES), http.StatusFound)
}
