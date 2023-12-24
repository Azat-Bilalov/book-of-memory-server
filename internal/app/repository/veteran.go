package repository

import (
	"github.com/Azat-Bilalov/book-of-memory-server/internal/app/ds"
	"gorm.io/gorm"
)

type InterfaceVeteranRepository interface {
	Store(veteran *ds.Veteran) (*ds.Veteran, error)
	FindAll(name string) ([]*ds.Veteran, error)
	FindByUUID(uuid string) (*ds.Veteran, error)
	UpdateByUUID(veteran *ds.Veteran) (*ds.Veteran, error)
	DeleteByUUID(uuid string) error
}

type VeteranRepository struct {
	db *gorm.DB
}

func NewVeteranRepository(db *gorm.DB) InterfaceVeteranRepository {
	return &VeteranRepository{db}
}

func (r *VeteranRepository) Store(veteran *ds.Veteran) (*ds.Veteran, error) {
	if err := r.db.Create(veteran).Error; err != nil {
		return nil, err
	}
	return veteran, nil
}

func (r *VeteranRepository) FindAll(name string) ([]*ds.Veteran, error) {
	veterans := make([]*ds.Veteran, 0)
	err := r.db.
		Find(&veterans).
		Where("first_name LIKE ? OR last_name LIKE ? OR patronymic LIKE ?", "%"+name+"%", "%"+name+"%", "%"+name+"%").
		Error
	if err != nil {
		return nil, err
	}
	return veterans, nil
}

func (r *VeteranRepository) FindByUUID(uuid string) (*ds.Veteran, error) {
	veteran := &ds.Veteran{}
	err := r.db.First(veteran, "Veteran_id = ?", uuid).Error
	if err != nil {
		return nil, err
	}
	return veteran, nil
}

func (r *VeteranRepository) UpdateByUUID(veteran *ds.Veteran) (*ds.Veteran, error) {
	if err := r.db.Save(veteran).Error; err != nil {
		return nil, err
	}
	return veteran, nil
}

func (r *VeteranRepository) DeleteByUUID(uuid string) error {
	if err := r.db.Delete(&ds.Veteran{}, "Veteran_id = ?", uuid).Error; err != nil {
		return err
	}
	return nil
}
