package aasemail
 
import (
 
"encoding/base64"
"net/smtp"
"fmt"
"strings"
"io/ioutil"
)
 
type Email struct {
	To []string
	From string
	FromName string
	CC []string
	BCC []string
	Attachments []string
	Body string
	Subject string
	Username string
	Password string
	Server string
	Port string
	Charset string
	AttachedData [][][]byte
}

func NewEmail() *Email {
	return &Email{Port: "25", Charset: "us-ascii"}
}


func (em *Email) AttachFile(fileName string) {
	em.Attachments = append(em.Attachments, fileName)
}

func (em *Email) AttachData(data []byte, fileName string) {
	em.AttachedData = append(em.AttachedData, [][]byte{data, []byte(fileName)})
}
	

func (em *Email) Send() error {
	auth := smtp.PlainAuth("", em.Username, em.Password, em.Server)
	toline := append(em.To, em.CC ...)
	toline = append(toline, em.BCC ...)
	var fromline string
	if em.FromName != "" {
		fromline = fmt.Sprintf(`%s <%s>`, em.FromName, em.From)
	} else {
		fromline = em.From
	}
        boundary := "=====a8jk5mfd9isr77grpv399====="
	conType := `text/plain; charset="` + em.Charset + `"`
	if len(em.Attachments) > 0 || len(em.AttachedData) > 0 {
		conType = `multipart/mixed; boundary="` + boundary + `"`
	}
	m := fmt.Sprintf("Content-Type: %s\r\n", conType)
	m += "MIME-Version: 1.0\r\n"
	m += fmt.Sprintf("From: %s\r\n", fromline)
	m += fmt.Sprintf("To: %s\r\n", strings.Join(em.To, ", "))
	if len(em.CC) > 0 {
		m += fmt.Sprintf("CC: %s\r\n", strings.Join(em.CC, ", "))
	}
	m += fmt.Sprintf("Subject: %s\r\n", em.Subject)

        if len(em.Attachments) > 0 || len(em.AttachedData) > 0 {
		m += "\r\n--" + boundary + "\r\n"
		m += `Content-Type: text/plain; charset="` + em.Charset + `"` + "\r\n"
		m += "MIME-Version: 1.0\r\n"
	}
	m += "\r\n"
	//m += "Content-Transfer-Encoding: 7bit\r\n\r\n"
	m += em.Body + "\r\n\r\n"
	if len(em.Attachments) > 0 {
		for _, e := range em.Attachments {
			m += "--" + boundary + "\r\n"
			fileContents, err := ioutil.ReadFile(e)
			if err != nil {
				return err
			}
			attln := strings.Split(e,"/")
			att := attln[len(attln)-1]
			m += "Content-Type: application/octet-stream\r\n"
			m += "MIME-Version: 1.0\r\n"
			m += "Content-Transfer-Encoding: base64\r\n"
			m += `Content-Disposition: attachment; filename="` + att + `"` + "\r\n"
			m += "\r\n" + base64.StdEncoding.EncodeToString(fileContents)
			m += "\r\n"
			m += "\r\n--" + boundary + "--"
		}
	}
	if len(em.AttachedData) > 0 {
		for _, e := range em.AttachedData {
			m += "--" + boundary + "\r\n"
			fileContents := e[0]
			att := string(e[1])
			m += "Content-Type: application/octet-stream\r\n"
			m += "MIME-Version: 1.0\r\n"
			m += "Content-Transfer-Encoding: base64\r\n"
			m += `Content-Disposition: attachment; filename="` + att + `"` + "\r\n"
			m += "\r\n" + base64.StdEncoding.EncodeToString(fileContents)
			m += "\r\n"
			m += "\r\n--" + boundary + "--"
		}
	}
	err := smtp.SendMail(em.Server + ":" + em.Port, auth, em.From, toline, []byte(m))
	if err != nil {
		return err
	}
        return nil
}
