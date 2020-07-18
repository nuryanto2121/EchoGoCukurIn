package useemailauth

import (
	"fmt"
	templateemail "nuryanto2121/dynamic_rest_api_go/pkg/email"
	util "nuryanto2121/dynamic_rest_api_go/pkg/utils"
	"strings"
)

type Register struct {
	Email      string `json:"email"`
	Name       string `json:"name"`
	GenerateNo string `json:"generate_no"`
}

func (R *Register) SendRegister() error {
	subjectEmail := "Verifikasi Code"
	fmt.Printf(subjectEmail)
	err := util.SendEmail(R.Email, subjectEmail, getVerifyBody(R))
	if err != nil {
		return err
	}
	return nil
}

func getVerifyBody(R *Register) string {
	verifyHTML := templateemail.VerifyCode

	verifyHTML = strings.ReplaceAll(verifyHTML, `{Name}`, R.Name)
	verifyHTML = strings.ReplaceAll(verifyHTML, `{GenerateCode}`, R.GenerateNo)
	return verifyHTML
}
