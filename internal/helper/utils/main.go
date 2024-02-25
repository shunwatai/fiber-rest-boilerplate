package utils

import (
	// "golang-api-starter/internal/config"
	"path"
	"path/filepath"
	"runtime"
)

// var cfg = config.Cfg

// ToPtr uses to return the pointer of the value without using one more line to declare a variable
// e.g.: helper.ToPtr("some string") returns the address of "some string"
func ToPtr[T any](v T) *T {
	return &v
}

// RootDir get the project base path
// ref: https://stackoverflow.com/a/58294680
func RootDir(level int) string {
	parentPath := ""
	for p := 0; p < level; p++ {
		parentPath += "../"
	}
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b), parentPath)
	return filepath.Dir(d)
}
