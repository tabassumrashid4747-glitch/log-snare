package server

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log-snare/web/data"
	"log-snare/web/service"
	"strings"
)

type DashboardHandler struct {
	DashboardService *service.DashboardService
}

// NewDashboardHandler initializes a new user handler with the given user service
func NewDashboardHandler(us *service.DashboardService) *DashboardHandler {
	return &DashboardHandler{DashboardService: us}
}

func (h *DashboardHandler) Dashboard(c *gin.Context) {

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

	summaryCounts, err := h.DashboardService.GetSummaryCounts(user.CompanyId)
	if err != nil {
		h.DashboardService.Logger.Error("unable get summary counts", zap.Error(err))
		c.Redirect(500, "/")
		return
	}

	challenges, err := h.DashboardService.GetCompletedChallenges()
	if err != nil {
		h.DashboardService.Logger.Error("unable get challenge counts", zap.Error(err))
		c.Redirect(500, "/")
		return
	}

	userInitial := ""
	if len(user.Username) > 0 {
		userInitial = string(strings.ToUpper(user.Username)[0])
	}

	c.HTML(200, "dashboard.html", gin.H{
		"CurrentRoute":        "/dashboard",
		"EmployeeCount":       summaryCounts.EmployeeCount,
		"AdminCount":          summaryCounts.AdminCount,
		"UserCount":           summaryCounts.UserCount,
		"HighSalaryCount":     summaryCounts.EmployeeHighSalaryCount,
		"CompletedOne":        challenges.OneComplete,
		"CompletedTwo":        challenges.TwoComplete,
		"CompletedThree":      challenges.ThreeComplete,
		"UserCompanyId":       user.CompanyId,
		"ManageUserCompanyId": manageUserCompanyIdentifier(user.CompanyId),

		// common data can be moved to middleware
		"CompanyName":       user.CompanyName,
		"UserInitial":       userInitial,
		"UserRole":          user.Role,
		"ValidationEnabled": data.ValidationEnabled(),
	})
}
