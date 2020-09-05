package useemailcapster

import (
	"fmt"
	templateemail "nuryanto2121/dynamic_rest_api_go/pkg/email"
	util "nuryanto2121/dynamic_rest_api_go/pkg/utils"
	"strings"
)

type RegisterCapster struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"passwrod"`
}

func (R *RegisterCapster) SendRegisterCapster() error {
	subjectEmail := "Activation Code"
	fmt.Printf(subjectEmail)
	err := util.SendEmail(R.Email, subjectEmail, getVerifyBody(R))
	if err != nil {
		return err
	}
	return nil
}

func getVerifyBody(R *RegisterCapster) string {
	verifyHTML := templateemail.SendPasswordCapster

	verifyHTML = strings.ReplaceAll(verifyHTML, `{Name}`, R.Name)
	verifyHTML = strings.ReplaceAll(verifyHTML, `{PasswordCode}`, R.Password)
	return verifyHTML
}
