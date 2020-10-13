package useemailauth

import (
	"fmt"
	templateemail "nuryanto2121/dynamic_rest_api_go/pkg/email"
	util "nuryanto2121/dynamic_rest_api_go/pkg/utils"
	"strings"
)

type Forgot struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	OTP   string `json:"otp"`
}

func (F *Forgot) SendRegister() error {
	subjectEmail := "Lupa Password"
	fmt.Printf(subjectEmail)
	err := util.SendEmail(F.Email, subjectEmail, getInformasiLoginBodyForgot(F))
	if err != nil {
		return err
	}
	return nil
}

func getInformasiLoginBodyForgot(F *Forgot) string {
	verifyHTML := templateemail.SendRegister

	verifyHTML = strings.ReplaceAll(verifyHTML, `{Name}`, F.Name)
	verifyHTML = strings.ReplaceAll(verifyHTML, `{Email}`, F.Email)
	verifyHTML = strings.ReplaceAll(verifyHTML, `{PasswordCode}`, F.OTP)
	return verifyHTML
}