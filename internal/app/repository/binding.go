package repository

import (
	"time"

	"github.com/Azat-Bilalov/book-of-memory-server/internal/app/ds"
	"gorm.io/gorm"
)

type InterfaceBindingRepository interface {
	Store(binding *ds.Binding) (*ds.Binding, error)
	FindAll(status string, timeFrom *time.Time, timeTo *time.Time) ([]*ds.Binding, error)
	FindAllByUserID(userID string, status string, timeFrom *time.Time, timeTo *time.Time) ([]*ds.Binding, error)
	FindAllByVeteranID(veteranID string, status string, dateFrom string, dateTo string) ([]*ds.Binding, error)
	FindByUUID(uuid string) (*ds.Binding, error)
	FindLastEnteredBindingByUserID(userID string) (*ds.Binding, error)
	UpdateByUUID(binding *ds.Binding) (*ds.Binding, error)
	DeleteByUUID(uuid string) error
}

type BindingRepository struct {
	db *gorm.DB
}

func NewBindingRepository(db *gorm.DB) InterfaceBindingRepository {
	return &BindingRepository{db}
}

func (r *BindingRepository) Store(binding *ds.Binding) (*ds.Binding, error) {
	if err := r.db.Create(binding).Error; err != nil {
		return nil, err
	}

	return binding, nil
}

func (r *BindingRepository) FindAll(status string, timeFrom *time.Time, timeTo *time.Time) ([]*ds.Binding, error) {
	bindings := make([]*ds.Binding, 0)
	if timeFrom == nil && timeTo == nil {
		err := r.db.Find(&bindings, "? = '' OR status = ?", status, status).Error
		if err != nil {
			return nil, err
		}
		return bindings, nil
	}
	query := r.db.
		Table("bindings").
		Where("? = '' OR status = ?", status, status).
		Where("formatted_at >= ?", timeFrom).
		Where("formatted_at <= ?", timeTo).
		Order("created_at DESC")
	if err := query.Find(&bindings).Error; err != nil {
		return nil, err
	}
	return bindings, nil
}

func (r *BindingRepository) FindAllByUserID(userID string, status string, timeFrom *time.Time, timeTo *time.Time) ([]*ds.Binding, error) {
	bindings := make([]*ds.Binding, 0)
	if timeFrom == nil && timeTo == nil {
		err := r.db.Find(&bindings, "user_id = ? AND (? = '' OR status = ?)", userID, status, status).Error
		if err != nil {
			return nil, err
		}
		return bindings, nil
	}
	query := r.db.Table("bindings").
		Where("user_id = ?", userID).
		Where("? = '' OR status = ?", status, status).
		Where("formatted_at >= ?", timeFrom).
		Where("formatted_at <= ?", timeTo).
		Order("created_at DESC")
	if err := query.Find(&bindings).Error; err != nil {
		return nil, err
	}
	return bindings, nil
}

func (r *BindingRepository) FindAllByVeteranID(veteranID string, status string, dateFrom string, dateTo string) ([]*ds.Binding, error) {
	bindings := make([]*ds.Binding, 0)
	err := r.db.Find(&bindings, "veteran_id = ?", veteranID).Error
	if err != nil {
		return nil, err
	}
	return bindings, nil
}

func (r *BindingRepository) FindByUUID(uuid string) (*ds.Binding, error) {
	binding := &ds.Binding{}
	err := r.db.
		Preload("Documents").
		Preload("Veteran").
		Preload("Moderator", func(db *gorm.DB) *gorm.DB {
			return db.Select("user_id, first_name, last_name, email")
		}).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("user_id, first_name, last_name, email")
		}).
		First(binding, "binding_id = ?", uuid).
		Error
	if err != nil {
		return nil, err
	}
	return binding, nil
}

func (r *BindingRepository) FindLastEnteredBindingByUserID(userID string) (*ds.Binding, error) {
	binding := &ds.Binding{}
	err := r.db.First(binding, "user_id = ? AND status = ?", userID, ds.BINDING_STATUS_ENTERED).Error
	if err != nil {
		return nil, err
	}
	return binding, nil
}

func (r *BindingRepository) UpdateByUUID(binding *ds.Binding) (*ds.Binding, error) {
	if err := r.db.Save(binding).Error; err != nil {
		return nil, err
	}
	return binding, nil
}

func (r *BindingRepository) DeleteByUUID(uuid string) error {
	binding, err := r.FindByUUID(uuid)
	if err != nil {
		return err
	}
	binding.Status = ds.BINDING_STATUS_DELETED
	_, err = r.UpdateByUUID(binding)
	return err
}
