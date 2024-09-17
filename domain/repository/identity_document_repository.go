package repository

import (
	"github.com/OmarBader7/web-service-jayeek/domain/entity"
)

// IdentityDocumentRepository defines the methods that a setting repository should implement
type IdentityDocumentRepository interface {
	CreateIdentityDocument(identityDocument *entity.IdentityDocument) (*entity.IdentityDocument, error)
}
