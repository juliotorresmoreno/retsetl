package mail

import "testing"
import "bitbucket.org/mlsdatatools/retsetl/config"

func TestSend(t *testing.T) {
	subject := "subject"
	to := config.EMAIL_ADMIN
	from := config.EMAIL_SEND
	pass := config.EMAIL_PASSWORD
	body := "body"
	err := Send(subject, to, from, pass, body)
	if err != nil {
		t.Error(err)
	}
}
