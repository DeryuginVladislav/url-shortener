package user

import "url_shortener/pkg/db"

type UserRepository struct {
	Database *db.Db
}

func NewUserRepository(database *db.Db) *UserRepository {
	return &UserRepository{
		Database: database,
	}
}

func (r *UserRepository) Create(user *User) (*User, error) {
	result := r.Database.DB.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}
func (r *UserRepository) FindByEmail(email string) (*User, error) {
	var u User
	result := r.Database.DB.Where("email=?", email).First(&u)
	if result.Error != nil {
		return nil, result.Error
	}
	return &u, nil
}
