package server

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log-snare/web/data"
	"log-snare/web/service"
	"strings"
)

type SettingsHandler struct {
	SettingsService *service.SettingsService
	UserService     *service.UserService
}

// NewSettingsHandler initializes a new user handler with the given user service
func NewSettingsHandler(us *service.SettingsService, uss *service.UserService) *SettingsHandler {
	return &SettingsHandler{SettingsService: us, UserService: uss}
}

func (h *SettingsHandler) Settings(c *gin.Context) {
	session := sessions.Default(c)
	val := session.Get("user")
	if val == nil {
		c.Redirect(302, "/login")
		return
	}

	user, ok := val.(data.UserSafe)
	if !ok {
		c.Redirect(302, "/login")
		return
	}

	users := h.UserService.UsersByCompanyId(user.CompanyId)

	userInitial := ""
	if len(user.Username) > 0 {
		userInitial = string(strings.ToUpper(user.Username)[0])
	}

	c.HTML(200, "settings.html", gin.H{
		"CurrentRoute":        "/settings",
		"Users":               users,
		"UserCompanyId":       user.CompanyId,
		"ManageUserCompanyId": manageUserCompanyIdentifier(user.CompanyId),

		// common data can be moved to middleware
		"CompanyName":       user.CompanyName,
		"UserInitial":       userInitial,
		"UserRole":          user.Role,
		"ValidationEnabled": data.ValidationEnabled(),
	})
}

func (h *SettingsHandler) Docs(c *gin.Context) {
	session := sessions.Default(c)
	val := session.Get("user")
	if val == nil {
		c.Redirect(302, "/login")
		return
	}

	user, ok := val.(data.UserSafe)
	if !ok {
		c.Redirect(302, "/login")
		return
	}

	userInitial := ""
	if len(user.Username) > 0 {
		userInitial = string(strings.ToUpper(user.Username)[0])
	}

	c.HTML(200, "docs.html", gin.H{
		"CurrentRoute":        "/docs",
		"UserCompanyId":       user.CompanyId,
		"ManageUserCompanyId": manageUserCompanyIdentifier(user.CompanyId),

		// common data can be moved to middleware
		"CompanyName":       user.CompanyName,
		"UserInitial":       userInitial,
		"UserRole":          user.Role,
		"ValidationEnabled": data.ValidationEnabled(),
	})
}

func (h *SettingsHandler) EnableValidation(c *gin.Context) {
	data.SetValidation(true)
	c.JSON(200, gin.H{
		"success": true,
	})
}

func (h *SettingsHandler) DisableValidation(c *gin.Context) {
	data.SetValidation(false)
	c.JSON(200, gin.H{
		"success": true,
	})
}
