package mapper

import (
	domainUser "github.com/prr133f/avito-backend-intership-2025/internal/domain/user"
	"github.com/prr133f/avito-backend-intership-2025/internal/rest/dto"
)

func UserDTOToDomain(u dto.User) (domainUser.User, error) {
	return domainUser.User{
		Id:       u.Id,
		Name:     u.Name,
		IsActive: u.IsActive,
	}, nil
}

func UserDomainToDTO(u domainUser.User) dto.User {
	return dto.User{
		Id:       u.Id,
		Name:     u.Name,
		IsActive: u.IsActive,
		TeamName: u.TeamName,
	}
}
