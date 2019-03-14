package common

import (
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

type Message struct {
	From        string
	To          []string
	Subject     string
	Body        string
	ContentType string
}

func Sendmail(key string, sm Message) error {
	if key == "" {
		log.Infof("NOT sending message %v", sm)
		return nil
	}
	m := gomail.NewMessage()
	m.SetHeader("From", sm.From)
	user := "apikey"
	pass := key

	m.SetHeader("To", sm.To...)
	m.SetHeader("Bcc", "jack@sofialondonmoskva.com")
	m.SetHeader("Subject", sm.Subject)
	if sm.ContentType == "" {
		sm.ContentType = "text/plain"
	}
	m.SetBody(sm.ContentType, sm.Body)

	d := gomail.NewDialer("smtp.sendgrid.net", 465, user, pass)

	return d.DialAndSend(m)
}

type EmailQueueItem struct {
	ID     uint64 `gorm:"primary_key"`
	UserID uint64 `gorm:"type:bigint not null REFERENCES users(id)"`
	UUID   string `gorm:"type:varchar(255) not null unique"`

	EmailText    string `gorm:"not null;type:text"`
	EmailSubject string `gorm:"not null;type:text"`
	LastError    string `gorm:"type:text"`

	UserScore      uint64
	SentAt         time.Time
	Sent           bool
	AcknowledgedAt time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
