package helper

import (
	"path"
	"path/filepath"
	"runtime"
)

// GetIdsMapCondition() return map[string]interface{} sth like map[string]interface{}{"id":[]string{x,y,z}}
// which can use for Get() to retrieve records by id(s)
// param "keyId" can be nil for default "id" in db, it can be other column(foreign key) like "todo_id"
func GetIdsMapCondition(keyId *string, ids []string) map[string]interface{} {
	condition := map[string]interface{}{}

	if keyId == nil {
		if cfg.DbConf.Driver == "mongodb" {
			condition["_id"] = ids
		} else {
			condition["id"] = ids
		}
	} else {
		condition[*keyId] = ids
	}

	return condition
}

// ToPtr() uses to return the pointer of the value without using one more line to declare a variable
// e.g.: helper.ToPtr("some string") returns the address of "some string"
func ToPtr[T any](v T) *T {
	return &v
}

// RootDir() get the project base path
// ref: https://stackoverflow.com/a/58294680
func RootDir() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b), "../")
	return filepath.Dir(d)
}
