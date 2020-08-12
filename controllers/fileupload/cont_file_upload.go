package contfileupload

import (
	"context"
	"fmt"
	"io"
	"net/http"
	ifileupload "nuryanto2121/dynamic_rest_api_go/interface/fileupload"
	midd "nuryanto2121/dynamic_rest_api_go/middleware"
	"nuryanto2121/dynamic_rest_api_go/models"
	app "nuryanto2121/dynamic_rest_api_go/pkg"
	"nuryanto2121/dynamic_rest_api_go/pkg/file"
	"nuryanto2121/dynamic_rest_api_go/pkg/logging"
	tool "nuryanto2121/dynamic_rest_api_go/pkg/tools"
	util "nuryanto2121/dynamic_rest_api_go/pkg/utils"
	"os"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
)

// ContFileUpload :
type ContFileUpload struct {
	useSaFileUpload ifileupload.UseCase
}

// NewContFileUpload :
func NewContFileUpload(e *echo.Echo, useSaFileUpload ifileupload.UseCase) {
	cont := &ContFileUpload{
		useSaFileUpload: useSaFileUpload,
	}

	e.Static("/wwwroot", "wwwroot")
	r := e.Group("/api/fileupload")
	// Configure middleware with custom claims
	r.Use(midd.Versioning)
	r.Use(midd.JWT)
	r.POST("", cont.CreateImage)

}

// CreateImage :
// @Summary File Upload
// @Security ApiKeyAuth
// @Description Upload file
// @Tags FileUpload
// @Accept  multipart/form-data
// @Produce json
// @Param OS header string true "OS Device"
// @Param Version header string true "OS Device"
// @Param upload_file formData file true "account image"
// @Param path formData string true "path images"
// @Success 200 {object} tool.ResponseModel
// @Router /api/fileupload [post]
func (u *ContFileUpload) CreateImage(e echo.Context) (err error) {
	ctx := e.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		appE          = tool.Res{R: e}
		imageFormList []models.SaFileUpload
		logger        = logging.Logger{}
	)

	form, err := e.MultipartForm()
	if err != nil {
		return err
	}
	images := form.File["upload_file"]

	pt := form.Value["path"]

	logger.Info(pt)
	//directory api
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	claims, err := app.GetClaims(e)
	if err != nil {
		return appE.ResponseError(http.StatusBadRequest, fmt.Sprintf("%v", err), nil)
	}
	for i, image := range images {
		// Source
		src, err := image.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		var dir_file = dir + "/wwwroot/uploads"
		var path_file = "/wwwroot/uploads"
		err = file.IsNotExistMkDir(dir_file)
		if err != nil {
			return err
		}

		// create folder directory if not exist from param
		if pt[i] != "" {

			dirx := dir_file
			for _, val := range strings.Split(pt[i], "/") {
				dirx = dirx + "/" + val
				fmt.Printf(dirx)
				err = file.IsNotExistMkDir(dirx)
				if err != nil {
					return err
				}
			}
			dir_file = fmt.Sprintf("%s/%s", dir_file, pt[i])

			path_file = fmt.Sprintf("%s/%s", path_file, pt[i])
		}

		fileNameAndUnix := fmt.Sprintf("%d_%s", util.GetTimeNow().Unix(), image.Filename)

		// Destination
		dest := fmt.Sprintf("%s/%s", dir_file, fileNameAndUnix)
		dst, err := os.Create(dest)
		if err != nil {
			return err
		}
		defer dst.Close()

		// Copy
		if _, err = io.Copy(dst, src); err != nil {
			return err
		}

		// r := c.Request()
		// user := c.Get("user").(*jwt.Token)
		// claims := user.Claims.(*util.Claims)

		var imageForm models.SaFileUpload
		// fileName := fmt.Sprintf("%s://%s/upload/%s", c.Scheme(), r.Host, fileNameAndUnix)
		imageForm.FileName = fileNameAndUnix
		imageForm.FilePath = fmt.Sprintf("%s/%s", path_file, fileNameAndUnix)
		imageForm.FileType = filepath.Ext(fileNameAndUnix)
		imageForm.UserInput = claims.UserName
		imageForm.UserEdit = claims.UserName
		// err = u.useSaFileUpload.CreateSaFileUpload(ctx, &imageForm)
		err = u.useSaFileUpload.CreateFileUpload(ctx, &imageForm)

		if err != nil {
			return err
		}
		imageFormList = append(imageFormList, imageForm)

	}
	return appE.Response(http.StatusOK, "Ok", imageFormList)

	// return c.JSON(http.StatusCreated, models.ResponseImage(http.StatusCreated, imageFormList))
}
