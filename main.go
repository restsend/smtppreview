package main

import (
	"embed"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/emersion/go-smtp"
	"github.com/gin-gonic/gin"
	"github.com/restsend/carrot"
)

const key_WRITE_TIMEOUT = "WRITE_TIMEOUT"
const key_READ_TIMEOUT = "READ_TIMEOUT"
const key_DOMAIN = "DOMAIN"
const key_MAX_MESSAGE_BYTES = "MAX_MESSAGE_BYTES"
const key_MAX_RECIPIENTS = "MAX_RECIPIENTS"
const key_AUTH_STRICT = "AUTH_STRICT"
const key_MAIL_DIR = "MAIL_DIR"
const key_MAIL_UUID_SIZE = "MAIL_UUID_SIZE"

const fileMode os.FileMode = 0666

func checkValue(key, defaultVal string) string {
	v := carrot.GetEnv(key)
	if v == "" {
		return defaultVal
	}
	return v
}

//go:embed ui/dist/assets/*
var assets embed.FS

//go:embed ui/dist/index.html
var indexHtml string

func main() {

	log.Default().SetFlags(logFlags)

	var smtpServerAddr string
	var httpServerAddr string
	var logFile string
	var dbName string

	defaultSMTPAddr := checkValue("SMTP_ADDR", "0.0.0.0") + ":" + checkValue("SMTP_PORT", "9025")
	defaultHTTPAddr := checkValue("HTTP_ADDR", "0.0.0.0") + ":" + checkValue("HTTP_PORT", "8000")
	defaultDb := checkValue("DB_NAME", "preview.db")

	flag.StringVar(&smtpServerAddr, "smtp", defaultSMTPAddr, "SMTP listen addr, default: "+defaultSMTPAddr)
	flag.StringVar(&httpServerAddr, "http", defaultHTTPAddr, "HTTP listen addr, default: "+defaultHTTPAddr)
	flag.StringVar(&logFile, "log", "", "log file")
	flag.StringVar(&dbName, "db", defaultDb, "DB file, default: "+defaultDb)

	flag.Parse()

	var lw io.Writer = os.Stdout
	var err error
	if logFile != "" {
		lw, err = os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, fileMode)
		if err != nil {
			log.Printf("open %s fail, %v\n", logFile, err)
		}
	}

	db, err := carrot.InitDatabase(lw, "", dbName)
	if err != nil {
		panic(err)
	}
	err = carrot.InitMigrate(db)
	if err != nil {
		log.Panic("init carrot migrate fail", err)
	}

	carrot.CheckValue(db, key_MAIL_DIR, "./mails/")
	carrot.CheckValue(db, key_DOMAIN, "localhost")
	carrot.CheckValue(db, key_AUTH_STRICT, "false")
	carrot.CheckValue(db, key_MAIL_UUID_SIZE, "10")

	be := &Backend{
		db:         db,
		MailDir:    carrot.GetValue(db, key_MAIL_DIR),
		AuthStrict: carrot.GetBoolValue(db, key_AUTH_STRICT),
		FilePerm:   fileMode,
		Logout:     lw,
	}

	err = be.Prepare()
	if err != nil {
		log.Panic("Backend prepare fail", err)
	}

	s := smtp.NewServer(be)
	s.Addr = smtpServerAddr
	s.Domain = carrot.GetValue(db, key_DOMAIN)
	s.WriteTimeout = time.Duration(carrot.GetIntValue(db, key_WRITE_TIMEOUT, 10)) * time.Second
	s.ReadTimeout = time.Duration(carrot.GetIntValue(db, key_READ_TIMEOUT, 10)) * time.Second
	s.MaxMessageBytes = carrot.GetIntValue(db, key_MAX_MESSAGE_BYTES, 1024*1024*10)
	s.MaxRecipients = carrot.GetIntValue(db, key_MAX_RECIPIENTS, 50)
	s.AllowInsecureAuth = true
	s.AuthDisabled = false

	log.Println("Starting SMTP server at", s.Addr)
	log.Println("Domain:", s.Domain)
	log.Println("MailDir:", be.MailDir)
	log.Println("AuthStrict:", be.AuthStrict, "AuthDisabled:", s.AuthDisabled, "AllowInsecureAuth:", s.AllowInsecureAuth)
	log.Println("WriteTimeout:", s.WriteTimeout, "ReadTimeout:", s.ReadTimeout)
	log.Println("MaxMessageBytes:", s.MaxMessageBytes, "MaxRecipients:", s.MaxRecipients)

	go func() {
		if err := s.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	r := gin.Default()
	if err := carrot.InitCarrot(db, r); err != nil {
		panic(err)
	}

	// Embed static file
	r.GET("/assets/*filepath", func(ctx *gin.Context) {
		p := path.Join("ui/dist/", ctx.Request.RequestURI)
		ctx.FileFromFS(p, http.FS(assets))
	})
	r.GET("/", func(ctx *gin.Context) {
		ctx.Data(http.StatusOK, "text/html", []byte(indexHtml))
	})

	RegisterHandlers(r.Group("/api/"), be)
	r.Run(httpServerAddr)
}
