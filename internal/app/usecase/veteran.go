package usecase

import (
	"log"

	"github.com/Azat-Bilalov/book-of-memory-server/internal/app/config"
	"github.com/Azat-Bilalov/book-of-memory-server/internal/app/ds"
	"github.com/Azat-Bilalov/book-of-memory-server/internal/app/repository"
)

type InterfaceVeteranUsecase interface {
	CreateVeteran(veteran ds.VeteranRequest) (*ds.Veteran, error)
	FindActiveVeterans(title string) ([]*ds.Veteran, error)
	FindVeteranByUUID(uuid string) (*ds.Veteran, error)
	UpdateVeteranByUUID(uuid string, veteran *ds.VeteranRequest) (*ds.Veteran, error)
	DeleteVeteranByUUID(uuid string) error
}

type VeteranUsecase struct {
	veteranRepository repository.InterfaceVeteranRepository
}

func NewVeteranUsecase(veteranRepository repository.InterfaceVeteranRepository) *VeteranUsecase {
	return &VeteranUsecase{veteranRepository}
}

func (u *VeteranUsecase) CreateVeteran(veteran ds.VeteranRequest) (*ds.Veteran, error) {
	file, err := veteran.Image.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()
	err = config.UploadObject(BUSKET, veteran.Image.Filename, file, veteran.Image.Size, veteran.Image.Header.Get("Content-Type"))
	if err != nil {
		return nil, err
	}

	v := &ds.Veteran{
		FirstName:  veteran.FirstName,
		LastName:   veteran.LastName,
		Patronymic: veteran.Patronymic,
		BirthDate:  veteran.BirthDate,
		ImageUrl:   "veterans/" + veteran.Image.Filename,
	}

	return u.veteranRepository.Store(v)
}

func (u *VeteranUsecase) FindVeterans(name string) ([]*ds.Veteran, error) {
	veterans, err := u.veteranRepository.FindAll(name)
	if err != nil {
		return nil, err
	}
	return veterans, nil
}

func (u *VeteranUsecase) FindVeteranByUUID(uuid string) (*ds.Veteran, error) {
	veteran, err := u.veteranRepository.FindByUUID(uuid)
	if err != nil {
		return nil, err
	}
	return veteran, nil
}

func (u *VeteranUsecase) UpdateVeteranByUUID(uuid string, veteran *ds.VeteranRequest) (*ds.Veteran, error) {
	file, err := veteran.Image.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()
	err = config.UploadObject(BUSKET, veteran.Image.Filename, file, veteran.Image.Size, veteran.Image.Header.Get("Content-Type"))
	if err != nil {
		return nil, err
	}

	// парсинг даты
	log.Println(veteran.BirthDate)

	v := &ds.Veteran{
		Veteran_id: uuid,
		FirstName:  veteran.FirstName,
		LastName:   veteran.LastName,
		Patronymic: veteran.Patronymic,
		BirthDate:  veteran.BirthDate,
		ImageUrl:   "veterans/" + veteran.Image.Filename,
	}

	return u.veteranRepository.UpdateByUUID(v)
}

func (u *VeteranUsecase) DeleteVeteranByUUID(uuid string) error {
	err := u.veteranRepository.DeleteByUUID(uuid)
	if err != nil {
		return err
	}
	return nil
}
