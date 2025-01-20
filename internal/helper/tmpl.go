package helper

import (
	"golang-api-starter/internal/helper/utils"
	"html/template"
	"strconv"
	"strings"
	"time"
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
		// IsMongo check whether cfg.DbConf.Driver is mongo
		"IsMongo": func() bool {
			return cfg.DbConf.Driver == "mongodb"
		},
		// GetActionByMethod for log list page to show the (readable)request aciont name for user
		"GetActionByMethod": func(method string) string {
			action, ok := MethodToPermType[method]
			if !ok || len(action) == 0 {
				return ""
			}
			return strings.ToUpper(action[0:1]) + action[1:]
		},
		// GetSucceedFailedByCode return succeed / failed by status code 2xx / 4xx
		"GetSucceedFailedByCode": func(code int64) string {
			if code >= 200 && code < 400 {
				return "Succeed"
			}
			if code >= 400 && code <= 500 {
				return "Failed"
			}
			return ""
		},
		// GetDuration converts duration to millisecond for display
		"GetDuration": func(duration int64) time.Duration {
			return time.Duration(duration) * time.Millisecond
		},
	}
}
