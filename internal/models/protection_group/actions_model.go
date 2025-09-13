package protection_group

const (
	ActionsName = "actions"
)

type ActionsModel struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Descritpion string `json:"description"`
}
