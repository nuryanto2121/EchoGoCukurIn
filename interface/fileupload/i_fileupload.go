package ifileupload

import (
	"context"
	"nuryanto2121/cukur_in_barber/models"
)

type Repository interface {
	CreateFileUpload(ctx context.Context, data *models.SaFileUpload) error
	GetBySaFileUpload(ctx context.Context, fileID int) (models.SaFileUpload, error)
	DeleteSaFileUpload(ctx context.Context, fileID int) error
}

type UseCase interface {
	CreateFileUpload(ctx context.Context, data *models.SaFileUpload) error
	GetBySaFileUpload(ctx context.Context, fileID int) (models.SaFileUpload, error)
}
