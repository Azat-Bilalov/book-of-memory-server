package usecase

import (
	"github.com/Azat-Bilalov/book-of-memory-server/internal/app/ds"
	"github.com/Azat-Bilalov/book-of-memory-server/internal/app/repository"
)

type InterfaceDocumentUsecase interface {
	CreateDocument(document ds.DocumentCreateRequest) (*ds.Document, error)
	FindActiveDocuments() ([]*ds.Document, error)
	FindActiveDocumentByUUID(uuid string) (*ds.Document, error)
	UpdateDocumentByUUID(uuid string, document ds.DocumentUpdateRequest) (*ds.Document, error)
	DeleteDocumentByUUID(uuid string) error
}

type DocumentUsecase struct {
	documentRepository repository.InterfaceDocumentRepository
}

func NewDocumentUsecase(documentRepository repository.InterfaceDocumentRepository) *DocumentUsecase {
	return &DocumentUsecase{documentRepository}
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

func (u *DocumentUsecase) FindActiveDocuments() ([]*ds.Document, error) {
	return u.documentRepository.FindWithStatus(ds.DOCUMENT_STATUS_ACTIVE)
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
