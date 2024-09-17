package persistence

import (
	"github.com/OmarBader7/web-service-jayeek/domain/entity"
	"gorm.io/gorm"
)

// IdentityDocumentRepository implements the repository.IdentityDocumentRepository interface
type IdentityDocumentRepository struct {
	// db is a pointer to the GORM DB instance
	db *gorm.DB
}

// NewIdentityDocumentRepository creates a new instance of the IdentityDocumentRepository
func NewIdentityDocumentRepository(db *gorm.DB) *IdentityDocumentRepository {
	return &IdentityDocumentRepository{db: db}
}

func (r *IdentityDocumentRepository) CreateIdentityDocument(identityDocument *entity.IdentityDocument) (*entity.IdentityDocument, error) {
	if err := r.db.Debug().Model(&identityDocument).Create(&identityDocument).Error; err != nil {
		return nil, err
	}
	if err := r.db.Debug().Model(&identityDocument).Take(&identityDocument).Error; err != nil {
		return nil, err
	}
	return identityDocument, nil
}
