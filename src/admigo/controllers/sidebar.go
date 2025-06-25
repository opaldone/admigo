package controllers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"

	"admigo/applog"
)

type sidebarMenu struct {
	Sidebar *[]sidebarTitle `json:"sidebar,omitempty"`
}

type sidebarTitle struct {
	Title string         `json:"title,omitempty"`
	Items *[]sidebarItem `json:"items,omitempty"`
}

type sidebarItem struct {
	Label string         `json:"label,omitempty"`
	Alias string         `json:"alias,omitempty"`
	Icon  string         `json:"icon,omitempty"`
	Auth  bool           `json:"auth,omitempty"`
	Items *[]sidebarItem `json:"items,omitempty"`
}

var (
	sidebar  *sidebarMenu
	tag_item *sidebarItem
)

func loadMenu() {
	file, err := os.Open("menu.json")
	if err != nil {
		applog.Danger("Cannot open menu file", err)
	}
	decoder := json.NewDecoder(file)
	sidebar = &sidebarMenu{}
	err = decoder.Decode(sidebar)
	if err != nil {
		applog.Danger("loadMenu malformed JSON", err)
	}
}

func getSidebarItems(_items *[]sidebarItem, menuitem string, _ulcls string) (tags string) {
	var cls string
	var has bool
	var act bool
	var cl string

	if len(_ulcls) > 0 {
		cls = " class=\"" + _ulcls + "\""
	}

	tags = "<ul" + cls + ">"
	for _, i := range *_items {
		cls = ""
		cl = ""
		act = false
		has = i.Items != nil && len(*i.Items) > 0
		if len(i.Alias) > 0 && i.Alias == menuitem {
			act = true
			tag_item = &i
		}

		if has {
			cl = " class=\"has-submenu init\""
		}

		tags += "<li" + cl + ">"

		if act {
			cls = " active"
		}

		if has {
			tags += fmt.Sprintf("<div class=\"sb-item%s clearfix\" >", cls)
		} else {
			hr := ro(i.Alias)

			tags += fmt.Sprintf("<a href=\"%s\" class=\"sb-item%s clearfix\">", hr, cls)
		}

		tags += "<span class=\"icon-admin\">"
		tags += "<i class=\"" + i.Icon + "\"></i>"
		tags += "</span>"
		tags += "<span class=\"name\">" + i.Label + "</span>"

		if has {
			tags += "</div>"
		} else {
			tags += "</a>"
		}

		if has {
			tags += getSidebarItems(i.Items, menuitem, "")
		}

		tags += "</li>"
	}
	tags += "</ul>"

	return
}

// GetSidebarHTML returns html for sidebar
func getSidebarHTML(menuitem string) (tags string) {
	for _, v := range *sidebar.Sidebar {
		tags += "<div class=\"title-admin\">" + v.Title + "</div>"
		if v.Items != nil && len(*v.Items) > 0 {
			tags += getSidebarItems(v.Items, menuitem, "zlev")
		}
	}
	return
}

// CreateSidebar creates sidebar menu
func CreateSidebarWeb(menuitem string) (template.HTML, *sidebarItem) {
	tag_item = nil
	loadMenu()
	sb_html := template.HTML(getSidebarHTML(menuitem))
	return sb_html, tag_item
}
