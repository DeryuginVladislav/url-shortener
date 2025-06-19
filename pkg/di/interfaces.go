package di

import (
	"url_shortener/internal/user"
)

type IStatRepository interface {
	AddClick(linkId uint)
}

type IUserReposytory interface {
	Create(user *user.User) (*user.User, error)
	FindByEmail(email string) (*user.User, error)
}
