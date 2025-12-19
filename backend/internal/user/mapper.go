package user

import (
	"app/domain"
	"app/ent"
)

func ToUserResponse(u *ent.User) domain.UserResponse {
	return domain.UserResponse{
		ID:       u.ID,
		Email:    u.Email,
		Username: u.Username,
	}
}
