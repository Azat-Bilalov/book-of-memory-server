package usecase

import (
	"errors"
	"math/rand"

	"github.com/Azat-Bilalov/book-of-memory-server/internal/app/config"
	"github.com/Azat-Bilalov/book-of-memory-server/internal/app/ds"
	"github.com/Azat-Bilalov/book-of-memory-server/internal/app/repository"
	"gorm.io/gorm"
)

const (
	BUSKET     = "documents"
	STATIC_URL = "http://localhost:8080/files/" + BUSKET
)

type InterfaceDocumentUsecase interface {
	CreateDocument(title string, description string, image_url string) (*ds.Document, error)
	FindActiveDocuments(title string) ([]*ds.Document, error)
	FindActiveDocumentByUUID(uuid string) (*ds.Document, error)
	UpdateDocumentByUUID(uuid string, document ds.DocumentRequest) (*ds.Document, error)
	DeleteDocumentByUUID(uuid string) error
	FindDocumentInBinding(bindingID string, documentID string) (*ds.DocBinding, error)
	AddDocumentToBindingByUUID(documentID string, userID string, docBinding ds.DocBindingRequest) error
	RemoveDocumentFromBindingByUUID(documentID string, userID string) error
}

type DocumentUsecase struct {
	documentRepository   repository.InterfaceDocumentRepository
	bindingRepository    repository.InterfaceBindingRepository
	docBindingRepository repository.InterfaceDocBindingRepository
	userRepository       repository.InterfaceUserRepository
}

func NewDocumentUsecase(
	documentRepository repository.InterfaceDocumentRepository,
	bindingRepository repository.InterfaceBindingRepository,
	docBindingRepository repository.InterfaceDocBindingRepository,
	userRepository repository.InterfaceUserRepository,
) *DocumentUsecase {
	return &DocumentUsecase{documentRepository, bindingRepository, docBindingRepository, userRepository}
}

func (u *DocumentUsecase) CreateDocument(document ds.DocumentRequest) (*ds.Document, error) {
	file, err := document.Image.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()
	err = config.UploadObject(BUSKET, document.Image.Filename, file, document.Image.Size, document.Image.Header.Get("Content-Type"))
	if err != nil {
		return nil, err
	}

	d := &ds.Document{
		Title:       document.Title,
		Description: document.Description,
		Image_url:   STATIC_URL + "/" + document.Image.Filename,
		Status:      ds.DOCUMENT_STATUS_ACTIVE,
	}

	return u.documentRepository.Store(d)
}

func (u *DocumentUsecase) FindActiveDocuments(title string) ([]*ds.Document, error) {
	return u.documentRepository.FindAll(ds.DOCUMENT_STATUS_ACTIVE, title)
}

func (u *DocumentUsecase) FindActiveDocumentByUUID(uuid string) (*ds.Document, error) {
	document, err := u.documentRepository.FindByUUID(uuid)
	if err != nil {
		return nil, err
	}
	if document.Status != ds.DOCUMENT_STATUS_ACTIVE {
		return nil, nil
	}
	return document, nil
}

func (u *DocumentUsecase) UpdateDocumentByUUID(uuid string, document ds.DocumentRequest) (*ds.Document, error) {
	d, err := u.documentRepository.FindByUUID(uuid)
	if err != nil {
		return nil, err
	}
	if d.Status != ds.DOCUMENT_STATUS_ACTIVE {
		return nil, errors.New("document not found")
	}
	countDocumentsWithImageUrl, err := u.documentRepository.CountWithImageUrl(d.Image_url)
	if err != nil {
		return nil, err
	}
	if countDocumentsWithImageUrl == 1 {
		imageName := d.Image_url[len(STATIC_URL)+1:]
		err = config.DeleteObject(BUSKET, imageName)
		if err != nil {
			return nil, err
		}
	}

	file, err := document.Image.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()
	err = config.UploadObject(BUSKET, document.Image.Filename, file, document.Image.Size, document.Image.Header.Get("Content-Type"))
	if err != nil {
		return nil, err
	}

	d = &ds.Document{
		Document_id: uuid,
		Title:       document.Title,
		Description: document.Description,
		Image_url:   STATIC_URL + "/" + document.Image.Filename,
		Status:      ds.DOCUMENT_STATUS_ACTIVE,
	}

	return u.documentRepository.UpdateByUUID(d)
}

func (u *DocumentUsecase) DeleteDocumentByUUID(uuid string) error {
	document, err := u.documentRepository.FindByUUID(uuid)
	if err != nil {
		return err
	}
	if document.Status != ds.DOCUMENT_STATUS_ACTIVE {
		return errors.New("document not found")
	}
	document.Status = ds.DOCUMENT_STATUS_DELETED
	_, err = u.documentRepository.UpdateByUUID(document)
	if err != nil {
		return err
	}
	countDocumentsWithImageUrl, err := u.documentRepository.CountWithImageUrl(document.Image_url)
	if err != nil {
		return err
	}
	if countDocumentsWithImageUrl > 1 {
		return nil
	}
	imageName := document.Image_url[len(STATIC_URL)+1:]
	err = config.DeleteObject(BUSKET, imageName)
	if err != nil {
		return err
	}
	return nil
}

func (u *DocumentUsecase) FindDocumentInBinding(bindingID string, documentID string) (*ds.DocBinding, error) {
	return u.docBindingRepository.Find(documentID, bindingID)
}

func (u *DocumentUsecase) AddDocumentToBindingByUUID(documentID string, userID string, docBindingRequest ds.DocBindingRequest) error {
	binding, err := u.bindingRepository.FindLastEnteredBindingByUserID(userID)
	if err == gorm.ErrRecordNotFound {
		moderators, err := u.userRepository.FindAllModerators()
		if err != nil {
			return err
		}
		if len(moderators) == 0 {
			return errors.New("moderators not found")
		}
		randomModerator := moderators[rand.Intn(len(moderators))]
		binding = &ds.Binding{
			UserID:      userID,
			ModeratorID: randomModerator.User_id,
			Status:      ds.BINDING_STATUS_ENTERED,
		}
		binding, err = u.bindingRepository.Store(binding)
		if err != nil {
			return err
		}
	}
	docBinding, err := u.docBindingRepository.Find(documentID, binding.Binding_id)
	if err == gorm.ErrRecordNotFound {
		docBinding = &ds.DocBinding{
			Binding_id:  binding.Binding_id,
			Document_id: documentID,
			File_url:    docBindingRequest.File_url,
			Info:        docBindingRequest.Info,
		}
		_, err = u.docBindingRepository.Store(docBinding)
		return err
	}
	if err != nil {
		return err
	}
	docBinding.File_url = docBindingRequest.File_url
	docBinding.Info = docBindingRequest.Info
	return u.docBindingRepository.Update(docBinding)
}

func (u *DocumentUsecase) RemoveDocumentFromBindingByUUID(documentID string, userID string) error {
	binding, err := u.bindingRepository.FindLastEnteredBindingByUserID(userID)
	if err != nil {
		return err
	}
	_, err = u.docBindingRepository.Find(documentID, binding.Binding_id)
	if err == gorm.ErrRecordNotFound {
		return errors.New("document not found in binding")
	}
	if err != nil {
		return err
	}
	return u.docBindingRepository.Delete(documentID, binding.Binding_id)
}
