package email

import (
	"crypto/tls"
	"time"

	gomail "gopkg.in/mail.v2"
)

type Message struct {
	To []string
	Cc []string
	Subject string
	BodyText string
	BodyHTML string
	AttachFiles []string
}

type Mailer struct {
	Host string
	Port int
	DisplayName string
	Username string
	Password string
}

func NewMailer(host string, port int, username, password, displayName string) *Mailer {
	return &Mailer{
		Host: host,
		Port: port,
		DisplayName: displayName,
		Username: username,
		Password: password,
	}
}

func (mailer Mailer) Send(m *Message) error {
	message := gomail.NewMessage()

	// Set E-Mail sender
	// message.SetHeader("From", fmt.Sprintf("%s <%s>", mailer.DisplayName, mailer.Username))
	message.SetHeader("From", message.FormatAddress(mailer.Username, mailer.DisplayName))
	
	// Set E-Mail receivers
	message.SetHeader("To", m.To...)

	// Set E-Mail receivers
	message.SetHeader("Cc", m.Cc...)

	// Set E-Mail subject
	message.SetHeader("Subject", m.Subject)

	for _, val := range m.AttachFiles {
		message.Attach(val)
	}

	// Set E-Mail body. You can set plain text or html with text/html
	if m.BodyText != "" {
		message.SetBody("text/plain", m.BodyText)
	} else {
		message.SetBody("text/html", m.BodyHTML)
	}

	// Settings for SMTP server
	d := gomail.NewDialer(mailer.Host, mailer.Port, mailer.Username, mailer.Password)

	d.Timeout = 20 * time.Second

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	return d.DialAndSend(message)
}

/*
func SendMail() {
	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", "tgl.info.2024@gmail.com")

	// Set E-Mail receivers
	m.SetHeader("To", "huynhtrongtien89@gmail.com", "ttien1989@gmail.com")

	// Set E-Mail subject
	m.SetHeader("Subject", "Gomail test subject")

	// Set E-Mail body. You can set plain text or html with text/html
	m.SetBody("text/plain", "Dear Customers! \r\nThis is internal mail to test new feature \r\n Please don't reply! \r\n Have a nice day! \r\n Thanks You!")

	// Settings for SMTP server
	d := gomail.NewDialer("smtp.gmail.com", 587, "tgl.info.2024@gmail.com", "vbxz xgpy isao juzy")

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		fmt.Println("Send mail failed", err)
		panic(err)
	}

	fmt.Printf("Send mail success")
}
*/
