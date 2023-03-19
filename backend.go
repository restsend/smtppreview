package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path"
	"time"

	"github.com/emersion/go-smtp"
	"github.com/restsend/carrot"
	"gorm.io/gorm"
)

type Attachment struct {
	Name string `json:"name,omitempty"`
	Path string `json:"path"`
	Size int    `json:"size"`
}

// Mail
type Mail struct {
	ID            string       `json:"id" gorm:"size:20;primarykey"`
	CreatedAt     time.Time    `json:"createdAt" gorm:"index"`
	UpdatedAt     time.Time    `json:"updatedAt" gorm:"index"`
	Opened        bool         `json:"opened" gorm:"index"`
	OpenAt        sql.NullTime `json:"openAt"  gorm:"index"`
	Score         int          `json:"score" gorm:"index"`
	From          string       `json:"from" gorm:"size:100;index"`
	To            string       `json:"to"`
	Size          int          `json:"size"`
	Subject       string       `json:"subject,omitempty"`
	Attachments   string       `json:"attachments,omitempty"`
	EmbeddedFiles string       `json:"embeddedFiles,omitempty"`
	RemoteAddr    string       `json:"remoteAddr,omitempty" gorm:"size:100"`
	AuthName      string       `json:"authName,omitempty" gorm:"size:200"`
	TextBody      string       `json:"textBody,omitempty"`
	RawBody       []byte       `json:"-" gorm:"-"`
}

func (m *Mail) IsValid() bool {
	return m.From != "" && m.To != ""
}

func (m *Mail) Flush(db *gorm.DB, mailDir string, perm os.FileMode) (string, error) {
	name := path.Join(mailDir, m.ID+".eml")
	err := os.WriteFile(name, m.RawBody, perm)
	if err != nil {
		return name, err
	}
	return name, db.Save(m).Error
}

func (m *Mail) AfterDelete(tx *gorm.DB) (err error) {
	mailDir := carrot.GetValue(tx, key_MAIL_DIR)
	name := path.Join(mailDir, m.ID+".eml")
	err = os.Remove(name)
	if err != nil {
		log.Println("remove file fail", name, err)
		return err
	}

	var attachments []Attachment
	json.Unmarshal([]byte(m.Attachments), &attachments)
	for _, att := range attachments {
		name = path.Join(mailDir, att.Path)
		os.Remove(name)
	}

	var embeddedFiles []Attachment
	json.Unmarshal([]byte(m.EmbeddedFiles), &embeddedFiles)
	for _, efile := range embeddedFiles {
		name = path.Join(mailDir, efile.Path)
		os.Remove(name)
	}
	return
}

// The Backend implements SMTP server methods.
type Backend struct {
	db         *gorm.DB
	FilePerm   os.FileMode
	MailDir    string
	AuthStrict bool
	Logout     io.Writer
}

func (b *Backend) Prepare() (err error) {
	err = carrot.MakeMigrates(b.db, []any{
		&Mail{},
	})

	if err != nil {
		return
	}

	st, err := os.Stat(b.MailDir)
	if err != nil {
		log.Println("Prepare maildir:", b.MailDir)
		return os.MkdirAll(b.MailDir, fs.ModePerm)
	}

	if !st.IsDir() {
		return fmt.Errorf("maildir: %s is not directory", b.MailDir)
	}
	return
}

// NewSession is called after client greeting (EHLO, HELO).
func (b *Backend) NewSession(c *smtp.Conn) (smtp.Session, error) {
	s := &Session{
		b:      b,
		c:      c,
		Logger: log.New(b.Logout, "", logFlags),
	}
	s.Logger.Println(s.ID(), "Incoming session")
	return s, nil
}

func (bkd *Backend) AuthUser(username, password string) (bool, error) {
	u, err := carrot.GetUserByEmail(bkd.db, username)
	if err != nil {
		return false, err
	}
	return carrot.CheckPassword(u, password), nil
}
