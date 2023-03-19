package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/DusanKasan/parsemail"
	"github.com/emersion/go-smtp"
	"github.com/restsend/carrot"
)

// A Session is returned after successful login.
type Session struct {
	b      *Backend
	c      *smtp.Conn
	Logger *log.Logger
	_mail  *Mail
}

func (s *Session) ID() string {
	return colorize(ColorBlue, s.c.Conn().RemoteAddr().String())
}
func (s *Session) Current() *Mail {
	if s._mail == nil {
		uuid := carrot.GenUniqueKey(s.b.db.Model(&Mail{}), "id", carrot.GetIntValue(s.b.db, key_MAIL_UUID_SIZE, 8))
		n := time.Now()
		s._mail = &Mail{
			ID:         uuid,
			CreatedAt:  n,
			UpdatedAt:  n,
			Opened:     false,
			RemoteAddr: s.c.Conn().RemoteAddr().String(),
		}
	}
	return s._mail
}

// AuthPlain implements authentication using SASL PLAIN.
func (s *Session) AuthPlain(username, password string) (err error) {
	s.Current().AuthName = username
	result := true
	if s.b.AuthStrict {
		result, err = s.b.AuthUser(username, password)
		if err != nil {
			s.Logger.Println(s.ID(), "Auth err", colorize(ColorGreen, username), err)
			return err
		}
	}
	s.Logger.Println(s.ID(), "Auth", colorize(ColorGreen, username), colorize(ColorRed, password), "->", result)
	return nil
}

func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	s.Logger.Println(s.ID(), "From", colorize(ColorGreen, from))
	s.Current().From = from
	return nil
}

func (s *Session) Rcpt(to string) error {
	s.Logger.Println(s.ID(), "To", colorize(ColorGreen, to))
	if s.Current().To != "" {
		s.Current().To += ","
	}
	s.Current().To += to
	return nil
}

func (s *Session) Data(r io.Reader) (err error) {
	const maxDataShow = 120
	var b []byte
	b, err = io.ReadAll(r)
	if err != nil {
		s.Logger.Println(s.ID(), "Data", "error", err)
		return err
	}

	current := s.Current()
	current.RawBody = b
	current.Size = len(b)

	data := strings.TrimSpace(strings.SplitN(string(b), "\n", 2)[0])
	if len(data) > maxDataShow {
		data = data[0:maxDataShow]
	}
	if len(b) != len(data) {
		data += " ..."
	}
	s.Logger.Println(s.ID(), "Data", "size", colorize(ColorBlue, strconv.Itoa(len(b))), "body", data)

	if current.IsValid() {
		//
		r := bytes.NewReader(current.RawBody)
		msg, err := parsemail.Parse(r) // returns Email struct and error
		if err != nil {
			s.Logger.Println(s.ID(), "Parse error:", err)
			s._mail = nil
			return err
		}

		subject := msg.Header.Get("Subject")
		s.Logger.Println(s.ID(), "Subject", colorize(ColorBlue, subject))

		s._mail.Subject = subject

		s.Logger.Println(s.ID(), "Attachments", colorize(ColorBlue, strconv.Itoa(len(msg.Attachments))))
		var attachments []Attachment

		s._mail.TextBoby = msg.TextBody

		for idx, att := range msg.Attachments {
			ext := path.Ext(att.Filename)
			storepath := fmt.Sprintf("%s-%d%s", s._mail.ID, idx, ext)

			attData, err := io.ReadAll(att.Data)
			if err != nil {
				s.Logger.Println(s.ID(), "Flush error attachment idx:", idx, "filename:", att.Filename, err)
				continue
			}

			err = os.WriteFile(path.Join(s.b.MailDir, storepath), attData, s.b.FilePerm)
			if err != nil {
				s.Logger.Println(s.ID(), "Flush error attachment idx:", idx, "filename:", att.Filename, err)
				continue
			}

			s.Logger.Println(s.ID(), "Attachment", colorize(ColorBlue, path.Join(s.b.MailDir, storepath)), "size", colorize(ColorBlue, strconv.Itoa(len(attData))))
			attachments = append(attachments, Attachment{
				Name: att.Filename,
				Path: storepath,
				Size: len(attData),
			})
		}
		data, _ := json.Marshal(attachments)
		s._mail.Attachments = string(data)
		name, err := current.Flush(s.b.db, s.b.MailDir, s.b.FilePerm)
		if err != nil {
			s.Logger.Println(s.ID(), "Flush error filename:", name, err)
		} else {
			s.Logger.Println(s.ID(), "Flush", colorize(ColorBlue, name), "size", colorize(ColorBlue, strconv.Itoa(len(current.RawBody))))
		}
		s._mail = nil //
	}

	return nil
}

func (s *Session) Reset() {
	s._mail = nil //
}

func (s *Session) Logout() error {
	return nil
}
