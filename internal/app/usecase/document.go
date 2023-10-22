package usecase

import (
	"errors"

	"github.com/Azat-Bilalov/book-of-memory-server/internal/app/ds"
	"github.com/Azat-Bilalov/book-of-memory-server/internal/app/repository"
	"gorm.io/gorm"
)

type InterfaceDocumentUsecase interface {
	CreateDocument(document ds.DocumentCreateRequest) (*ds.Document, error)
	FindActiveDocuments(title string) ([]*ds.Document, error)
	FindActiveDocumentByUUID(uuid string) (*ds.Document, error)
	UpdateDocumentByUUID(uuid string, document ds.DocumentUpdateRequest) (*ds.Document, error)
	DeleteDocumentByUUID(uuid string) error
	FindDocumentInBinding(bindingID string, documentID string) (*ds.DocBinding, error)
	AddDocumentToBindingByUUID(documentID string, userID string, docBinding ds.DocBindingRequest) error
	RemoveDocumentFromBindingByUUID(documentID string, userID string) error
}

type DocumentUsecase struct {
	documentRepository   repository.InterfaceDocumentRepository
	bindingRepository    repository.InterfaceBindingRepository
	docBindingRepository repository.InterfaceDocBindingRepository
}

func NewDocumentUsecase(
	documentRepository repository.InterfaceDocumentRepository,
	bindingRepository repository.InterfaceBindingRepository,
	docBindingRepository repository.InterfaceDocBindingRepository,
) *DocumentUsecase {
	return &DocumentUsecase{documentRepository, bindingRepository, docBindingRepository}
}

func (u *DocumentUsecase) CreateDocument(document ds.DocumentCreateRequest) (*ds.Document, error) {
	d := &ds.Document{
		Title:       document.Title,
		Description: document.Description,
		Image_url:   document.Image_url,
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

func (u *DocumentUsecase) UpdateDocumentByUUID(uuid string, document ds.DocumentUpdateRequest) (*ds.Document, error) {
	d := &ds.Document{
		Document_id: uuid,
		Title:       document.Title,
		Description: document.Description,
		Image_url:   document.Image_url,
		Status:      ds.DOCUMENT_STATUS_ACTIVE,
	}

	return u.documentRepository.UpdateByUUID(d)
}

func (u *DocumentUsecase) DeleteDocumentByUUID(uuid string) error {
	document, err := u.documentRepository.FindByUUID(uuid)
	if err != nil {
		return err
	}
	document.Status = ds.DOCUMENT_STATUS_DELETED
	_, err = u.documentRepository.UpdateByUUID(document)
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
		binding = &ds.Binding{
			UserID:      userID,
			ModeratorID: "e35910a0-87b0-453b-b582-1460bc01d23c", // todo: remove hardcode
			Status:      ds.BINDING_STATUS_ENTERED,
		}
		binding, err = u.bindingRepository.Store(binding)
		if err != nil {
			return err
		}
	}
	if err != nil {
		return err
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
