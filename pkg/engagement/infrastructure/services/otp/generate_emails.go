package otp

import (
	"bytes"
	"html/template"
)

//GenerateEmailFunc generates the custom email to be sent to the user when sending the OTP through email
func GenerateEmailFunc(code string) string {

	t := template.Must(template.New("sendmail").Parse(SendOtpToEmailTemplate))
	buf := new(bytes.Buffer)
	_ = t.Execute(buf, code)
	return buf.String()
}
