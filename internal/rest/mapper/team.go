package mapper

import (
	"fmt"

	domainTeam "github.com/prr133f/avito-backend-intership-2025/internal/domain/team"
	domainUser "github.com/prr133f/avito-backend-intership-2025/internal/domain/user"
	"github.com/prr133f/avito-backend-intership-2025/internal/rest/dto"
)

func TeamDTOToDomain(d dto.Team) (domainTeam.Team, error) {
	var out domainTeam.Team

	out.Name = d.Name

	users := make([]domainUser.User, 0, len(d.Users))
	for _, u := range d.Users {
		user, err := UserDTOToDomain(u)
		if err != nil {
			return out, fmt.Errorf("invalid user %q: %w", u.Id, err)
		}
		users = append(users, user)
	}
	out.Users = users

	return out, nil
}

func TeamDomainToDTO(t domainTeam.Team) dto.Team {
	users := make([]dto.User, 0, len(t.Users))
	for _, u := range t.Users {
		users = append(users, UserDomainToDTO(u))
	}
	return dto.Team{
		Name:  t.Name,
		Users: users,
	}
}
