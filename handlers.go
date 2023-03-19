package main

import (
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/restsend/carrot"
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

	carrot.RegisterObject(routes, &carrot.WebObject[Mail]{
		Model:     Mail{},
		Name:      "mail",
		Editables: []string{"Opened", "OpenAt"},
		Filters:   []string{"Opened", "OpenAt", "Score"},
		Orders:    []string{"CreatedAt", "OpenAt", "Score"},
		Searchs:   []string{"Subject", "From", "To"},
		GetDB: func(ctx *gin.Context, isCreate bool) *gorm.DB {
			return be.db
		},
		AllowMethods: carrot.GET | carrot.DELETE | carrot.EDIT | carrot.QUERY,
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
