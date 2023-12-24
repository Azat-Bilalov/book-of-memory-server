package usecase

import (
	"github.com/Azat-Bilalov/book-of-memory-server/internal/app/config"
	"github.com/Azat-Bilalov/book-of-memory-server/internal/app/ds"
	"github.com/Azat-Bilalov/book-of-memory-server/internal/app/repository"
)

type InterfaceVeteranUsecase interface {
	CreateVeteran(veteran ds.VeteranRequest) (*ds.Veteran, error)
	FindActiveVeterans(title string) ([]*ds.Veteran, error)
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

func (u *VeteranUsecase) FindVeterans(title string) ([]*ds.Veteran, error) {
	veterans, err := u.veteranRepository.FindAll(title)
	if err != nil {
		return nil, err
	}
	return veterans, nil
}
