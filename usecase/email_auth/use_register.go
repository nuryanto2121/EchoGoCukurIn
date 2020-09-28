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
	PasswordCd string `json:"generate_no"`
}

func (R *Register) SendRegister() error {
	subjectEmail := "Informasi Login"
	fmt.Printf(subjectEmail)
	err := util.SendEmail(R.Email, subjectEmail, getInformasiLoginBody(R))
	if err != nil {
		return err
	}
	return nil
}

func getInformasiLoginBody(R *Register) string {
	verifyHTML := templateemail.SendRegister

	verifyHTML = strings.ReplaceAll(verifyHTML, `{Name}`, R.Name)
	verifyHTML = strings.ReplaceAll(verifyHTML, `{Email}`, R.Email)
	verifyHTML = strings.ReplaceAll(verifyHTML, `{PasswordCode}`, R.PasswordCd)
	return verifyHTML
}
