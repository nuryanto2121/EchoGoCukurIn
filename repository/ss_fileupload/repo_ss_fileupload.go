package repofileupload

import (
	"context"
	ifileupload "nuryanto2121/dynamic_rest_api_go/interface/fileupload"
	"nuryanto2121/dynamic_rest_api_go/models"
	"nuryanto2121/dynamic_rest_api_go/pkg/logging"
	queryfileupload "nuryanto2121/dynamic_rest_api_go/query/fileupload"

	"github.com/jmoiron/sqlx"
)

type repoAuth struct {
	DB *sqlx.DB
}

func NewRepoFileUpload(Conn *sqlx.DB) ifileupload.Repository {
	return &repoAuth{Conn}
}

func (m *repoAuth) CreateFileUpload(ctx context.Context, data models.SaFileUpload) error {
	var logger = logging.Logger{}
	logger.Query(queryfileupload.QuerySave, data)
	_, err := m.DB.NamedExecContext(ctx, queryfileupload.QuerySave, data)
	if err != nil {
		return err
	}
	return nil
}
