package useemailauth

import (
	"fmt"
	templateemail "nuryanto2121/cukur_in_barber/pkg/email"
	util "nuryanto2121/cukur_in_barber/pkg/utils"
	"strings"
)

type Forgot struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	OTP   string `json:"otp"`
}

func (F *Forgot) SendForgot() error {
	subjectEmail := "Permintaan Lupa Password"
	fmt.Printf(subjectEmail)
	err := util.SendEmail(F.Email, subjectEmail, getInformasiLoginBodyForgot(F))
	if err != nil {
		return err
	}
	return nil
}

func getInformasiLoginBodyForgot(F *Forgot) string {
	verifyHTML := templateemail.VerifyCode

	verifyHTML = strings.ReplaceAll(verifyHTML, `{Name}`, F.Name)
	verifyHTML = strings.ReplaceAll(verifyHTML, `{Email}`, F.Email)
	verifyHTML = strings.ReplaceAll(verifyHTML, `{OTP}`, F.OTP)
	return verifyHTML
}
