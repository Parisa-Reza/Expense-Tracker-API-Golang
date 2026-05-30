package controllers

import "github.com/beego/beego/v2/server/web"

// HealthController handles health check related endpoints.
// It embeds web.Controller to gain access to Beego's
// request/response handling capabilities.
type HealthController struct {
	web.Controller
}

// Check handles GET /api/v1/health
// It returns the current status of the server.
//
// @Title       Health Check
// @Summary     Returns server health status
// @Description Checks whether the API server is up and running
// @Tags        Health
// @Produce     json
// @Success     200  {object}  map[string]interface{}  "Server is running"
// @Router      /api/v1/health [get]
func (c *HealthController) Check() {
	c.Data["json"] = map[string]interface{}{
		"success": true,
		"message": "Server is running",
	}
	c.ServeJSON()
}