package helper

var MethodToPermType = map[string]string{
	"GET":    "read",
	"POST":   "add",
	"PATCH":  "edit",
	"PUT":    "edit",
	"DELETE": "delete",
}
