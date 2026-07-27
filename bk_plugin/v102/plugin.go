package v102

import (
	"fmt"

	"github.com/TencentBlueKing/bk-plugin-framework-go/kit"
)

var InputsFormJSON = []byte(`{
	"type": "object",
	"properties": {
		"hello": {
			"type": "string",
			"title": "Hello",
			"ui:component": {"name": "bk-input", "props": {"type": "textarea"}}
		},
		"data_api_tenant_id": {
			"type": "string",
			"title": "Data API 获取的租户 ID",
			"ui:component": {
				"name": "select",
				"props": {
					"remoteConfig": {
						"url": "{{ $context.get('bk_plugin_api_host')['bkplugin-demo-go'] + 'bk_plugin/plugin_api/tenant_id?app_tenant_mode=global' }}"
					},
					"clearable": false
				}
			},
			"ui:reactions": [
				{
					"lifetime": "init",
					"then": {
						"actions": ["{{ $loadDataSource }}"]
					}
				}
			]
		}
	},
	"required": ["hello", "data_api_tenant_id"]
}`)

type Inputs struct {
	Hello           string `json:"hello" jsonschema:"title=Hello"`
	DataAPITenantID string `json:"data_api_tenant_id" jsonschema:"title=Data API 获取的租户 ID"`
}

type ContextInputs struct {
	Executor string `json:"executor" jsonschema:"title=任务执行人"`
	TenantID string `json:"tenant_id" jsonschema:"title=上下文租户 ID"`
}

type Outputs struct {
	World string `json:"world" jsonschema:"title=World"`
}

type Plugin struct{}

func (p *Plugin) Version() string {
	return "1.0.2"
}

func (p *Plugin) Desc() string {
	return "蓝鲸 Go 插件 Demo - 上下文与 Data API 租户 ID 验证"
}

func (p *Plugin) Execute(c *kit.Context) error {
	var inputs Inputs
	if err := c.ReadInputs(&inputs); err != nil {
		return err
	}

	var contextInputs ContextInputs
	if err := c.ReadContextInputs(&contextInputs); err != nil {
		return err
	}

	logger := c.WithField("tenant_id", contextInputs.TenantID).
		WithField("data_api_tenant_id", inputs.DataAPITenantID)
	if contextInputs.TenantID == "" || inputs.DataAPITenantID == "" {
		logger.Error("tenant context validation failed")
		return fmt.Errorf("tenant_id and data_api_tenant_id are required")
	}
	if contextInputs.TenantID != inputs.DataAPITenantID {
		logger.Error("tenant context validation failed")
		return fmt.Errorf(
			"tenant_id mismatch: context=%q, data_api=%q",
			contextInputs.TenantID,
			inputs.DataAPITenantID,
		)
	}

	logger.Info("tenant context validation succeeded")
	return c.WriteOutputs(&Outputs{World: inputs.Hello})
}
