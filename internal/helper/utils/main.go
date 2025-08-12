package utils

import (
	"crypto/rand"
	"fmt"
	"net/url"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// ToPtr uses to return the pointer of the value without using one more line to declare a variable
// e.g.: helper.ToPtr("some string") returns the address of "some string"
func ToPtr[T any](v T) *T {
	return &v
}

// Deref get the value of pointer. I use it as a helper func in html template.
func Deref[T any](v *T) T {
	return *v
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

// GetRandString generate random string by given length
// ref: https://gist.github.com/arxdsilva/8caeca47b126a290c4562a25464895e8
func GetRandString(length int) string {
	if length < 1 {
		length = 1
	}
	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

// HashPassword hash the given plain string by bcrypt
func HashPassword(pwd string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		fmt.Printf("failed to bcrypt.GenerateFromPassword, err: %+v", err)
	}

	return string(hash)
}

// GetDomain extracts the primary domain from a URL
func GetDomain(urlStr string) string {
	// Parse the URL
	parsedURL, err := url.Parse(urlStr)
	// if err != nil {
	// 	return ""
	// }

	host := parsedURL.Hostname()
	if err != nil || len(host) == 0 {
		return urlStr
	}

	// Split the host into parts
	parts := strings.Split(host, ".")
	n := len(parts)

	// Return last two parts for main domain (e.g., example.com)
	if n > 1 {
		return strings.Join(parts[n-2:], ".")
	}

	return host // for cases like localhost
}
