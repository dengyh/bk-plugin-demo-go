package v103

import (
	v102 "bk-plugin-demo-go/bk_plugin/v102"

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
						"url": "{{ $context.get('bk_plugin_api_host')['bkplugin-go2'] + 'bk_plugin/plugin_api/tenant_id' }}"
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

type Inputs = v102.Inputs

type ContextInputs = v102.ContextInputs

type Outputs = v102.Outputs

type Plugin struct{}

func (p *Plugin) Version() string {
	return "1.0.3"
}

func (p *Plugin) Desc() string {
	return "蓝鲸 Go 插件 Demo - 修正 Data API SaaS App Code"
}

func (p *Plugin) Execute(c *kit.Context) error {
	return (&v102.Plugin{}).Execute(c)
}
