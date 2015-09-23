package user

type User struct {
	ID      string            `json:"user"`
	Profile map[string]string `json:"profile"`
}
