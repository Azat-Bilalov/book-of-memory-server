package usecase

import (
	"errors"
	"time"

	"github.com/Azat-Bilalov/book-of-memory-server/internal/app/ds"
	"github.com/Azat-Bilalov/book-of-memory-server/internal/app/repository"
	"gorm.io/gorm"
)

const timeFormat = "2006-01-02 15:04:05"

type InterfaceBindingUsecase interface {
	FindBindingsByUserID(userID string, status string, dateFrom string, dateTo string) ([]*ds.Binding, error)
	FindBindingByVeteranID(veteranID string, status string, dateFrom string, dateTo string) ([]*ds.Binding, error)
	FindBindingByUUID(uuid string) (*ds.Binding, error)
	UpdateBindingByUUID(uuid string, binding ds.BindingUpdateRequest) (*ds.Binding, error)
	SubmitBindingByUUID(uuid string) (*ds.Binding, error)
	AcceptBindingByUUID(uuid string) (*ds.Binding, error)
	RejectBindingByUUID(uuid string) (*ds.Binding, error)
	DeleteBindingByUUID(uuid string) error
}

type BindingUsecase struct {
	bindingRepository repository.InterfaceBindingRepository
	userRepository    repository.InterfaceUserRepository
	veteranRepository repository.InterfaceVeteranRepository
}

func NewBindingUsecase(
	bindingRepository repository.InterfaceBindingRepository,
	userRepository repository.InterfaceUserRepository,
	veteranRepository repository.InterfaceVeteranRepository,
) *BindingUsecase {
	return &BindingUsecase{bindingRepository, userRepository, veteranRepository}
}

func (u *BindingUsecase) FindBindingsByUserID(userID string, status string, dateFrom string, dateTo string) ([]*ds.Binding, error) {
	if userID == "" {
		return nil, errors.New("идентификатор пользователя пустой")
	}
	user, err := u.userRepository.FindByUUID(userID)
	if err == gorm.ErrRecordNotFound {
		return nil, errors.New("пользователь не найден")
	}
	if err != nil {
		return nil, err
	}

	if dateFrom == "" && dateTo == "" {
		if user.Role == ds.USER_ROLE_MODERATOR {
			return u.bindingRepository.FindAll(status, nil, nil)
		}
		return u.bindingRepository.FindAllByUserID(userID, status, nil, nil)
	}

	timeFrom, err := time.Parse(timeFormat, dateFrom)
	if err != nil {
		timeFrom = time.Unix(0, 0)
	}
	timeTo, err := time.Parse(timeFormat, dateTo)
	if err != nil {
		timeTo = time.Now()
	}
	if user.Role == ds.USER_ROLE_MODERATOR {
		return u.bindingRepository.FindAll(status, &timeFrom, &timeTo)
	}
	return u.bindingRepository.FindAllByUserID(userID, status, &timeFrom, &timeTo)
}

func (u *BindingUsecase) FindBindingByVeteranID(veteranID string, status string, dateFrom string, dateTo string) ([]*ds.Binding, error) {
	return u.bindingRepository.FindAllByVeteranID(veteranID, status, dateFrom, dateTo)
}

func (u *BindingUsecase) FindBindingByUUID(uuid string) (*ds.Binding, error) {
	binding, err := u.bindingRepository.FindByUUID(uuid)
	if err == gorm.ErrRecordNotFound {
		return nil, errors.New("заявка не найдена")
	}
	if err != nil {
		return nil, err
	}
	return binding, nil
}

func (u *BindingUsecase) UpdateBindingByUUID(uuid string, binding ds.BindingUpdateRequest) (*ds.Binding, error) {
	b, err := u.bindingRepository.FindByUUID(uuid)
	if err != nil {
		return nil, err
	}
	v, err := u.veteranRepository.FindByUUID(binding.VeteranID)
	if err == gorm.ErrRecordNotFound {
		return nil, errors.New("ветеран не найден")
	}
	if err != nil {
		return nil, err
	}
	b.VeteranID = &v.Veteran_id
	return u.bindingRepository.UpdateByUUID(b)
}

func (u *BindingUsecase) SubmitBindingByUUID(uuid string) (*ds.Binding, error) {
	b, err := u.bindingRepository.FindByUUID(uuid)
	if err != nil {
		return nil, err
	}
	b.Status = ds.BINDING_STATUS_IN_PROGRESS
	now := time.Now()
	b.FormattedAt = &now
	return u.bindingRepository.UpdateByUUID(b)
}

func (u *BindingUsecase) AcceptBindingByUUID(uuid string) (*ds.Binding, error) {
	b, err := u.bindingRepository.FindByUUID(uuid)
	if err != nil {
		return nil, err
	}
	b.Status = ds.BINDING_STATUS_COMPLETED
	now := time.Now()
	b.EndedAt = &now
	return u.bindingRepository.UpdateByUUID(b)
}

func (u *BindingUsecase) RejectBindingByUUID(uuid string) (*ds.Binding, error) {
	b, err := u.bindingRepository.FindByUUID(uuid)
	if err != nil {
		return nil, err
	}
	b.Status = ds.BINDING_STATUS_CANCELED
	now := time.Now()
	b.EndedAt = &now
	return u.bindingRepository.UpdateByUUID(b)
}

func (u *BindingUsecase) DeleteBindingByUUID(uuid string) error {
	b, err := u.bindingRepository.FindByUUID(uuid)
	if err != nil {
		return err
	}
	b.Status = ds.BINDING_STATUS_DELETED
	now := time.Now()
	b.EndedAt = &now
	_, err = u.bindingRepository.UpdateByUUID(b)
	return err
}
