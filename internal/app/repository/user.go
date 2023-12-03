package repository

import (
	"github.com/Azat-Bilalov/book-of-memory-server/internal/app/ds"
	"gorm.io/gorm"
)

type InterfaceUserRepository interface {
	Store(user *ds.User) (*ds.User, error)
	FindAllUsers() ([]*ds.User, error)
	FindAllModerators() ([]*ds.User, error)
	FindByUUID(uuid string) (*ds.User, error)
	FindByEmail(email string) (*ds.User, error)
	UpdateByUUID(user *ds.User) (*ds.User, error)
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) InterfaceUserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) Store(user *ds.User) (*ds.User, error) {
	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) FindAllUsers() ([]*ds.User, error) {
	users := make([]*ds.User, 0)
	err := r.db.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) FindAllModerators() ([]*ds.User, error) {
	users := make([]*ds.User, 0)
	err := r.db.Find(&users, "role = ?", ds.USER_ROLE_MODERATOR).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) FindByUUID(uuid string) (*ds.User, error) {
	user := &ds.User{}
	err := r.db.First(user, "user_id = ?", uuid).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) FindByEmail(email string) (*ds.User, error) {
	user := &ds.User{}
	err := r.db.First(user, "email = ?", email).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) UpdateByUUID(user *ds.User) (*ds.User, error) {
	err := r.db.Save(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}
