package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	_ "embed"

	"github.com/gin-gonic/gin"
	"github.com/restsend/carrot"
	"github.com/restsend/gormpher"
	"github.com/restsend/parsemail"
	"gorm.io/gorm"
)

type Summary struct {
	TotalCount int `json:"totalCount"`
	TotalSize  int `json:"totalSize"`
}

func RegisterHandlers(r *gin.RouterGroup, be *Backend) {
	routes := r.Use(func(ctx *gin.Context) {
		ctx.Set("backend", be)
		ctx.Next()
	})

	routes.StaticFS("/raw", http.Dir(be.MailDir))
	routes.POST("/summary", handleSummary)
	routes.POST("/config", handleGetConfig)
	routes.POST("/config/edit", handleEditConfig)
	routes.GET("/render/:id", handleRender)

	gormpher.RegisterObject(routes, &gormpher.WebObject{
		Model:       Mail{},
		Name:        "mail",
		Editables:   []string{"Opened", "OpenAt"},
		Filterables: []string{"Opened", "OpenAt", "Score"},
		Orderables:  []string{"CreatedAt", "OpenAt", "Score"},
		Searchables: []string{"Subject", "From", "To"},
		GetDB: func(ctx *gin.Context, isCreate bool) *gorm.DB {
			return be.db
		},
		AllowMethods: gormpher.GET | gormpher.DELETE | gormpher.EDIT | gormpher.QUERY,
	})
}

func getBackend(c *gin.Context) *Backend {
	obj, _ := c.Get("backend")
	return obj.(*Backend)
}

func handleSummary(c *gin.Context) {
	be := getBackend(c)

	summary := Summary{}
	filepath.Walk(be.MailDir, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(info.Name(), ".eml") {
			summary.TotalCount++
		}
		// All file
		summary.TotalSize += int(info.Size())
		return nil
	})
	c.JSON(http.StatusOK, summary)
}

func handleGetConfig(c *gin.Context) {
	be := getBackend(c)
	vals := map[string]any{}

	vals["WRITE_TIMEOUT"] = carrot.GetValue(be.db, key_WRITE_TIMEOUT)
	vals["READ_TIMEOUT"] = carrot.GetValue(be.db, key_READ_TIMEOUT)
	vals["DOMAIN"] = carrot.GetValue(be.db, key_DOMAIN)
	vals["MAX_MESSAGE_BYTES"] = carrot.GetValue(be.db, key_MAX_MESSAGE_BYTES)
	vals["MAX_RECIPIENTS"] = carrot.GetValue(be.db, key_MAX_RECIPIENTS)
	vals["AUTH_STRICT"] = carrot.GetValue(be.db, key_AUTH_STRICT)
	vals["MAIL_DIR"] = carrot.GetValue(be.db, key_MAIL_DIR)

	c.JSON(http.StatusOK, vals)
}

func handleEditConfig(c *gin.Context) {
	be := getBackend(c)
	var vals map[string]string
	err := c.BindJSON(&vals)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	allows := map[string]bool{
		"WRITE_TIMEOUT":     true,
		"READ_TIMEOUT":      true,
		"DOMAIN":            true,
		"MAX_MESSAGE_BYTES": true,
		"MAX_RECIPIENTS":    true,
		"AUTH_STRICT":       true,
		"MAIL_DIR":          true,
	}

	for k, v := range vals {
		if _, ok := allows[k]; ok {
			carrot.SetValue(be.db, k, v)
		}
	}
	c.JSON(http.StatusOK, true)
}

func handleRender(c *gin.Context) {
	be := getBackend(c)
	mailid := c.Param("id")
	var mail Mail
	result := be.db.Take(&mail, "id", mailid)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err": result.Error.Error(),
		})
		return
	}

	name := path.Join(be.MailDir, mailid+".eml")
	data, err := os.ReadFile(name)
	if err != nil {
		log.Println("open eml fail", name, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"err": result.Error.Error(),
		})
		return
	}
	msg, err := parsemail.Parse(bytes.NewReader(data)) // returns Email struct and error
	if err != nil {
		log.Println("parse eml fail", name, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"err": result.Error.Error(),
		})
	}

	if msg.HTMLBody == "" {
		c.Data(http.StatusOK, "text/plain", []byte(msg.TextBody))
		return
	}

	htmlBody := msg.HTMLBody
	for _, efile := range msg.EmbeddedFiles {
		key := "cid:" + efile.CID
		url := fmt.Sprintf("/api/raw/%s-%s", mail.ID, efile.CID)
		htmlBody = strings.ReplaceAll(htmlBody, key, url)
	}
	c.Data(http.StatusOK, "text/html", []byte(htmlBody))
}
