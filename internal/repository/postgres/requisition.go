package postgres

import (
	"FinMa/internal/domain"
	"FinMa/internal/repository"
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RequisitionRepository struct {
	db *gorm.DB
}

func NewRequisitionRepository(db *gorm.DB) *RequisitionRepository {
	return &RequisitionRepository{
		db: db,
	}
}

func (r *RequisitionRepository) Create(ctx context.Context, requisition *domain.Requisition) error {
	if err := r.db.WithContext(ctx).Create(requisition).Error; err != nil {
		return repository.NewRequisitionError("create", err, map[string]interface{}{
			"user_id":        requisition.UserID,
			"institution_id": requisition.InstitutionID,
		})
	}
	return nil
}
func (r *RequisitionRepository) Update(ctx context.Context, requisition *domain.Requisition) error {
	if err := r.db.WithContext(ctx).Where("id = ?", requisition.ID).Updates(requisition).Error; err != nil {
		return repository.NewRequisitionError("update", err, map[string]interface{}{
			"requisition_id": requisition.ID,
		})
	}
	return nil
}

func (r *RequisitionRepository) GetByReference(ctx context.Context, reference string) (*domain.Requisition, error) {
	var requisition domain.Requisition
	result := r.db.WithContext(ctx).Where("reference = ?", reference).First(&requisition)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, repository.NewRequisitionError("get_by_reference", repository.ErrRequisitionNotFound, map[string]interface{}{
				"reference": reference,
			})
		}
		return nil, repository.NewRequisitionError("get_by_reference", result.Error, map[string]interface{}{
			"reference": reference,
		})
	}
	return &requisition, nil
}

func (r *RequisitionRepository) GetByID(ctx context.Context, id string) (*domain.Requisition, error) {
	var requisition domain.Requisition
	result := r.db.WithContext(ctx).First(&requisition, "requisition_id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, repository.NewRequisitionError("get_by_id", repository.ErrRequisitionNotFound, map[string]interface{}{
				"requisition_id": id,
			})
		}
		return nil, repository.NewRequisitionError("get_by_id", result.Error, map[string]interface{}{
			"requisition_id": id,
		})
	}
	return &requisition, nil
}

func (r *RequisitionRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Requisition, error) {
	var requisitions []domain.Requisition
	result := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&requisitions)
	if result.Error != nil {
		return nil, repository.NewRequisitionError("get_by_user_id", result.Error, map[string]interface{}{
			"user_id": userID,
		})
	}
	return requisitions, nil
}

func (r *RequisitionRepository) GetByUserIDAndInstitutionID(ctx context.Context, userID uuid.UUID, institutionID string) (*domain.Requisition, error) {
	var requisition domain.Requisition
	result := r.db.WithContext(ctx).Where("user_id = ? AND institution_id = ?", userID, institutionID).First(&requisition)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, repository.NewRequisitionError("get_by_user_and_institution", repository.ErrRequisitionNotFound, map[string]interface{}{
				"user_id":        userID,
				"institution_id": institutionID,
			})
		}
		return nil, repository.NewRequisitionError("get_by_user_and_institution", result.Error, map[string]interface{}{
			"user_id":        userID,
			"institution_id": institutionID,
		})
	}
	return &requisition, nil
}
