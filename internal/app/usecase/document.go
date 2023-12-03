package usecase

import (
	"errors"
	"log"
	"math/rand"

	"github.com/Azat-Bilalov/book-of-memory-server/internal/app/config"
	"github.com/Azat-Bilalov/book-of-memory-server/internal/app/ds"
	"github.com/Azat-Bilalov/book-of-memory-server/internal/app/repository"
	"gorm.io/gorm"
)

const (
	BUSKET = "documents"
)

type InterfaceDocumentUsecase interface {
	CreateDocument(title string, description string, image_url string) (*ds.Document, error)
	FindActiveDocuments(title string) ([]*ds.Document, error)
	FindActiveDocumentByUUID(uuid string) (*ds.Document, error)
	UpdateDocumentByUUID(uuid string, document ds.DocumentRequest) (*ds.Document, error)
	DeleteDocumentByUUID(uuid string) error
	AddDocumentToBindingByUUID(documentID string, userID string, docBinding ds.DocBinding) error
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
		Image_url:   BUSKET + "/" + document.Image.Filename,
		Status:      ds.DOCUMENT_STATUS_ACTIVE,
	}

	return u.documentRepository.Store(d)
}

func (u *DocumentUsecase) FindActiveDocuments(title string, userID string) (*ds.DocumentListResponse, error) {
	log.Println("FindActiveDocuments")
	documents, err := u.documentRepository.FindAll(ds.DOCUMENT_STATUS_ACTIVE, title)
	log.Println(documents)
	if err != nil {
		return nil, err
	}
	documentList := &ds.DocumentListResponse{}
	documentList.Documents = documents
	documentList.EnteredBinding_id = nil
	if userID != "" {
		enteredBinding, err := u.bindingRepository.FindLastEnteredBindingByUserID(userID)
		if err == nil {
			documentList.EnteredBinding_id = &enteredBinding.Binding_id
		}
	}
	return documentList, nil
}

func (u *DocumentUsecase) FindActiveDocumentByUUID(uuid string) (*ds.Document, error) {
	document, err := u.documentRepository.FindByUUID(uuid)
	if err != nil {
		return nil, err
	}
	return document, nil
}

func (u *DocumentUsecase) UpdateDocumentByUUID(uuid string, document ds.DocumentRequest) (*ds.Document, error) {
	d, err := u.documentRepository.FindByUUID(uuid)
	if err != nil || d.Status != ds.DOCUMENT_STATUS_ACTIVE {
		return nil, errors.New("document not found")
	}
	countDocumentsWithImageUrl, err := u.documentRepository.CountWithImageUrl(d.Image_url)
	if err != nil {
		return nil, err
	}
	if countDocumentsWithImageUrl == 1 {
		imageName := d.Image_url[len(BUSKET)+1:]
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
		Image_url:   BUSKET + "/" + document.Image.Filename,
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
	imageName := document.Image_url[len(BUSKET)+1:]
	err = config.DeleteObject(BUSKET, imageName)
	if err != nil {
		return err
	}
	return nil
}

func (u *DocumentUsecase) AddDocumentToBindingByUUID(documentID string, userID string, docBindingRequest ds.DocBinding) error {
	binding, err := u.bindingRepository.FindLastEnteredBindingByUserID(userID)
	if err == gorm.ErrRecordNotFound {
		moderators, err := u.userRepository.FindAllModerators()
		if err != nil {
			return err
		}
		if len(moderators) == 0 {
			return errors.New("модераторы не найдены")
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
	_, err = u.docBindingRepository.Find(documentID, binding.Binding_id)
	if err != gorm.ErrRecordNotFound {
		return errors.New("документ уже добавлен в заявку")
	}
	docBinding := &ds.DocBinding{
		Binding_id:  binding.Binding_id,
		Document_id: documentID,
	}
	_, err = u.docBindingRepository.Store(docBinding)
	return err
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
