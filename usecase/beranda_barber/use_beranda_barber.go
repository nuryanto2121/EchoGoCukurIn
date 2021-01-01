package useberandabarber

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"

	iberandabarber "nuryanto2121/cukur_in_barber/interface/beranda_barber"
	ifileupload "nuryanto2121/cukur_in_barber/interface/fileupload"
	"nuryanto2121/cukur_in_barber/models"
	util "nuryanto2121/cukur_in_barber/pkg/utils"
	"time"
)

type useBarber struct {
	repoBeranda    iberandabarber.Repository
	repoFile       ifileupload.Repository
	contextTimeOut time.Duration
}

func NewUserMBarber(a iberandabarber.Repository, c ifileupload.Repository, timeout time.Duration) iberandabarber.Usecase {
	return &useBarber{
		repoBeranda:    a,
		repoFile:       c,
		contextTimeOut: timeout}
}

func (u *useBarber) GetStatusOrder(ctx context.Context, Claims util.Claims, queryparam models.ParamDynamicList) (interface{}, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	ID, err := strconv.Atoi(Claims.UserID)
	if err != nil {
		return nil, err
	}
	result, err := u.repoBeranda.GetStatusOrder(queryparam.ParamView, ID)
	if err != nil {
		if err == models.ErrNotFound {
			return result, nil
		}
		return nil, err
	}

	return result, nil
}
func (u *useBarber) GetListOrder(ctx context.Context, Claims util.Claims, queryparam models.ParamDynamicList) (interface{}, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()
	var (
		result models.ResponseModelList
		err    error
	)
	// var tUser = models.Barber{}
	/*membuat Where like dari struct*/
	if queryparam.Search != "" {
		queryparam.Search = strings.ToLower(fmt.Sprintf("%%%s%%", queryparam.Search))
	}

	queryparam.InitSearch = fmt.Sprintf("owner_id = %s", Claims.UserID)
	result.Data, err = u.repoBeranda.GetListOrder(queryparam)
	if err != nil {
		return result, err
	}

	result.Total, err = u.repoBeranda.Count(queryparam)
	if err != nil {
		return result, err
	}

	// d := float64(result.Total) / float64(queryparam.PerPage)
	result.LastPage = int(math.Ceil(float64(result.Total) / float64(queryparam.PerPage)))
	result.Page = queryparam.Page

	return result, nil
}
