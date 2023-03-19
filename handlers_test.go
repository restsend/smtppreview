package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/restsend/carrot"
	"github.com/stretchr/testify/assert"
)

func setupRouter() (*gin.Engine, *Backend) {
	r := gin.Default()
	db, err := carrot.InitDatabase(nil, "", "")
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
	}

	err = be.Prepare()
	if err != nil {
		log.Panic("Backend prepare fail", err)
	}

	RegisterHandlers(r.Group("/api"), be)
	return r, be
}
func TestSummary(t *testing.T) {
	router, _ := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/summary", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
