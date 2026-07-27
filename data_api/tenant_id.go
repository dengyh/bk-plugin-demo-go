package data_api

import (
	"strings"

	"github.com/TencentBlueKing/beego-runtime/controllers"
	"github.com/TencentBlueKing/beego-runtime/routers"
	log "github.com/sirupsen/logrus"
)

const (
	HeaderTenantID  = "X-Bkapi-Tenant-Id"
	HeaderRequestID = "X-Bkapi-Request-Id"
)

type TenantOption struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type TenantIDResponse struct {
	Result  bool           `json:"result"`
	Message string         `json:"message"`
	Data    []TenantOption `json:"data"`
}

type TenantIDController struct {
	controllers.PluginApiController
}

func (c *TenantIDController) Get() {
	tenantID := strings.TrimSpace(c.Ctx.Request.Header.Get(HeaderTenantID))
	logger := log.WithFields(log.Fields{
		"request_id": c.Ctx.Request.Header.Get(HeaderRequestID),
		"tenant_id":  tenantID,
	})

	if tenantID == "" {
		logger.Error("data api tenant context missing")
		c.Data["json"] = &TenantIDResponse{
			Result:  false,
			Message: HeaderTenantID + " header is required",
			Data:    nil,
		}
		c.ServeJSON()
		return
	}

	logger.Info("data api tenant context received")
	c.Data["json"] = &TenantIDResponse{
		Result:  true,
		Message: "",
		Data: []TenantOption{
			{
				Label: "当前租户：" + tenantID,
				Value: tenantID,
			},
		},
	}
	c.ServeJSON()
}

func init() {
	routers.PluginApiNamespace.Router("tenant_id", &TenantIDController{})
}
