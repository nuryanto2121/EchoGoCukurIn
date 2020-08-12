package usesafileupload

import (
	"context"
	ifileupload "nuryanto2121/dynamic_rest_api_go/interface/fileupload"
	"nuryanto2121/dynamic_rest_api_go/models"
	"time"
)

type useSaFileUpload struct {
	repoSaFileUpload ifileupload.Repository
	contextTimeOut   time.Duration
}

func NewSaFileUpload(a ifileupload.Repository, timeout time.Duration) ifileupload.UseCase {
	return &useSaFileUpload{
		repoSaFileUpload: a,
		contextTimeOut:   timeout,
	}
}

func (u *useSaFileUpload) CreateFileUpload(ctx context.Context, data *models.SaFileUpload) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()
	var (
		err error
	)

	err = u.repoSaFileUpload.CreateFileUpload(ctx, data)
	if err != nil {
		return err
	}

	return nil
}
func (u *useSaFileUpload) GetBySaFileUpload(ctx context.Context, fileID int) (models.SaFileUpload, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()
	var (
		err    error
		result models.SaFileUpload
	)
	result, err = u.GetBySaFileUpload(ctx, fileID)
	if err != nil {
		return result, err
	}
	return result, nil
}
