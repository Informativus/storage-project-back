package roles_model

type Role int8

const (
	Admin Role = iota + 1
	User
	Owner
)

var Roles = map[Role]string{
	Admin: "admin",
	User:  "user",
	Owner: "owner",
}
