package apis

import (
	"net/smtp"
	"simpleRest/app"
	"simpleRest/utils"
	"strconv"
)

func SendEmailCode(rs app.RequestScope, to, code string, userId int) {
	template := app.Config.EmailVerificationTemplate
	subject := "Email verification"

	body := utils.ReplacePlaceholders(template, map[string]interface{}{
		"url": app.Config.ServerAddress + "/v1/users/verify/" + code + "?user_id=" + strconv.Itoa(userId),
	})

	go send(rs, to, subject, body)
}

func SendPasswordCode(rs app.RequestScope, to, code string, userId int) {
	template := app.Config.ResetPasswordTemplate
	subject := "Reset password"

	body := utils.ReplacePlaceholders(template, map[string]interface{}{
		"url": app.Config.ServerAddress + "/v1/users/reset/" + code + "?user_id=" + strconv.Itoa(userId),
	})

	go send(rs, to, subject, body)
}

func send(rs app.RequestScope, to, subject, body string) {
	from := app.Config.EmailUser
	pass := app.Config.EmailPassword
	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))
	if err != nil {
		rs.Errorf("%s", err)
	}
}
