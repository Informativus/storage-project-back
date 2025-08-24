package roles_model

type Role int8

const (
	Admin Role = iota + 1
	User
	Owner
	Reader
	Editor
)

var Roles = map[Role]string{
	Admin:  "admin",
	User:   "user",
	Owner:  "owner",
	Reader: "reader",
	Editor: "editor",
}
