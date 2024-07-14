package helper

import (
	"golang-api-starter/internal/helper/utils"
	"html/template"
	"strconv"
	"strings"
)

// TmplCustomFuncs contains list of  custom functions for go template
func TmplCustomFuncs() template.FuncMap {
	return template.FuncMap{
		// Iterate for looping by a given number, I use for rendering the dropdown of page nubmers
		"Iterate": func(totalPages *int64) []int64 {
			var i int64
			var page []int64
			for i = 0; i < (*totalPages); i++ {
				page = append(page, i+1)
			}
			return page
		},
		// Contains is checking the partial match of a given keyword
		"Contains": func(word, key string) bool {
			return strings.Contains(word, key)
		},
		// ShowSortingDirection shows the "arrow" icon at the table's header cell indicating the sorting order asc | desc
		"ShowSortingDirection": func(tableCellName string, orderBy map[string]string) string {
			if orderBy["key"] != tableCellName {
				return ""
			}
			if orderBy["by"] == "asc" {
				return "↑"
			} else {
				return "↓"
			}
		},
		// DerefBool get the value of *bool
		"DerefBool": utils.Deref[bool],
		// GetId get either int ID or mongo ID
		"GetId": func(mongoId *string, id *FlexInt) string {
			if mongoId == nil {
				return strconv.Itoa(int(*id))
			}
			return *mongoId
		},
		// GetIdKeyName return either _id or id for the html's attribute based on mongoId
		"GetIdKeyName": func() string {
			if cfg.DbConf.Driver == "mongodb" {
				return "_id"
			}
			return "id"
		},
	}
}
