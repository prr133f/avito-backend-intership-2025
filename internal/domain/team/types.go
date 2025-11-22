package team

import "github.com/prr133f/avito-backend-intership-2025/internal/domain/user"

type Team struct {
	Name  string `db:"name"`
	Users []user.User
}
