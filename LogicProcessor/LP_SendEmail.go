package LogicProcessor
import (
	"net/smtp"
	"strings"
)

func SendToMail(user, pwd, host, to, subject, body, mailtype string) error {
	hp := strings.Split(host, ":")

	auth := smtp.PlainAuth("", user, pwd, hp[0])

	var content_type string

	if mailtype == "html" {

		content_type = "Content-Type: text/" + mailtype + "; charset=UTF-8"

	} else {
		content_type = "Content-Type: text/plain" + "; charset=UTF-8"

	}

	msg := []byte("To: " + to + "\r\nFrom: " + user + "\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + body)
	send_to := strings.Split(to, ";")
	err := smtp.SendMail(host, auth, user, send_to, msg)
	return err

}
