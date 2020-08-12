package repofileupload

import (
	"context"
	"fmt"
	ifileupload "nuryanto2121/dynamic_rest_api_go/interface/fileupload"
	"nuryanto2121/dynamic_rest_api_go/models"
	"nuryanto2121/dynamic_rest_api_go/pkg/logging"

	"github.com/jinzhu/gorm"
)

type repoAuth struct {
	Conn *gorm.DB
}

func NewRepoFileUpload(Conn *gorm.DB) ifileupload.Repository {
	return &repoAuth{Conn}
}

func (m *repoAuth) CreateFileUpload(ctx context.Context, data *models.SaFileUpload) (err error) {
	var logger = logging.Logger{}
	query := m.Conn.Create(&data)
	logger.Query(fmt.Sprintf("%v", query.QueryExpr())) //cath to log query string
	err = query.Error
	// err = db.Conn.Create(userData).Error
	if err != nil {
		return err
	}
	return nil
}

func (m *repoAuth) GetBySaFileUpload(ctx context.Context, fileID int) (models.SaFileUpload, error) {
	var (
		dataFileUpload = models.SaFileUpload{}
		logger         = logging.Logger{}
		err            error
	)
	query := m.Conn.Where("file_id = ?", fileID).First(&dataFileUpload)
	logger.Query(fmt.Sprintf("%v", query.QueryExpr()))
	err = query.Error

	if err != nil {
		//
		if err == gorm.ErrRecordNotFound {
			return dataFileUpload, models.ErrNotFound
		}
		return dataFileUpload, err
	}

	return dataFileUpload, err
}
func (m *repoAuth) DeleteSaFileUpload(ctx context.Context, fileID int) error {
	var (
		logger = logging.Logger{}
		err    error
	)
	userData := &models.SaFileUpload{}
	userData.FileID = fileID

	query := m.Conn.Delete(&userData)
	logger.Query(fmt.Sprintf("%v", query.QueryExpr())) //cath to log query string
	err = query.Error

	if err != nil {
		return err
	}
	return nil
}
