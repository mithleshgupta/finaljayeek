package application

import (
	"github.com/OmarBader7/web-service-jayeek/domain/entity"
	"github.com/OmarBader7/web-service-jayeek/domain/repository"
)

// IdentityDocumentApplication struct is responsible for handling setting-related business logic
type IdentityDocumentApplication struct {
	identityDocumentRepo repository.IdentityDocumentRepository
}

var _ IdentityDocumentApplicationInterface = &IdentityDocumentApplication{}

// IdentityDocumentApplicationInterface defines the methods that IdentityDocumentApplication should implement
type IdentityDocumentApplicationInterface interface {
	CreateIdentityDocument(identityDocument *entity.IdentityDocument) (*entity.IdentityDocument, error)
}

func (a *IdentityDocumentApplication) CreateIdentityDocument(identityDocument *entity.IdentityDocument) (*entity.IdentityDocument, error) {
	return a.identityDocumentRepo.CreateIdentityDocument(identityDocument)
}
