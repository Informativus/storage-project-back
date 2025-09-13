package roles_model

const (
	TableName = "roles"
)

type RolesModel struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Note string `json:"note"`
}
