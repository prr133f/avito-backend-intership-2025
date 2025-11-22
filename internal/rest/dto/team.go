package dto

type Team struct {
	Name  string `json:"team_name"`
	Users []User `json:"members"`
}
