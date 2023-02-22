package mail

import (
	"testing"

	"github.com/hoangtk0100/simple-bank/util"
	"github.com/stretchr/testify/require"
)

func TestSendEmailWithGmail(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	config, err := util.LoadConfig("..")
	require.NoError(t, err)

	sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)

	subject := "A test mail"
	content := `
	<h1>Hello world!</h1>
	<p>I'm coming :v</p>
	`
	to := []string{"hoangtk.0100@gmail.com"}
	attachments := []string{"../README.md"}

	err = sender.SendEmail(subject, content, to, nil, nil, attachments)
	require.NoError(t, err)
}
